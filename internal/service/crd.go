package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CRDService struct {
	pb.UnimplementedCRDServer
	uc *biz.CRDUseCase
	uf *biz.UserFollowUseCase
}

func (x *CRDService) ListCRD(ctx context.Context, request *pb.ListCRDRequest) (*pb.ListCRDResponse, error) {
	crds, total, err := x.uc.ListCRD(ctx, &biz.ListCRDRequest{
		ClusterName: request.GetClusterName(),
		Namespace:   request.GetNamespace(),
		Keyword:     request.GetKeyword(),
		Page:        request.GetPage(),
		PageSize:    request.GetPageSize(),
		ListAll:     utils.IntToBool(int(request.ListAll)),
		ApiGroup:    request.GetApiGroup(),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.CRDItem, 0, len(crds))
	for _, crd := range crds {
		list = append(list, x.convertDO2DTO(crd))
	}
	return &pb.ListCRDResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *CRDService) convertDO2DTO(crd *apiextensionsv1.CustomResourceDefinition) *pb.CRDItem {
	yamlStr, _ := utils.ConvertResourceTOYaml(crd)
	return &pb.CRDItem{
		Name:       crd.Name,
		ApiGroup:   crd.Spec.Group,
		ApiVersion: crd.Spec.Versions[0].Name,
		Scope:      string(crd.Spec.Scope),
		CreateTime: uint64(crd.CreationTimestamp.UnixNano() / 1e6),
		Kind:       crd.Spec.Names.Kind,
		Yaml:       yamlStr,
	}

}

func (x *CRDService) ListCRDInstance(ctx context.Context, request *pb.ListCRDInstanceRequest) (*pb.ListCRDInstanceResponse, error) {
	instances, total, err := x.uc.ListCRDInstances(ctx, &biz.ListCRDInstancesRequest{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
		Name:        request.Name,
		ApiVersion:  request.ApiVersion,
		ApiGroup:    request.ApiGroup,
		Keyword:     request.Keyword,
		Page:        request.Page,
		PageSize:    request.PageSize,
		ListAll:     utils.IntToBool(int(request.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.CRDInstanceItem, 0, len(instances))
	for _, instance := range instances {
		list = append(list, x.convertInstanceDO2DTO(instance))
	}
	return &pb.ListCRDInstanceResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *CRDService) convertInstanceDO2DTO(instance *unstructured.Unstructured) *pb.CRDInstanceItem {
	yamlStr, _ := utils.ConvertResourceTOYaml(instance)
	return &pb.CRDInstanceItem{
		Name:       instance.GetName(),
		Namespace:  instance.GetNamespace(),
		CreateTime: uint64(instance.GetCreationTimestamp().UnixNano() / 1e6),
		ApiVersion: instance.GetAPIVersion(),
		Yaml:       yamlStr,
	}
}

func (x *CRDService) DeleteCRD(ctx context.Context, request *pb.DeleteCRDRequest) (*pb.DeleteCRDResponse, error) {
	err := x.uc.DeleteCRD(ctx, &biz.DeleteCRDRequest{
		CRDCommonParams: biz.CRDCommonParams{
			ClusterName: request.ClusterName,
			Name:        request.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCRDResponse{}, nil
}

func NewCRDService(uc *biz.CRDUseCase, uf *biz.UserFollowUseCase) *CRDService {
	return &CRDService{
		uc: uc,
		uf: uf,
	}
}
