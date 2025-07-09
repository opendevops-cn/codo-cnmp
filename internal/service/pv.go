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

type PvService struct {
	pb.UnimplementedPersistentVolumeServer
	uc *biz.PersistentVolumeUseCase
	uf *biz.UserFollowUseCase
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *PvService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_PersistentVolume,
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
func (x *PvService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.PersistentVolumeItem) error {
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

func (x *PvService) ListPersistentVolume(ctx context.Context, req *pb.ListPersistentVolumeRequest) (*pb.ListPersistentVolumeResponse, error) {
	pv, total, err := x.uc.ListPersistentVolume(ctx, &biz.ListPersistentVolumeRequest{
		ClusterName: req.ClusterName,
		Keyword:     req.Keyword,
		PageSize:    req.PageSize,
		Page:        req.Page,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.PersistentVolumeItem, 0, len(pv))
	for _, p := range pv {
		list = append(list, x.convertDO2DTO(p))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListPersistentVolumeResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *PvService) convertDO2DTO(pv *corev1.PersistentVolume) *pb.PersistentVolumeItem {
	if pv == nil {
		return &pb.PersistentVolumeItem{}
	}
	if pv.APIVersion == "" {
		pv.APIVersion = "v1"
	}
	if pv.Kind == "" {
		pv.Kind = "PersistentVolume"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(pv)
	if err != nil {
		yamlStr = ""
	}
	return &pb.PersistentVolumeItem{
		Name:             pv.Name,
		Yaml:             yamlStr,
		Capacity:         pv.Spec.Capacity.Storage().String(),
		AccessModes:      getAccessModes(pv.Spec.AccessModes),
		Status:           string(pv.Status.Phase),
		StorageClassName: pv.Spec.StorageClassName,
		VolumeMode:       string(*pv.Spec.VolumeMode),
		ClaimRef:         x.convertClaimRef2DTO(pv.Spec.ClaimRef),
		CreateTime:       uint64(pv.CreationTimestamp.UnixNano() / 1e6),
	}
}

func (x *PvService) convertClaimRef2DTO(obj *corev1.ObjectReference) *pb.ClaimRef {
	if obj == nil {
		return &pb.ClaimRef{}
	}
	return &pb.ClaimRef{
		Name:      obj.Name,
		Namespace: obj.Namespace,
		Kind:      obj.Kind,
	}
}

func NewPVService(uc *biz.PersistentVolumeUseCase, uf *biz.UserFollowUseCase) *PvService {
	return &PvService{uc: uc, uf: uf}
}
