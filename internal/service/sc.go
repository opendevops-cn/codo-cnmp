package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	storagev1 "k8s.io/api/storage/v1"
	"sort"
)

type ScService struct {
	pb.UnimplementedStorageClassServer
	uc *biz.StorageClassUseCase
	uf *biz.UserFollowUseCase
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *ScService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_StorageClass,
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
func (x *ScService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.StorageClassItem) error {
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

func (x *ScService) ListStorageClass(ctx context.Context, req *pb.ListStorageClassRequest) (*pb.ListStorageClassResponse, error) {
	storageClasses, total, err := x.uc.ListStorageClass(ctx, &biz.ListStorageClassRequest{
		ClusterName: req.ClusterName,
		Keyword:     req.Keyword,
		ListAll:     utils.IntToBool(int(req.ListAll)),
		Page:        req.Page,
		PageSize:    req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.StorageClassItem, 0, len(storageClasses))
	for _, storageClass := range storageClasses {
		list = append(list, x.convertDO2DTO(storageClass))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListStorageClassResponse{
		List:  list,
		Total: total,
	}, nil

}

func (x *ScService) convertDO2DTO(class *storagev1.StorageClass) *pb.StorageClassItem {
	if class == nil {
		return &pb.StorageClassItem{}
	}
	if class.APIVersion == "" {
		class.APIVersion = "storage.k8s.io/v1"
	}
	if class.Kind == "" {
		class.Kind = "StorageClass"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(class)
	if err != nil {
		yamlStr = ""
	}
	var isDefault bool
	if class.GetAnnotations() != nil {
		if class.GetAnnotations()["storageclass.kubernetes.io/is-default-class"] == "true" {
			isDefault = true
		}
	}
	return &pb.StorageClassItem{
		Name:              class.Name,
		Provisioner:       class.Provisioner,
		CreateTime:        uint64(class.CreationTimestamp.UnixNano() / 1e6),
		ReclaimPolicy:     string(*class.ReclaimPolicy),
		VolumeBindingMode: string(*class.VolumeBindingMode),
		Yaml:              yamlStr,
		IsDefault:         isDefault,
	}
}

func NewScService(uc *biz.StorageClassUseCase, uf *biz.UserFollowUseCase) *ScService {
	return &ScService{
		uc: uc,
		uf: uf,
	}
}
