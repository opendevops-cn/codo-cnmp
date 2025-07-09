package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapService struct {
	pb.UnimplementedConfigMapServer
	uc *biz.ConfigMapUseCase
	uf *biz.UserFollowUseCase
}

func NewConfigMapService(uc *biz.ConfigMapUseCase, uf *biz.UserFollowUseCase) *ConfigMapService {
	return &ConfigMapService{
		uc: uc,
		uf: uf,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *ConfigMapService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_ConfigMap,
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
func (x *ConfigMapService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.ConfigMapItem) error {
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

func (x *ConfigMapService) convertDO2DTO(configMap *corev1.ConfigMap) *pb.ConfigMapItem {
	var (
		createTime, updateTime time.Time
		yamlStr                string
	)
	for _, managedField := range configMap.ManagedFields {
		if managedField.Operation == metav1.ManagedFieldsOperationUpdate {
			if updateTime.IsZero() || managedField.Time.Time.After(updateTime) {
				updateTime = managedField.Time.Time
			}
		}
	}
	createTime = configMap.CreationTimestamp.Time

	// 如果没有更新记录，则更新时间等于创建时间
	if updateTime.IsZero() {
		updateTime = createTime
	}
	if configMap.Kind == "" {
		configMap.Kind = "ConfigMap"
	}
	if configMap.APIVersion == "" {
		configMap.APIVersion = "v1"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(configMap)
	if err != nil {
		yamlStr = ""
	}

	return &pb.ConfigMapItem{
		Name:        configMap.Name,
		Labels:      configMap.Labels,
		Annotations: configMap.Annotations,
		CreateTime:  uint64(createTime.UnixNano() / 1e6),
		UpdateTime:  uint64(updateTime.UnixNano() / 1e6),
		Data:        configMap.Data,
		RefCount:    0, // todo
		Yaml:        yamlStr,
	}
}

func (x *ConfigMapService) ListConfigMap(ctx context.Context, req *pb.ListConfigMapsRequest) (*pb.ListConfigMapsResponse, error) {
	configMaps, total, err := x.uc.ListConfigMap(ctx, &biz.ListConfigMapRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Keyword:     req.Keyword,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ConfigMapItem, 0, len(configMaps))
	for _, configMap := range configMaps {
		dto := x.convertDO2DTO(configMap)
		refCount, refs, err := x.uc.GetConfigMapReferences(ctx, &biz.ConfigMapCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   configMap.Namespace,
			Name:        configMap.Name,
		})
		if err != nil {
			continue
		}
		for _, ref := range refs {
			dto.Refs = append(dto.Refs, &pb.ConfigMapReference{
				Kind: ref["kind"],
				Name: ref["name"],
			})

		}
		dto.RefCount = uint32(refCount)
		list = append(list, dto)
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListConfigMapsResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *ConfigMapService) CreateOrUpdateConfigMapByYaml(ctx context.Context, req *pb.CreateOrUpdateConfigMapByYamlRequest) (*pb.CreateOrUpdateConfigMapByYamlResponse, error) {
	err := x.uc.CreateOrUpdateConfigMapByYaml(ctx, &biz.CreateOrUpdateConfigMapByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateConfigMapByYamlResponse{}, nil
}

func (x *ConfigMapService) CreateConfigMap(ctx context.Context, req *pb.CreateConfigMapRequest) (*pb.CreateConfigMapResponse, error) {
	err := x.uc.CreateConfigMap(ctx, &biz.CreateConfigMapRequest{
		ConfigMapCommonParams: biz.ConfigMapCommonParams{
			ClusterName: req.ClusterName,
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
		Labels:      req.Labels,
		Annotations: req.Annotations,
		Data:        req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateConfigMapResponse{}, nil
}

func (x *ConfigMapService) UpdateConfigMap(ctx context.Context, req *pb.UpdateConfigMapRequest) (*pb.UpdateConfigMapResponse, error) {
	err := x.uc.UpdateConfigMap(ctx, &biz.UpdateConfigMapRequest{
		ConfigMapCommonParams: biz.ConfigMapCommonParams{
			ClusterName: req.ClusterName,
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
		Labels:      req.Labels,
		Annotations: req.Annotations,
		Data:        req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateConfigMapResponse{}, nil
}

func (x *ConfigMapService) DeleteConfigMap(ctx context.Context, req *pb.DeleteConfigMapRequest) (*pb.DeleteConfigMapResponse, error) {
	err := x.uc.DeleteConfigMap(ctx, &biz.DeleteConfigMapRequest{
		ConfigMapCommonParams: biz.ConfigMapCommonParams{
			ClusterName: req.ClusterName,
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteConfigMapResponse{}, nil
}

func (x *ConfigMapService) GetConfigMapDetail(ctx context.Context, req *pb.ConfigMapDetailRequest) (*pb.ConfigMapDetailResponse, error) {
	configMap, err := x.uc.GetConfigMap(ctx, &biz.ConfigMapCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}
	detail := x.convertDO2DTO(configMap)
	refCount, refs, err := x.uc.GetConfigMapReferences(ctx, &biz.ConfigMapCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   configMap.Namespace,
		Name:        configMap.Name,
	})
	if err != nil {
		return nil, err
	}
	for _, ref := range refs {
		detail.Refs = append(detail.Refs, &pb.ConfigMapReference{
			Kind: ref["kind"],
			Name: ref["name"],
		})

	}
	detail.RefCount = uint32(refCount)
	return &pb.ConfigMapDetailResponse{
		Detail: detail,
	}, nil
}
