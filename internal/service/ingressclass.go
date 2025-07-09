package service

import (
	"context"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
)

type IngressClassService struct {
	pb.UnimplementedIngressClassServer
	uc *biz.IngressClassUseCase
}

func NewIngressClassService(uc *biz.IngressClassUseCase) *IngressClassService {
	return &IngressClassService{
		uc: uc,
	}
}

func (x *IngressClassService) ListIngressClass(ctx context.Context, request *pb.ListIngressClassRequest) (*pb.ListIngressClassResponse, error) {
	ingressClasses, total, err := x.uc.ListIngressClass(ctx, &biz.ListIngressClassRequest{
		IngressClassCommonParams: biz.IngressClassCommonParams{
			ClusterName: request.GetClusterName(),
		},
		Keyword:  request.GetKeyword(),
		Page:     request.GetPage(),
		PageSize: request.GetPageSize(),
		ListAll:  utils.IntToBool(int(request.GetListAll())),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.IngressClassItem, 0, total)
	for _, ingressClass := range ingressClasses {
		list = append(list, &pb.IngressClassItem{
			Name: ingressClass.Name,
		})
	}
	return &pb.ListIngressClassResponse{
		List: list,
	}, nil
}
