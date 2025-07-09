package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strconv"
)

type DeploymentService struct {
	pb.UnimplementedDeploymentServer
	uc *biz.DeploymentUseCase
	uf *biz.UserFollowUseCase
}

func NewDeploymentService(uc *biz.DeploymentUseCase, uf *biz.UserFollowUseCase) *DeploymentService {
	return &DeploymentService{
		uc: uc,
		uf: uf,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *DeploymentService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Deployment,
		},
		ListAll: true,
	})
	if err != nil {
		return nil, err
	}

	followMap := make(map[string]bool)
	for _, follow := range userFollows {
		followKey := fmt.Sprintf("%s.%s", follow.ClusterName, follow.FollowValue)
		followMap[followKey] = true
	}
	return followMap, nil
}

// setFollowedStatus 设置关注状态
func (x *DeploymentService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.DeploymentItem) error {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return err
	}
	followMap, err := x.getUserFollowMap(ctx, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		followKey := fmt.Sprintf("%s.%s", clusterName, item.Name)
		item.IsFollowed = followMap[followKey]
	}
	return nil
}

// convertContainers converts the containers of deployment to the pb.Container type
func (x *DeploymentService) convertContainers(containers []corev1.Container) []*corev1.Container {
	result := make([]*corev1.Container, 0, len(containers))
	for _, container := range containers {
		result = append(result, &container)
	}
	return result
}

func (x *DeploymentService) convertDeploymentStrategy(strategy appsv1.DeploymentStrategy) *pb.DeploymentStrategy {
	switch strategy.Type {
	case appsv1.RecreateDeploymentStrategyType:
		return &pb.DeploymentStrategy{
			UpdateStrategyType: pb.DeploymentStrategy_Recreate,
		}
	case appsv1.RollingUpdateDeploymentStrategyType:
		var maxUnavailable, maxSurge string
		if strategy.RollingUpdate.MaxUnavailable == nil {
			maxUnavailable = "25%"
		}
		if strategy.RollingUpdate.MaxSurge == nil {
			maxSurge = "25%"
		} else {
			maxUnavailable = strategy.RollingUpdate.MaxUnavailable.String()
			maxSurge = strategy.RollingUpdate.MaxSurge.String()
		}
		return &pb.DeploymentStrategy{
			UpdateStrategyType: pb.DeploymentStrategy_RollingUpdate,
			MaxUnavailable:     maxUnavailable,
			MaxSurge:           maxSurge,
		}
	default:
		return &pb.DeploymentStrategy{}
	}
}

// convertContainers converts the containers of deployment to the pb.Container type
func (x *DeploymentService) convertStatus(status appsv1.DeploymentStatus) *pb.DeploymentStatus {
	return &pb.DeploymentStatus{
		Replicas:            uint32(status.Replicas),
		UpdatedReplicas:     uint32(status.UpdatedReplicas),
		AvailableReplicas:   uint32(status.AvailableReplicas),
		ReadyReplicas:       uint32(status.ReadyReplicas),
		UnavailableReplicas: uint32(status.UnavailableReplicas),
	}
}

// convertConditions converts the conditions of deployment to the pb.DeploymentCondition type
func (x *DeploymentService) convertConditions(conditions []appsv1.DeploymentCondition) []*pb.DeploymentCondition {
	if len(conditions) == 0 || conditions == nil {
		return []*pb.DeploymentCondition{}
	}
	result := make([]*pb.DeploymentCondition, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, &pb.DeploymentCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastUpdateTime:     uint64(condition.LastUpdateTime.UnixNano() / 1e6),
			LastTransitionTime: uint64(condition.LastTransitionTime.UnixNano() / 1e6),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return result
}

// convertDO2DTO converts the deployment object to the pb.DeploymentItem type
func (x *DeploymentService) convertDO2DTO(deployment *appsv1.Deployment) *pb.DeploymentItem {
	if deployment.APIVersion == "" {
		deployment.APIVersion = "apps/v1"
	}
	if deployment.Kind == "" {
		deployment.Kind = "Deployment"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(deployment)
	return &pb.DeploymentItem{
		Name:           deployment.Name,
		Namespace:      deployment.Namespace,
		Status:         x.convertStatus(deployment.Status),
		Containers:     x.convertContainers(deployment.Spec.Template.Spec.Containers),
		CreateTime:     uint64(deployment.CreationTimestamp.UnixNano() / 1e6),
		Labels:         deployment.Labels,
		Annotations:    deployment.Annotations,
		Selector:       deployment.Spec.Selector,
		UpdateStrategy: x.convertDeploymentStrategy(deployment.Spec.Strategy),
		Replicas:       uint32(deployment.Status.Replicas),
		SpecReplicas:   uint32(*deployment.Spec.Replicas),
		Conditions:     x.convertConditions(deployment.Status.Conditions),
		Yaml:           yamlStr,
	}
}

func (x *DeploymentService) ConvertReplicaSet2DTO(replicaset *appsv1.ReplicaSet) *pb.ReplicaSetItem {
	images := make([]string, 0, len(replicaset.Spec.Template.Spec.Containers))
	for _, container := range replicaset.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}
	return &pb.ReplicaSetItem{
		Name:       replicaset.Name,
		Images:     images,
		CreateTime: uint64(replicaset.CreationTimestamp.UnixNano() / 1e6),
		Version:    replicaset.Annotations["deployment.kubernetes.io/revision"],
	}
}

func (x *DeploymentService) ListDeployment(ctx context.Context, req *pb.ListDeploymentRequest) (*pb.ListDeploymentResponse, error) {
	deployments, total, err := x.uc.ListDeployment(ctx, &biz.ListDeploymentRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ListAll:     utils.IntToBool(int(req.ListAll)),
		Keyword:     req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.DeploymentItem, 0, len(deployments))
	for _, deployment := range deployments {
		list = append(list, x.convertDO2DTO(deployment))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListDeploymentResponse{
		List:  list,
		Total: int32(total),
	}, nil
}

func (x *DeploymentService) CreateOrUpdateDeploymentByYaml(ctx context.Context, req *pb.CreateOrUpdateDeploymentByYamlRequest) (*pb.CreateOrUpdateDeploymentByYamlResponse, error) {
	err := x.uc.CreateOrUpdateDeploymentByYaml(ctx, &biz.CreateOrUpdateDeploymentByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateDeploymentByYamlResponse{}, nil
}

func (x *DeploymentService) DeleteDeployment(ctx context.Context, req *pb.DeleteDeploymentRequest) (*pb.DeleteDeploymentResponse, error) {
	err := x.uc.DeleteDeployment(ctx, &biz.DeploymentCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})

	if err != nil {
		return nil, err
	}
	return &pb.DeleteDeploymentResponse{}, nil
}

func (x *DeploymentService) RestartDeployment(ctx context.Context, req *pb.RestartDeploymentRequest) (*pb.RestartDeploymentResponse, error) {
	err := x.uc.RestartDeployment(ctx, &biz.DeploymentCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})

	if err != nil {
		return nil, err
	}
	return &pb.RestartDeploymentResponse{}, nil
}

func (x *DeploymentService) ScaleDeployment(ctx context.Context, req *pb.ScaleDeploymentRequest) (*pb.ScaleDeploymentResponse, error) {
	err := x.uc.ScaleDeployment(ctx, &biz.ScaleDeploymentRequest{
		DeploymentCommonParams: biz.DeploymentCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Replicas: req.Replicas,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ScaleDeploymentResponse{}, nil
}

func (x *DeploymentService) GetDeploymentDetail(ctx context.Context, req *pb.DeploymentDetailRequest) (*pb.DeploymentDetailResponse, error) {
	deploymentDetail, err := x.uc.GetDeploymentDetail(ctx, &biz.DeploymentCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeploymentDetailResponse{
		Detail: x.convertDO2DTO(deploymentDetail.Deployment),
	}, nil
}

func (x *DeploymentService) RollbackDeployment(ctx context.Context, req *pb.RollbackDeploymentRequest) (*pb.RollbackDeploymentResponse, error) {
	version, err := strconv.ParseUint(req.Version, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %v", err)
	}
	err = x.uc.RollbackDeployment(ctx, &biz.RollbackDeploymentRequest{
		DeploymentCommonParams: biz.DeploymentCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Revision: uint32(version),
	})
	if err != nil {
		return nil, err
	}
	return &pb.RollbackDeploymentResponse{}, nil
}

func (x *DeploymentService) ListReplicaSet(ctx context.Context, req *pb.DeploymentDetailRequest) (*pb.ListReplicaSetResponse, error) {
	params := &biz.DeploymentCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	}
	deployment, err := x.uc.GetDeployment(ctx, params)
	if err != nil {
		return nil, err
	}
	currentRevision := ""
	if deployment.Annotations != nil {
		currentRevision = deployment.Annotations["deployment.kubernetes.io/revision"]
	}
	replicaSets, err := x.uc.GetDeploymentHistory(ctx, params)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ReplicaSetItem, 0, len(replicaSets))
	for _, rs := range replicaSets {
		item := x.ConvertReplicaSet2DTO(rs)
		rs.TypeMeta.Kind = "ReplicaSet"
		rs.TypeMeta.APIVersion = "apps/v1"
		yaml, err := utils.ConvertResourceTOYaml(rs)
		if err != nil {
			continue
		}
		item.Yaml = yaml
		// Check if this ReplicaSet is current based on multiple conditions
		if rs.Annotations != nil {
			rsRevision := rs.Annotations["deployment.kubernetes.io/revision"]
			item.IsCurrent = rsRevision == currentRevision
		}
		list = append(list, item)
	}

	return &pb.ListReplicaSetResponse{
		List:  list,
		Total: uint32(len(list)),
	}, nil
}

func (x *DeploymentService) UpdateDeploymentStrategy(ctx context.Context, req *pb.UpdateDeploymentStrategyRequest) (*pb.UpdateDeploymentStrategyResponse, error) {
	success, err := x.uc.UpdateDeploymentStrategy(ctx, &biz.UpdateDeploymentStrategyRequest{
		DeploymentCommonParams: biz.DeploymentCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		MaxSurge:           req.MaxSurge,
		MaxUnavailable:     req.MaxUnavailable,
		UpdateStrategyType: uint32(req.UpdateStrategyType),
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateDeploymentStrategyResponse{
		Success: success,
	}, nil
}
