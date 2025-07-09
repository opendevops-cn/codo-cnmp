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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sort"
	"strconv"
)

type StatefulSetService struct {
	pb.UnimplementedStatefulSetServer
	uc *biz.StatefulSetUseCase
	uf *biz.UserFollowUseCase
}

func NewStatefulSetService(uc *biz.StatefulSetUseCase, uf *biz.UserFollowUseCase) *StatefulSetService {
	return &StatefulSetService{
		uc: uc,
		uf: uf,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *StatefulSetService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_StatefulSet,
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
func (x *StatefulSetService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.StatefulSetItem) error {
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

func (x *StatefulSetService) convertStatus(status appsv1.StatefulSetStatus) *pb.StatefulSetStatus {
	return &pb.StatefulSetStatus{
		Replicas:            uint32(status.Replicas),
		UpdatedReplicas:     uint32(status.UpdatedReplicas),
		AvailableReplicas:   uint32(status.AvailableReplicas),
		ReadyReplicas:       uint32(status.ReadyReplicas),
		UnavailableReplicas: uint32(status.Replicas - status.AvailableReplicas),
	}
}

func (x *StatefulSetService) convertContainers(containers []corev1.Container) []*corev1.Container {
	result := make([]*corev1.Container, 0, len(containers))
	for _, container := range containers {
		result = append(result, &container)
	}
	return result
}

func (x *StatefulSetService) convertConditions(conditions []appsv1.StatefulSetCondition) []*pb.StatefulSetCondition {
	if len(conditions) == 0 || conditions == nil {
		return []*pb.StatefulSetCondition{}
	}
	result := make([]*pb.StatefulSetCondition, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, &pb.StatefulSetCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: uint64(condition.LastTransitionTime.UnixNano() / 1e6),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return result
}

func (x *StatefulSetService) convertStatefulSetStrategy(strategy appsv1.StatefulSetUpdateStrategy) *pb.StatefulSetStrategy {
	if &strategy == nil {
		return &pb.StatefulSetStrategy{}
	}
	switch strategy.Type {
	case appsv1.OnDeleteStatefulSetStrategyType:
		return &pb.StatefulSetStrategy{
			UpdateStrategyType: pb.StatefulSetStrategy_OnDelete,
		}
	case appsv1.RollingUpdateStatefulSetStrategyType:
		var partition uint32
		if strategy.RollingUpdate != nil {
			partition = uint32(*strategy.RollingUpdate.Partition)
		}
		return &pb.StatefulSetStrategy{
			UpdateStrategyType: pb.StatefulSetStrategy_RollingUpdate,
			Partition:          partition,
		}
	default:
		return &pb.StatefulSetStrategy{}
	}

}

func (x *StatefulSetService) convertDO2DTO(statefulSet *appsv1.StatefulSet) *pb.StatefulSetItem {
	if statefulSet.APIVersion == "" {
		statefulSet.APIVersion = "apps/v1"
	}
	if statefulSet.Kind == "" {
		statefulSet.Kind = "StatefulSet"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(statefulSet)
	updateStrategy := make(map[string]string)
	updateStrategy["type"] = string(statefulSet.Spec.UpdateStrategy.Type)

	return &pb.StatefulSetItem{
		Name:           statefulSet.Name,
		Namespace:      statefulSet.Namespace,
		Status:         x.convertStatus(statefulSet.Status),
		Containers:     x.convertContainers(statefulSet.Spec.Template.Spec.Containers),
		CreateTime:     uint64(statefulSet.CreationTimestamp.UnixNano() / 1e6),
		Labels:         statefulSet.Labels,
		Annotations:    statefulSet.Annotations,
		Selector:       statefulSet.Spec.Selector,
		UpdateStrategy: x.convertStatefulSetStrategy(statefulSet.Spec.UpdateStrategy),
		Replicas:       uint32(statefulSet.Status.Replicas),
		SpecReplicas:   uint32(*statefulSet.Spec.Replicas),
		Conditions:     x.convertConditions(statefulSet.Status.Conditions),
		Yaml:           yamlStr,
	}
}

func (x *StatefulSetService) GetStatefulSetDetail(ctx context.Context, req *pb.GetStatefulSetDetailRequest) (*pb.GetStatefulSetDetailResponse, error) {
	statefulSet, err := x.uc.GetStatefulSet(ctx, &biz.GetStatefulSetRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetStatefulSetDetailResponse{
		Detail: x.convertDO2DTO(statefulSet),
	}, nil
}

func (x *StatefulSetService) ListStatefulSet(ctx context.Context, req *pb.ListStatefulSetRequest) (*pb.ListStatefulSetResponse, error) {
	statefulSets, total, err := x.uc.ListStatefulSet(ctx, &biz.ListStatefulSetRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Keyword:     req.Keyword,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if errors.IsNotFound(err) {
		return &pb.ListStatefulSetResponse{
			List:  []*pb.StatefulSetItem{},
			Total: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	list := make([]*pb.StatefulSetItem, 0, len(statefulSets))
	for _, item := range statefulSets {
		list = append(list, x.convertDO2DTO(item))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListStatefulSetResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *StatefulSetService) RestartStatefulSet(ctx context.Context, req *pb.RestartStatefulSetRequest) (*pb.RestartStatefulSetResponse, error) {
	err := x.uc.RestartStatefulSet(ctx, &biz.RestartStatefulSetRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.RestartStatefulSetResponse{}, nil
}

func (x *StatefulSetService) ScaleStatefulSet(ctx context.Context, req *pb.ScaleStatefulSetRequest) (*pb.ScaleStatefulSetResponse, error) {
	err := x.uc.ScaleStatefulSet(ctx, &biz.ScaleStatefulSetRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Replicas: req.Replicas,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ScaleStatefulSetResponse{}, nil
}

func (x *StatefulSetService) DeleteStatefulSet(ctx context.Context, req *pb.DeleteStatefulSetRequest) (*pb.DeleteStatefulSetResponse, error) {
	err := x.uc.DeleteStatefulSet(ctx, &biz.DeleteStatefulSetRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteStatefulSetResponse{}, nil
}

func (x *StatefulSetService) GetStatefulSetRevisions(ctx context.Context, req *pb.GetStatefulSetHistoryRequest) (*pb.GetStatefulSetHistoryResponse, error) {
	revisions, err := x.uc.GetStatefulSetRevisions(ctx, &biz.StatefulSetCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}
	currentStatefulSet, err := x.uc.GetStatefulSet(ctx, &biz.GetStatefulSetRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.StatefulSetControllerReversionItem, 0, len(revisions))
	for _, revision := range revisions {
		statefulSet, err := DecodeControllerRevisionData2StatefulSet(revision)
		if err != nil {
			continue
		}
		statefulSet.ObjectMeta = revision.ObjectMeta
		images := make([]string, 0, len(statefulSet.Spec.Template.Spec.Containers))
		for _, container := range statefulSet.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}

		if statefulSet.APIVersion == "" {
			statefulSet.APIVersion = "apps/v1"
		}
		if statefulSet.Kind == "" {
			statefulSet.Kind = "StatefulSet"
		}
		yamlStr, err := utils.ConvertResourceTOYaml(statefulSet)
		if err != nil {
			continue
		}
		isCurrent := revision.Name == currentStatefulSet.Status.CurrentRevision

		list = append(list, &pb.StatefulSetControllerReversionItem{
			Name:       revision.Name,
			CreateTime: uint64(revision.CreationTimestamp.UnixNano() / 1e6),
			Images:     images,
			Version:    strconv.FormatInt(revision.Revision, 10),
			Yaml:       yamlStr,
			IsCurrent:  isCurrent,
		})
	}

	return &pb.GetStatefulSetHistoryResponse{
		List:  list,
		Total: uint32(len(list)),
	}, nil
}

// DecodeControllerRevisionData2StatefulSet 解析ControllerRevision的数据
func DecodeControllerRevisionData2StatefulSet(controllerRevision *appsv1.ControllerRevision) (*appsv1.StatefulSet, error) {
	statefulSet := &appsv1.StatefulSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(controllerRevision.Data.Raw), 4096)
	if err := decoder.Decode(&statefulSet); err != nil {
		return nil, fmt.Errorf("解析失败, 请检查格式是否正确: %v", err)
	}
	return statefulSet, nil
}

// RollbackStatefulSet 回滚StatefulSet
func (x *StatefulSetService) RollbackStatefulSet(ctx context.Context, req *pb.RollbackStatefulSetRequest) (*pb.RollbackStatefulSetResponse, error) {
	version, err := strconv.ParseUint(req.Version, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的版本号: %v", err)
	}
	err = x.uc.RollbackStatefulSet(ctx, &biz.RollbackStatefulSetRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Revision: uint32(version),
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.RollbackStatefulSetResponse{}, nil
}

func (x *StatefulSetService) CreateOrUpdateStatefulSetByYaml(ctx context.Context, req *pb.CreateOrUpdateStatefulSetByYamlRequest) (*pb.CreateOrUpdateStatefulSetByYamlResponse, error) {
	err := x.uc.CreateOrUpdateStatefulSetByYaml(ctx, &biz.CreateOrUpdateStatefulSetByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateStatefulSetByYamlResponse{}, nil
}

func (x *StatefulSetService) UpdateStatefulSetUpdateStrategy(ctx context.Context, req *pb.UpdateStatefulSetUpdateStrategyRequest) (*pb.UpdateStatefulSetUpdateStrategyResponse, error) {
	_, err := x.uc.UpdateStatefulSetStrategy(ctx, &biz.UpdateStatefulSetStrategyRequest{
		StatefulSetCommonParams: biz.StatefulSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		UpdateStrategyType: uint32(req.UpdateStrategyType),
		Partition:          req.Partition,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateStatefulSetUpdateStrategyResponse{}, nil
}
