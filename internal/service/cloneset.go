package service

import (
	"bytes"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	kruiseappsv1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sort"
	"strconv"
)

type CloneSetService struct {
	pb.UnimplementedCloneSetServer
	uc *biz.CloneSetUseCase
	uf *biz.UserFollowUseCase
}

func NewCloneSetService(uc *biz.CloneSetUseCase, uf *biz.UserFollowUseCase) *CloneSetService {
	return &CloneSetService{
		uc: uc,
		uf: uf,
	}
}

// convertContainers converts the containers of deployment to the pb.Container type
func (x *CloneSetService) convertContainers(containers []corev1.Container) []*corev1.Container {
	result := make([]*corev1.Container, 0, len(containers))
	for _, container := range containers {
		result = append(result, &container)
	}
	return result
}

// convertContainers converts the status of deployment to the pb.DeploymentStatus type
func (x *CloneSetService) convertStatus(status kruiseappsv1alpha1.CloneSetStatus) *pb.CloneSetStatus {
	return &pb.CloneSetStatus{
		Replicas:            uint32(status.Replicas),
		UpdatedReplicas:     uint32(status.UpdatedReplicas),
		AvailableReplicas:   uint32(status.AvailableReplicas),
		ReadyReplicas:       uint32(status.ReadyReplicas),
		UnavailableReplicas: uint32(status.Replicas - status.AvailableReplicas),
	}
}

// convertConditions converts the conditions of deployment to the pb.DeploymentCondition type
func (x *CloneSetService) convertConditions(conditions []kruiseappsv1alpha1.CloneSetCondition) []*pb.CloneSetCondition {
	if len(conditions) == 0 || conditions == nil {
		return []*pb.CloneSetCondition{}
	}
	result := make([]*pb.CloneSetCondition, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, &pb.CloneSetCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: uint64(condition.LastTransitionTime.Time.UnixNano() / 1e6),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return result
}

func (x *CloneSetService) convertCloneSetStrategy(strategy kruiseappsv1alpha1.CloneSetUpdateStrategy) *pb.CloneSetUpdateStrategy {
	switch strategy.Type {
	case kruiseappsv1alpha1.RecreateCloneSetUpdateStrategyType:
		return &pb.CloneSetUpdateStrategy{
			UpdateStrategyType: pb.CloneSetUpdateStrategy_Recreate,
		}
	case kruiseappsv1alpha1.InPlaceIfPossibleCloneSetUpdateStrategyType:
		var gracePeriodSeconds uint32
		if strategy.InPlaceUpdateStrategy != nil {
			gracePeriodSeconds = uint32(strategy.InPlaceUpdateStrategy.GracePeriodSeconds)
		}
		return &pb.CloneSetUpdateStrategy{
			UpdateStrategyType: pb.CloneSetUpdateStrategy_InPlaceIfPossible,
			GracePeriodSeconds: gracePeriodSeconds,
			MaxSurge:           strategy.MaxSurge.String(),
			MaxUnavailable:     strategy.MaxUnavailable.String(),
		}
	case kruiseappsv1alpha1.InPlaceOnlyCloneSetUpdateStrategyType:
		return &pb.CloneSetUpdateStrategy{
			UpdateStrategyType: pb.CloneSetUpdateStrategy_InPlaceOnly,
			MaxUnavailable:     "",
		}
	default:
		return &pb.CloneSetUpdateStrategy{}
	}
}

func (x *CloneSetService) convertScaleStrategy(cloneSet *kruiseappsv1alpha1.CloneSet) *pb.ScaleStrategyItem {
	if &cloneSet.Spec.ScaleStrategy == nil {
		return &pb.ScaleStrategyItem{}
	}
	var (
		maxUnavailable string
	)
	if cloneSet.Spec.ScaleStrategy.MaxUnavailable != nil {
		maxUnavailable = cloneSet.Spec.ScaleStrategy.MaxUnavailable.String()
	}
	return &pb.ScaleStrategyItem{
		MinReadySeconds: fmt.Sprintf("%d", cloneSet.Spec.MinReadySeconds),
		MaxUnavailable:  maxUnavailable,
	}

}

// convertDO2DTO converts the cloneSet object to the pb.CloneSetItem type
func (x *CloneSetService) convertDO2DTO(cloneSet *kruiseappsv1alpha1.CloneSet) *pb.CloneSetItem {
	if cloneSet.APIVersion == "" {
		cloneSet.APIVersion = "apps.kruise.io/v1alpha1"
	}
	if cloneSet.Kind == "" {
		cloneSet.Kind = "CloneSet"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(cloneSet)
	return &pb.CloneSetItem{
		Name:           cloneSet.Name,
		Namespace:      cloneSet.Namespace,
		Status:         x.convertStatus(cloneSet.Status),
		Containers:     x.convertContainers(cloneSet.Spec.Template.Spec.Containers),
		CreateTime:     uint64(cloneSet.CreationTimestamp.UnixNano() / 1e6),
		Labels:         cloneSet.Labels,
		Annotations:    cloneSet.Annotations,
		Selector:       cloneSet.Spec.Selector,
		UpdateStrategy: x.convertCloneSetStrategy(cloneSet.Spec.UpdateStrategy),
		Replicas:       uint32(cloneSet.Status.Replicas),
		SpecReplicas:   uint32(*cloneSet.Spec.Replicas),
		Conditions:     x.convertConditions(cloneSet.Status.Conditions),
		Yaml:           yamlStr,
		ScaleStrategy:  x.convertScaleStrategy(cloneSet),
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *CloneSetService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_CloneSet,
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
func (x *CloneSetService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.CloneSetItem) error {
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

func (x *CloneSetService) ListCloneSet(ctx context.Context, req *pb.ListCloneSetRequest) (*pb.ListCloneSetResponse, error) {
	cloneSets, total, err := x.uc.ListCloneSets(ctx, &biz.LisCloneSetRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
		Keyword:  req.Keyword,
	})
	if errors.IsNotFound(err) {
		return &pb.ListCloneSetResponse{
			List:  []*pb.CloneSetItem{},
			Total: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	list := make([]*pb.CloneSetItem, 0, len(cloneSets))
	for _, cloneSet := range cloneSets {
		yamlStr, err := utils.ConvertResourceTOYaml(cloneSet)
		if err != nil {
			fmt.Println(err)
		}
		dto := x.convertDO2DTO(cloneSet)
		dto.Yaml = yamlStr
		list = append(list, dto)
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListCloneSetResponse{
		List:  list,
		Total: int32(total),
	}, nil
}

func (x *CloneSetService) CreateOrUpdateCloneSetByYaml(ctx context.Context, req *pb.CreateOrUpdateCloneSetByYamlRequest) (*pb.CreateOrUpdateCloneSetByYamlResponse, error) {
	err := x.uc.CreateOrUpdateCloneSetByYaml(ctx, &biz.CreateOrUpdateCloneSetByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateCloneSetByYamlResponse{}, nil
}

// DeleteCloneSet deletes the cloneSet by name
func (x *CloneSetService) DeleteCloneSet(ctx context.Context, req *pb.DeleteCloneSetRequest) (*pb.DeleteCloneSetResponse, error) {
	err := x.uc.DeleteCloneSet(ctx, &biz.DeleteCloneSetRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		}})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCloneSetResponse{}, nil
}

// RestartCloneSet restarts the cloneSet by name
func (x *CloneSetService) RestartCloneSet(ctx context.Context, req *pb.RestartCloneSetRequest) (*pb.RestartCloneSetResponse, error) {
	err := x.uc.RestartCloneSet(ctx, &biz.RestartCloneSetRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		}})
	if err != nil {
		return nil, err
	}
	return &pb.RestartCloneSetResponse{}, nil
}

// ScaleCloneSet scales the cloneSet by name
func (x *CloneSetService) ScaleCloneSet(ctx context.Context, req *pb.ScaleCloneSetRequest) (*pb.ScaleCloneSetResponse, error) {
	err := x.uc.ScaleCloneSet(ctx, &biz.ScaleCloneSetRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Replicas: req.Replicas,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ScaleCloneSetResponse{}, nil
}

// RollbackCloneSet rolls back the cloneSet by name
func (x *CloneSetService) RollbackCloneSet(ctx context.Context, req *pb.RollbackCloneSetRequest) (*pb.RollbackCloneSetResponse, error) {
	version, err := strconv.ParseUint(req.Version, 10, 32)
	if err != nil {
		return nil, err
	}
	err = x.uc.RollbackCloneSet(ctx, &biz.RollbackCloneSetRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Revision: uint32(version),
	})
	if err != nil {
		return nil, err
	}
	return &pb.RollbackCloneSetResponse{}, nil
}

// GetCloneSetDetail gets the cloneSet detail by name
func (x *CloneSetService) GetCloneSetDetail(ctx context.Context, req *pb.CloneSetDetailRequest) (*pb.CloneSetDetailResponse, error) {
	cloneSet, err := x.uc.GetCloneSetDetail(ctx, &biz.GetCloneSetDetailRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name}})
	if err != nil {
		return nil, err
	}
	yamlStr, err := utils.ConvertResourceTOYaml(cloneSet)
	if err != nil {
		fmt.Println(err)
	}
	dto := x.convertDO2DTO(cloneSet)
	dto.Yaml = yamlStr
	return &pb.CloneSetDetailResponse{
		Detail: dto,
	}, nil
}

// DeleteCloneSetPods deletes the cloneSet pod by name
func (x *CloneSetService) DeleteCloneSetPods(ctx context.Context, req *pb.DeleteCloneSetPodRequest) (*pb.DeleteCloneSetPodResponse, error) {
	err := x.uc.DeleteCloneSetPods(ctx, &biz.DeleteCloneSetPodsRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		PodNames:     req.PodNames,
		DeletePolicy: uint32(req.DeletePolicy),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCloneSetPodResponse{}, nil
}

// UpdateScaleStrategy updates the scale strategy of the cloneSet by name
func (x *CloneSetService) UpdateScaleStrategy(ctx context.Context, req *pb.UpdateScaleStrategyRequest) (*pb.UpdateScaleStrategyResponse, error) {
	err := x.uc.UpdateScaleStrategy(ctx, &biz.UpdateScaleStrategyRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		MaxUnavailable:  req.MaxUnavailable,
		MinReadySeconds: req.MinReadySeconds,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateScaleStrategyResponse{}, nil
}

// UpdateUpgradeStrategy updates the upgrade strategy of the cloneSet by name
func (x *CloneSetService) UpdateUpgradeStrategy(ctx context.Context, req *pb.UpdateUpgradeStrategyRequest) (*pb.UpdateUpgradeStrategyResponse, error) {
	err := x.uc.UpdateUpgradeStrategy(ctx, &biz.UpdateUpgradeStrategyRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		UpdateStrategyType: uint32(req.UpdateStrategyType),
		GracePeriodSeconds: req.GracePeriodSeconds,
		MaxSurge:           req.MaxSurge,
		MaxUnavailable:     req.MaxUnavailable,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUpgradeStrategyResponse{}, nil
}

// ListCloneSetControllerRevision lists the cloneSet controller revisions by name
func (x *CloneSetService) ListCloneSetControllerRevision(ctx context.Context, req *pb.ListCloneSetReversionRequest) (*pb.ListCloneSetReversionResponse, error) {
	revisions, err := x.uc.ListCloneSetControllerRevisions(ctx, &biz.CloneSetCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}
	currentCloneSet, err := x.uc.GetCloneSetDetail(ctx, &biz.GetCloneSetDetailRequest{
		CloneSetCommonParams: biz.CloneSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	currentRevision := currentCloneSet.Status.CurrentRevision
	list := make([]*pb.CloneSetControllerReversionItem, 0, len(revisions))
	for _, revision := range revisions {
		cloneSet, err := DecodeControllerRevisionData2CloneSet(revision)
		if err != nil {
			continue
		}
		images := make([]string, 0, len(cloneSet.Spec.Template.Spec.Containers))
		for _, container := range cloneSet.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}

		if cloneSet.APIVersion == "" {
			cloneSet.APIVersion = "apps.kruise.io/v1alpha1"
		}
		if cloneSet.Kind == "" {
			cloneSet.Kind = "CloneSet"
		}
		yamlStr, err := utils.ConvertResourceTOYaml(cloneSet)
		if err != nil {
			continue
		}
		var IsCurrent bool
		if revision.Name == currentRevision {
			IsCurrent = true
		} else {
			IsCurrent = false
		}

		list = append(list, &pb.CloneSetControllerReversionItem{
			Name:       revision.Name,
			CreateTime: uint64(revision.CreationTimestamp.UnixNano() / 1e6),
			Images:     images,
			Version:    strconv.FormatInt(revision.Revision, 10),
			Yaml:       yamlStr,
			IsCurrent:  IsCurrent,
		})
	}
	return &pb.ListCloneSetReversionResponse{
		List:  list,
		Total: uint32(len(revisions)),
	}, nil
}

// DecodeControllerRevisionData2CloneSet 解析ControllerRevision的数据
func DecodeControllerRevisionData2CloneSet(controllerRevision *appsv1.ControllerRevision) (*kruiseappsv1alpha1.CloneSet, error) {
	cloneSet := &kruiseappsv1alpha1.CloneSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(controllerRevision.Data.Raw), 4096)
	if err := decoder.Decode(&cloneSet); err != nil {
		return nil, fmt.Errorf("解析失败, 请检查格式是否正确: %v", err)
	}
	return cloneSet, nil
}
