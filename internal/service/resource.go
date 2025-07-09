package service

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
)

type ResourceService struct {
	pb.UnimplementedResourceServer
	uc *biz.ResourceUseCase
}

func (x *ResourceService) DryRunResource(ctx context.Context, req *pb.CreateOrUpdateResourceRequest) (*pb.CreateOrUpdateResourceResponse, error) {
	err := x.uc.DryRunResourceByYaml(ctx, &biz.CreateOrUpdateResourceRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	var (
		success bool
		message string
	)
	if err != nil {
		success = false
		message = err.Error()
	} else {
		success = true
		message = "success"
	}
	return &pb.CreateOrUpdateResourceResponse{
		Success: success,
		Message: message,
	}, nil
}

func NewResourceService(uc *biz.ResourceUseCase) *ResourceService {
	return &ResourceService{uc: uc}
}

func (x *ResourceService) CreateOrUpdateResource(ctx context.Context, req *pb.CreateOrUpdateResourceRequest) (*pb.CreateOrUpdateResourceResponse, error) {
	err := x.uc.CreateOrUpdateResourceByYaml(ctx, &biz.CreateOrUpdateResourceRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	return &pb.CreateOrUpdateResourceResponse{
		Success: err == nil,
	}, err
}
