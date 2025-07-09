package service

import (
	"bytes"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	kruiseappsv1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	"github.com/openkruise/kruise-api/apps/v1beta1"
	gamekruisev1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sort"
	"strconv"
)

type GameServerSetService struct {
	pb.UnimplementedGameServerSetServer
	uc *biz.GameServerSetUseCase
	uf *biz.UserFollowUseCase
}

func (x *GameServerSetService) ScaleGameServerSet(ctx context.Context, req *pb.ScaleGameServerSetRequest) (*pb.ScaleGameServerSetResponse, error) {
	err := x.uc.ScaleGameServerSet(ctx, &biz.ScaleGameServerSetRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Replicas: int32(req.Replicas),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ScaleGameServerSetResponse{}, nil
}

// getUserFollowMap 获取用户关注的GameServerSet列表
func (x *GameServerSetService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_EzRollout,
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
func (x *GameServerSetService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.GameServerSetItem) error {
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

func (x *GameServerSetService) GetGameServerSetDetail(ctx context.Context, req *pb.GameServerSetDetailRequest) (*pb.GameServerSetDetailResponse, error) {
	gameServerSet, err := x.uc.GetGameServerSetDetail(ctx, &biz.GetGameServerSetDetailRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.GameServerSetDetailResponse{Detail: x.convertDO2DTO(gameServerSet)}, nil
}

func (x *GameServerSetService) RollbackGameServerSet(ctx context.Context, req *pb.RollbackGameServerSetRequest) (*pb.RollbackGameServerSetResponse, error) {
	version, err := strconv.ParseUint(req.Version, 10, 32)
	if err != nil {
		return nil, err
	}
	err = x.uc.RollbackGameServerSet(ctx, &biz.RollbackGameServerSetRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Revision: uint32(version),
	})
	if err != nil {
		return nil, err
	}
	return &pb.RollbackGameServerSetResponse{}, nil
}

func (x *GameServerSetService) ListGameServerSetControllerRevision(ctx context.Context, req *pb.ListGameServerSetReversionRequest) (*pb.ListGameServerSetReversionResponse, error) {
	revisions, err := x.uc.ListGameServerSetControllerRevisions(ctx, &biz.GameServerSetCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})

	if err != nil {
		return nil, err
	}
	currentGameServerSet, err := x.uc.GetGameServerSetDetail(ctx, &biz.GetGameServerSetDetailRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}

	list := make([]*pb.GameServerSetControllerReversionItem, 0, len(revisions))
	for _, revision := range revisions {
		gss, err := DecodeControllerRevisionData2GameServerSet(revision)
		if err != nil {
			continue
		}
		gss.ObjectMeta = revision.ObjectMeta
		images := make([]string, 0, len(gss.Spec.GameServerTemplate.Spec.Containers))
		for _, container := range gss.Spec.GameServerTemplate.Spec.Containers {
			images = append(images, container.Image)
		}

		if gss.APIVersion == "" {
			gss.APIVersion = "apps.kruise.io/v1alpha1"
		}
		if gss.Kind == "" {
			gss.Kind = "GameServerSet"
		}
		yamlStr, err := utils.ConvertResourceTOYaml(gss)
		if err != nil {
			continue
		}
		var IsCurrent bool
		if revision.UID == currentGameServerSet.UID {
			IsCurrent = true
		} else {
			IsCurrent = false
		}

		list = append(list, &pb.GameServerSetControllerReversionItem{
			Name:       revision.Name,
			CreateTime: uint64(revision.CreationTimestamp.UnixNano() / 1e6),
			Images:     images,
			Version:    strconv.FormatInt(revision.Revision, 10),
			Yaml:       yamlStr,
			IsCurrent:  IsCurrent,
		})
	}

	return &pb.ListGameServerSetReversionResponse{
		List:  list,
		Total: uint32(len(list)),
	}, nil
}

func NewGameServerSetService(uc *biz.GameServerSetUseCase, uf *biz.UserFollowUseCase) *GameServerSetService {
	return &GameServerSetService{
		uc: uc,
		uf: uf,
	}
}

// convertContainers converts the containers of deployment to the pb.Container type
func (x *GameServerSetService) convertContainers(containers []corev1.Container) []*corev1.Container {
	result := make([]*corev1.Container, 0, len(containers))
	for _, container := range containers {
		result = append(result, &container)
	}
	return result
}

// convertContainers converts the status of deployment to the pb.DeploymentStatus type
func (x *GameServerSetService) convertStatus(status gamekruisev1alpha1.GameServerSetStatus) *pb.GameServerSetStatus {
	return &pb.GameServerSetStatus{
		Replicas:            uint32(status.Replicas),
		UpdatedReplicas:     uint32(status.UpdatedReplicas),
		AvailableReplicas:   uint32(status.AvailableReplicas),
		ReadyReplicas:       uint32(status.ReadyReplicas),
		UnavailableReplicas: uint32(status.Replicas - status.AvailableReplicas),
	}
}

func (x *GameServerSetService) convertGameServerSetStrategy(strategy gamekruisev1alpha1.UpdateStrategy) *pb.GameServerSetStrategy {
	switch strategy.Type {
	case appsv1.RollingUpdateStatefulSetStrategyType:
		// 处理滚动更新策略
		var podUpdateStrategy pb.GameServerSetStrategy_PodUpdateStrategy
		switch strategy.RollingUpdate.PodUpdatePolicy {
		case v1beta1.InPlaceIfPossiblePodUpdateStrategyType:
			podUpdateStrategy = pb.GameServerSetStrategy_InPlaceIfPossible
		case v1beta1.InPlaceOnlyPodUpdateStrategyType:
			podUpdateStrategy = pb.GameServerSetStrategy_InPlaceOnly
		case v1beta1.RecreatePodUpdateStrategyType:
			podUpdateStrategy = pb.GameServerSetStrategy_Recreate
		default:
			// 默认策略
			podUpdateStrategy = pb.GameServerSetStrategy_InPlaceIfPossible
		}

		return &pb.GameServerSetStrategy{
			UpdateStrategyType:    pb.GameServerSetStrategy_RollingUpdate,
			PodUpdateStrategyType: podUpdateStrategy,
			GracePeriodSeconds:    uint32(strategy.RollingUpdate.InPlaceUpdateStrategy.GracePeriodSeconds),
			MaxUnavailable:        strategy.RollingUpdate.MaxUnavailable.String(),
			Partition:             uint32(*strategy.RollingUpdate.Partition),
		}

	case appsv1.OnDeleteStatefulSetStrategyType:
		// 处理删除后更新策略
		return &pb.GameServerSetStrategy{
			UpdateStrategyType:    pb.GameServerSetStrategy_OnDelete,
			PodUpdateStrategyType: pb.GameServerSetStrategy_Recreate, // 默认使用重建策略
		}

	default:
		// 如果类型不匹配，返回 nil
		return &pb.GameServerSetStrategy{}
	}
}

// convertConditions converts the conditions of deployment to the pb.DeploymentCondition type
func (x *GameServerSetService) convertConditions(conditions []kruiseappsv1alpha1.CloneSetCondition) []*pb.CloneSetCondition {
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

// convertDO2DTO converts the cloneSet object to the pb.CloneSetItem type
func (x *GameServerSetService) convertDO2DTO(gameServerSet *gamekruisev1alpha1.GameServerSet) *pb.GameServerSetItem {
	if gameServerSet.APIVersion == "" {
		gameServerSet.APIVersion = "game.kruise.io/v1alpha1"
	}
	if gameServerSet.Kind == "" {
		gameServerSet.Kind = "GameServerSet"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(gameServerSet)
	conditions := make([]*pb.GameServerSetCondition, 0)
	selector := &metav1.LabelSelector{}
	return &pb.GameServerSetItem{
		Name:           gameServerSet.Name,
		Namespace:      gameServerSet.Namespace,
		Status:         x.convertStatus(gameServerSet.Status),
		Containers:     x.convertContainers(gameServerSet.Spec.GameServerTemplate.Spec.Containers),
		CreateTime:     uint64(gameServerSet.CreationTimestamp.UnixNano() / 1e6),
		Labels:         gameServerSet.Labels,
		Annotations:    gameServerSet.Annotations,
		Selector:       selector,
		UpdateStrategy: x.convertGameServerSetStrategy(gameServerSet.Spec.UpdateStrategy),
		Replicas:       uint32(gameServerSet.Status.Replicas),
		SpecReplicas:   uint32(*gameServerSet.Spec.Replicas),
		Conditions:     conditions,
		Yaml:           yamlStr,
	}
}

func (x *GameServerSetService) ConvertReplicaSet2DTOV2(replicaset *appsv1.ReplicaSet) *pb.ReplicaSetItem {
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

func (x *GameServerSetService) ListGameServerSet(ctx context.Context, req *pb.ListGameServerSetRequest) (*pb.ListGameServerSetResponse, error) {
	gameServerSets, total, err := x.uc.ListGameServerSet(ctx, &biz.LisGameServerSetRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
		Keyword:  req.Keyword,
	})
	if errors.IsNotFound(err) {
		return &pb.ListGameServerSetResponse{
			List:  make([]*pb.GameServerSetItem, 0),
			Total: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	list := make([]*pb.GameServerSetItem, 0, len(gameServerSets))
	for _, gameServerSet := range gameServerSets {
		dto := x.convertDO2DTO(gameServerSet)
		list = append(list, dto)
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListGameServerSetResponse{
		List:  list,
		Total: int32(total),
	}, nil
}

// CreateOrUpdateGameServerSetByYaml creates or updates the cloneSet by yaml
func (x *GameServerSetService) CreateOrUpdateGameServerSetByYaml(ctx context.Context, req *pb.CreateOrUpdateGameServerSetByYamlRequest) (*pb.CreateOrUpdateGameServerSetByYamlResponse, error) {
	err := x.uc.CreateOrUpdateGameServerSetByYaml(ctx, &biz.CreateOrUpdateGameServerSetByYamlRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
		},
		Yaml: req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateGameServerSetByYamlResponse{}, nil
}

func (x *GameServerSetService) DeleteGameServerSet(ctx context.Context, req *pb.DeleteGameServerSetRequest) (*pb.DeleteGameServerSetResponse, error) {
	err := x.uc.DeleteGameServerSet(ctx, &biz.DeleteGameServerSetRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		}})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteGameServerSetResponse{}, nil
}

func (x *GameServerSetService) RestartGameServerSet(ctx context.Context, req *pb.RestartGameServerSetRequest) (*pb.RestartGameServerSetResponse, error) {
	err := x.uc.RestartGameServerSet(ctx, &biz.RestartGameServerSetRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		}})
	if err != nil {
		return nil, err
	}
	return &pb.RestartGameServerSetResponse{}, nil
}

// DecodeControllerRevisionData2GameServerSet 解析ControllerRevision的数据
func DecodeControllerRevisionData2GameServerSet(controllerRevision *appsv1.ControllerRevision) (*gamekruisev1alpha1.GameServerSet, error) {
	spec := appsv1.StatefulSet{}
	gameServerSet := &gamekruisev1alpha1.GameServerSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(controllerRevision.Data.Raw), 4096)
	if err := decoder.Decode(&spec); err != nil {
		return nil, fmt.Errorf("解析失败, 请检查格式是否正确: %v", err)
	}
	gameServerSet.Spec.GameServerTemplate.PodTemplateSpec = spec.Spec.Template

	return gameServerSet, nil
}

func (x *GameServerSetService) UpdateUpgradeStrategy(ctx context.Context, req *pb.UpdateGameServerSetUpgradeStrategyRequest) (*pb.UpdateGameServerSetUpgradeStrategyResponse, error) {
	err := x.uc.UpdateUpgradeStrategy(ctx, &biz.UpdateGameServerSetUpgradeStrategyRequest{
		GameServerSetCommonParams: biz.GameServerSetCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		UpdateStrategyType:    uint32(req.UpdateStrategyType),
		PodUpdateStrategyType: uint32(req.PodUpdateStrategyType),
		GracePeriodSeconds:    req.GracePeriodSeconds,
		Partition:             req.Partition,
		MaxUnavailable:        req.MaxUnavailable,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateGameServerSetUpgradeStrategyResponse{}, nil
}
