package service

import (
	"context"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApiGroupService struct {
	pb.UnimplementedAPIGroupServer
	uc *biz.ApiGroupUseCase
	uf *biz.UserFollowUseCase
}

func (x *ApiGroupService) ListAPIGroup(ctx context.Context, request *pb.ListAPIGroupRequest) (*pb.ListAPIGroupResponse, error) {
	apiGroups, _, err := x.uc.ListApiGroup(ctx, &biz.ListApiGroupRequest{
		ClusterName: request.GetClusterName(),
		Keyword:     request.GetKeyword(),
		Page:        request.GetPage(),
		PageSize:    request.GetPageSize(),
		ListAll:     utils.IntToBool(int(request.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.APIGroupItem, 0, len(apiGroups))
	for _, apiGroup := range apiGroups {
		items := x.convertDO2DTO(apiGroup)
		list = append(list, items...)
	}
	return &pb.ListAPIGroupResponse{
		List:  list,
		Total: uint32(len(list)),
	}, nil
}

func (x *ApiGroupService) convertDO2DTO(group *metav1.APIGroup) []*pb.APIGroupItem {
	var result []*pb.APIGroupItem
	for _, version := range group.Versions {
		if group.Name == "" {
			continue
		}
		result = append(result, &pb.APIGroupItem{
			Name:       group.Name,
			ApiVersion: version.Version,
		})
	}
	return result
}

func NewApiGroupService(uc *biz.ApiGroupUseCase, uf *biz.UserFollowUseCase) *ApiGroupService {
	return &ApiGroupService{
		uc: uc,
		uf: uf,
	}
}
