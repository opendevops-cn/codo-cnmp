package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sort"
)

type PvcService struct {
	pb.UnimplementedPersistentVolumeServer
	uc *biz.PersistentVolumeClaimUseCase
	uf *biz.UserFollowUseCase
}

func (x *PvcService) DeletePersistentVolumeClaim(ctx context.Context, request *pb.DeletePersistentVolumeClaimRequest) (*pb.DeletePersistentVolumeClaimResponse, error) {
	err := x.uc.DeletePersistentVolumeClaim(ctx, &biz.DeletePersistentVolumeClaimRequest{
		ClusterName: request.ClusterName,
		NameSpace:   request.Namespace,
		Name:        request.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeletePersistentVolumeClaimResponse{}, nil
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *PvcService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_PersistentVolumeClaim,
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
func (x *PvcService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.PersistentVolumeClaimItem) error {
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

func (x *PvcService) ListPersistentVolumeClaim(ctx context.Context, req *pb.ListPersistentVolumeClaimRequest) (*pb.ListPersistentVolumeClaimResponse, error) {
	pvcs, total, err := x.uc.ListPersistentVolumeClaim(ctx, &biz.ListPersistentVolumeClaimRequest{
		ClusterName: req.ClusterName,
		NameSpace:   req.GetNamespace(),
		Keyword:     req.Keyword,
		PageSize:    req.PageSize,
		Page:        req.Page,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.PersistentVolumeClaimItem, 0, len(pvcs))
	for _, pvc := range pvcs {
		list = append(list, x.convertDO2DTO(pvc))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListPersistentVolumeClaimResponse{
		List:  list,
		Total: total,
	}, nil
}

func getAccessModes(AccessModes []corev1.PersistentVolumeAccessMode) []string {
	var accessModes []string
	if len(AccessModes) == 0 {
		return accessModes
	}
	for _, mode := range AccessModes {
		accessModes = append(accessModes, string(mode))
	}
	return accessModes
}

func (x *PvcService) convertDO2DTO(pvc *corev1.PersistentVolumeClaim) *pb.PersistentVolumeClaimItem {
	if pvc == nil {
		return &pb.PersistentVolumeClaimItem{}
	}
	if pvc.APIVersion == "" {
		pvc.APIVersion = "v1"
	}
	if pvc.Kind == "" {
		pvc.Kind = "PersistentVolumeClaim"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(pvc)
	if err != nil {
		yamlStr = ""
	}
	var storageClassName string
	if pvc.Spec.StorageClassName != nil {
		storageClassName = *pvc.Spec.StorageClassName
	}
	return &pb.PersistentVolumeClaimItem{
		Name:             pvc.Name,
		Yaml:             yamlStr,
		Capacity:         pvc.Spec.Resources.Requests.Storage().String(),
		AccessModes:      getAccessModes(pvc.Spec.AccessModes),
		Status:           string(pvc.Status.Phase),
		StorageClassName: storageClassName,
		VolumeName:       pvc.Spec.VolumeName,
		CreateTime:       uint64(pvc.CreationTimestamp.UnixNano() / 1e6),
	}
}

func NewPvcService(uc *biz.PersistentVolumeClaimUseCase, uf *biz.UserFollowUseCase) *PvcService {
	return &PvcService{uc: uc, uf: uf}
}
