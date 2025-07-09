package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	kruiseappsv1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sort"
)

type SideCarSetService struct {
	pb.UnimplementedSidecarSetServer
	uc *biz.SideCarSetUseCase
	uf *biz.UserFollowUseCase
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *SideCarSetService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_SideCarSet,
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
func (x *SideCarSetService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.SidecarSetItem) error {
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

func (x *SideCarSetService) DeleteSidecarSet(ctx context.Context, request *pb.DeleteSidecarSetRequest) (*pb.DeleteSidecarSetResponse, error) {
	result, err := x.uc.DeleteSideCarSet(ctx, &biz.DeleteSideCarSetRequest{
		ClusterName: request.GetClusterName(),
		Name:        request.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteSidecarSetResponse{
		Success: result,
	}, nil
}

func (x *SideCarSetService) UpdateSideCarSetStrategy(ctx context.Context, request *pb.UpdateSideCarSetStrategyRequest) (*pb.UpdateSideCarSetStrategyResponse, error) {
	result, err := x.uc.UpdateSideCarSetStrategy(ctx, &biz.UpdateSideCarSetStrategyRequest{
		ClusterName:        request.GetClusterName(),
		MaxUnavailable:     request.GetMaxUnavailable(),
		Name:               request.GetName(),
		Partition:          request.GetPartition(),
		Paused:             request.GetPause(),
		UpdateStrategyType: uint32(request.GetUpdateStrategyType()),
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSideCarSetStrategyResponse{
		Success: result,
	}, nil
}

func (x *SideCarSetService) ListSidecarSet(ctx context.Context, request *pb.ListSidecarSetRequest) (*pb.ListSidecarSetResponse, error) {
	sideCarSets, total, err := x.uc.ListSideCarSet(ctx, &biz.ListSideCarSetRequest{
		ClusterName: request.GetClusterName(),
		Keyword:     request.GetKeyword(),
		Page:        request.GetPage(),
		PageSize:    request.GetPageSize(),
		ListAll:     utils.IntToBool(int(request.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.SidecarSetItem, 0, len(sideCarSets))
	for _, sideCarSet := range sideCarSets {
		list = append(list, x.convertDO2DTO(sideCarSet))
	}
	if err := x.setFollowedStatus(ctx, request.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListSidecarSetResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *SideCarSetService) convertStatus(status kruiseappsv1alpha1.SidecarSetStatus) *pb.SideCarSetStatus {
	return &pb.SideCarSetStatus{
		MatchedPods:      uint32(status.MatchedPods),
		ReadyPods:        uint32(status.ReadyPods),
		UpdatedPods:      uint32(status.UpdatedPods),
		UpdatedReadyPods: uint32(status.UpdatedReadyPods),
	}
}

func (x *SideCarSetService) convertDO2DTO(sidecarSet *kruiseappsv1alpha1.SidecarSet) *pb.SidecarSetItem {
	if sidecarSet == nil {
		return &pb.SidecarSetItem{}
	}
	status := sidecarSet.Status
	if &status == nil {
		return &pb.SidecarSetItem{}
	}
	if sidecarSet.Kind == "" {
		sidecarSet.Kind = "SidecarSet"
	}
	if sidecarSet.APIVersion == "" {
		sidecarSet.APIVersion = "apps.kruise.io/v1alpha1"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(sidecarSet)

	containers := make([]*corev1.Container, 0, len(sidecarSet.Spec.Containers))
	for _, container := range sidecarSet.Spec.Containers {
		containers = append(containers, &container.Container)
	}

	return &pb.SidecarSetItem{
		Name:             sidecarSet.GetName(),
		ReadyPods:        uint32(status.ReadyPods),
		MatchedPods:      uint32(status.MatchedPods),
		UpdatedPods:      uint32(status.UpdatedPods),
		UpdatedReadyPods: uint32(status.UpdatedReadyPods),
		CreateTime:       uint64(sidecarSet.CreationTimestamp.Time.UnixNano() / 1e6),
		Selector:         sidecarSet.Spec.Selector,
		Labels:           sidecarSet.GetLabels(),
		Annotations:      sidecarSet.GetAnnotations(),
		Containers:       containers,
		Yaml:             yamlStr,
		Status:           x.convertStatus(status),
		UpdateStrategy:   x.convertUpdateStrategy(sidecarSet.Spec.UpdateStrategy),
	}
}

func (x *SideCarSetService) GetSidecarSet(ctx context.Context, request *pb.GetSidecarSetRequest) (*pb.GetSidecarSetResponse, error) {
	sideCarSet, err := x.uc.GetSideCarSet(ctx, &biz.GetSideCarSetRequest{
		ClusterName: request.GetClusterName(),
		Name:        request.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetSidecarSetResponse{
		Detail: x.convertDO2DTO(sideCarSet),
	}, nil
}

func (x *SideCarSetService) convertUpdateStrategy(strategy kruiseappsv1alpha1.SidecarSetUpdateStrategy) *pb.UpdateSideCarSetStrategyRequest {
	if strategy.Type == kruiseappsv1alpha1.NotUpdateSidecarSetStrategyType {
		return &pb.UpdateSideCarSetStrategyRequest{
			UpdateStrategyType: pb.UpdateSideCarSetStrategyRequest_NotUpdate,
		}
	} else if strategy.Type == kruiseappsv1alpha1.RollingUpdateSidecarSetStrategyType {
		return &pb.UpdateSideCarSetStrategyRequest{
			UpdateStrategyType: pb.UpdateSideCarSetStrategyRequest_RollingUpdate,
			Pause:              strategy.Paused,
			Partition:          strategy.Partition.String(),
			MaxUnavailable:     strategy.MaxUnavailable.String(),
		}
	} else {
		return &pb.UpdateSideCarSetStrategyRequest{}
	}
}

func NewSideCarSetService(uc *biz.SideCarSetUseCase, uf *biz.UserFollowUseCase) *SideCarSetService {
	return &SideCarSetService{
		uc: uc,
		uf: uf,
	}
}
