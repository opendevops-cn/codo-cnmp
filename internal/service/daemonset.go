package service

import (
	"bytes"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sort"
	"strconv"
)

type DaemonSetService struct {
	pb.UnimplementedDaemonSetServer
	uc *biz.DaemonSetUseCase
	uf *biz.UserFollowUseCase
}

func NewDaemonSetService(uc *biz.DaemonSetUseCase, uf *biz.UserFollowUseCase) *DaemonSetService {
	return &DaemonSetService{uc: uc, uf: uf}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *DaemonSetService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_DaemonSet,
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
func (x *DaemonSetService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.DaemonSetItem) error {
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

func (x *DaemonSetService) CreateOrUpdateDaemonSetByYaml(ctx context.Context, request *pb.CreateOrUpdateDaemonSetByYamlRequest) (*pb.CreateOrUpdateDaemonSetByYamlResponse, error) {
	err := x.uc.CreateOrUpdateDaemonSetByYaml(
		ctx, &biz.CreateOrUpdateDaemonSetByYamlRequest{
			ClusterName: request.ClusterName,
			Yaml:        request.Yaml,
		})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateDaemonSetByYamlResponse{}, nil
}

func (x *DaemonSetService) DeleteDaemonSet(ctx context.Context, request *pb.DeleteDaemonSetRequest) (*pb.DeleteDaemonSetResponse, error) {
	err := x.uc.DeleteDaemonSet(ctx, &biz.DaemonSetCommonParams{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
		Name:        request.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteDaemonSetResponse{}, nil
}

func (x *DaemonSetService) RestartDaemonSet(ctx context.Context, request *pb.RestartDaemonSetRequest) (*pb.RestartDaemonSetResponse, error) {
	err := x.uc.RestartDaemonSet(ctx, &biz.DaemonSetCommonParams{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
		Name:        request.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RestartDaemonSetResponse{}, nil
}

func (x *DaemonSetService) GetDaemonSetDetail(ctx context.Context, request *pb.GetDaemonSetDetailRequest) (*pb.GetDaemonSetDetailResponse, error) {
	daemonSet, err := x.uc.GetDaemonSet(ctx, &biz.DaemonSetCommonParams{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
		Name:        request.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetDaemonSetDetailResponse{
		Detail: x.convertDO2DTO(daemonSet),
	}, nil
}

func (x *DaemonSetService) GetDaemonSetRevisions(ctx context.Context, request *pb.GetDaemonSetHistoryRequest) (*pb.GetDaemonSetHistoryResponse, error) {
	revisions, err := x.uc.GetDaemonSetHistory(ctx, &biz.DaemonSetCommonParams{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
		Name:        request.Name,
	})
	if err != nil {
		return nil, err
	}
	currentDaemonSet, err := x.uc.GetDaemonSet(ctx, &biz.DaemonSetCommonParams{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
		Name:        request.Name,
	})

	if err != nil {
		return nil, err
	}
	list := make([]*pb.DaemonSetControllerReversionItem, 0, len(revisions))
	for _, revision := range revisions {
		daemonSet, err := DecodeDaemonSetRevisionData2DaemonSet(revision)
		if err != nil {
			continue
		}
		daemonSet.ObjectMeta = revision.ObjectMeta
		images := make([]string, 0, len(daemonSet.Spec.Template.Spec.Containers))
		for _, container := range daemonSet.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}

		if daemonSet.APIVersion == "" {
			daemonSet.APIVersion = "apps/v1"
		}
		if daemonSet.Kind == "" {
			daemonSet.Kind = "DaemonSet"
		}
		yamlStr, err := utils.ConvertResourceTOYaml(daemonSet)
		if err != nil {
			continue
		}
		var IsCurrent bool
		if revision.OwnerReferences == nil {
			IsCurrent = false
		} else if len(revision.OwnerReferences) == 0 {
			IsCurrent = true
		} else {
			if currentDaemonSet.UID == revision.OwnerReferences[0].UID {
				IsCurrent = true
			}
		}

		list = append(list, &pb.DaemonSetControllerReversionItem{
			Name:       revision.Name,
			CreateTime: uint64(revision.CreationTimestamp.UnixNano() / 1e6),
			Images:     images,
			Version:    strconv.FormatInt(revision.Revision, 10),
			Yaml:       yamlStr,
			IsCurrent:  IsCurrent,
		})
	}

	return &pb.GetDaemonSetHistoryResponse{
		List:  list,
		Total: uint32(len(list)),
	}, nil

}

func (x *DaemonSetService) RollbackDaemonSet(ctx context.Context, request *pb.RollbackDaemonSetRequest) (*pb.RollbackDaemonSetResponse, error) {
	version, err := strconv.ParseUint(request.Version, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %v", err)
	}
	err = x.uc.RollbackDaemonSet(ctx, &biz.RollbackDaemonSetRequest{
		DaemonSetCommonParams: biz.DaemonSetCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name,
		},
		Revision: uint32(version),
	})
	if err != nil {
		return nil, err
	}
	return &pb.RollbackDaemonSetResponse{}, nil
}

func (x *DaemonSetService) UpdateStatefulSetUpdateStrategy(ctx context.Context, request *pb.UpdateDaemonSetUpdateStrategyRequest) (*pb.UpdateDaemonSetUpdateStrategyResponse, error) {
	_, err := x.uc.UpdateDaemonSetStrategy(ctx, &biz.UpdateDaemonSetStrategyRequest{
		DaemonSetCommonParams: biz.DaemonSetCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name,
		},
		MaxSurge:       request.MaxSurge,
		MaxUnavailable: request.MaxUnavailable,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateDaemonSetUpdateStrategyResponse{}, nil
}

func (x *DaemonSetService) convertStatus(status appsv1.DaemonSetStatus) *pb.DaemonSetStatus {
	return &pb.DaemonSetStatus{
		Replicas:            uint32(status.CurrentNumberScheduled),
		UpdatedReplicas:     uint32(status.UpdatedNumberScheduled),
		AvailableReplicas:   uint32(status.NumberAvailable),
		ReadyReplicas:       uint32(status.NumberReady),
		UnavailableReplicas: uint32(status.NumberUnavailable),
	}
}

func (x *DaemonSetService) convertContainers(containers []corev1.Container) []*corev1.Container {
	result := make([]*corev1.Container, 0, len(containers))
	for _, container := range containers {
		result = append(result, &container)
	}
	return result
}

func (x *DaemonSetService) convertConditions(conditions []appsv1.DaemonSetCondition) []*pb.DaemonSetCondition {
	if len(conditions) == 0 || conditions == nil {
		return []*pb.DaemonSetCondition{}
	}
	result := make([]*pb.DaemonSetCondition, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, &pb.DaemonSetCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: uint64(condition.LastTransitionTime.UnixNano() / 1e6),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return result
}

func (x *DaemonSetService) convertDaemonSetStrategy(strategy *appsv1.DaemonSetUpdateStrategy) *pb.DaemonSetStrategy {
	switch strategy.Type {
	case appsv1.RollingUpdateDaemonSetStrategyType:
		return &pb.DaemonSetStrategy{
			UpdateStrategyType: pb.DaemonSetStrategy_RollingUpdate,
			MaxUnavailable:     strategy.RollingUpdate.MaxUnavailable.String(),
		}
	case appsv1.OnDeleteDaemonSetStrategyType:
		return &pb.DaemonSetStrategy{
			UpdateStrategyType: pb.DaemonSetStrategy_OnDelete,
		}
	}
	return &pb.DaemonSetStrategy{}

}

func (x *DaemonSetService) convertDO2DTO(daemonSet *appsv1.DaemonSet) *pb.DaemonSetItem {
	if daemonSet.APIVersion == "" {
		daemonSet.APIVersion = "apps/v1"
	}
	if daemonSet.Kind == "" {
		daemonSet.Kind = "DaemonSet"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(daemonSet)
	updateStrategy := make(map[string]string)
	updateStrategy["type"] = string(daemonSet.Spec.UpdateStrategy.Type)

	return &pb.DaemonSetItem{
		Name:           daemonSet.Name,
		Namespace:      daemonSet.Namespace,
		Status:         x.convertStatus(daemonSet.Status),
		Containers:     x.convertContainers(daemonSet.Spec.Template.Spec.Containers),
		CreateTime:     uint64(daemonSet.CreationTimestamp.UnixNano() / 1e6),
		Labels:         daemonSet.Labels,
		Annotations:    daemonSet.Annotations,
		Selector:       daemonSet.Spec.Selector,
		UpdateStrategy: x.convertDaemonSetStrategy(&daemonSet.Spec.UpdateStrategy),
		Replicas:       uint32(daemonSet.Status.CurrentNumberScheduled),
		SpecReplicas:   uint32(daemonSet.Status.DesiredNumberScheduled),
		Conditions:     x.convertConditions(daemonSet.Status.Conditions),
		Yaml:           yamlStr,
	}
}

func (x *DaemonSetService) ListDaemonSet(ctx context.Context, req *pb.ListDaemonSetRequest) (*pb.ListDaemonSetResponse, error) {
	daemonSets, total, err := x.uc.ListDaemonSet(ctx, &biz.ListDaemonSetRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Page:        req.Page,
		PageSize:    req.PageSize,
		Keyword:     req.Keyword,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.DaemonSetItem, 0, len(daemonSets))
	for _, daemonSet := range daemonSets {
		dto := x.convertDO2DTO(daemonSet)
		list = append(list, dto)
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListDaemonSetResponse{
		Total: total,
		List:  list,
	}, nil

}

// DecodeDaemonSetRevisionData2DaemonSet 解析ControllerRevision的数据
func DecodeDaemonSetRevisionData2DaemonSet(controllerRevision *appsv1.ControllerRevision) (*appsv1.DaemonSet, error) {
	daemonSet := &appsv1.DaemonSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(controllerRevision.Data.Raw), 4096)
	if err := decoder.Decode(&daemonSet); err != nil {
		return nil, fmt.Errorf("解析失败, 请检查格式是否正确: %v", err)
	}
	return daemonSet, nil
}
