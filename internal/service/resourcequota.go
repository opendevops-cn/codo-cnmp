package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceQuotaService struct {
	pb.UnimplementedResourceQuotaServer
	uc *biz.ResourceQuotaUseCase
}

func (x *ResourceQuotaService) GetResourceQuotaDetail(ctx context.Context, request *pb.ResourceQuotaDetailRequest) (*pb.ResourceQuotaDetailResponse, error) {
	resourceQuota, err := x.uc.GetResourceQuota(ctx, &biz.GetResourceQuotaRequest{
		ResourceQuotaCommonParams: biz.ResourceQuotaCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name}})
	if err != nil {
		return nil, err
	}
	return &pb.ResourceQuotaDetailResponse{
		Detail: x.convertResourceQuota2DTO(resourceQuota),
	}, nil
}

func NewResourceQuotaService(uc *biz.ResourceQuotaUseCase) *ResourceQuotaService {
	return &ResourceQuotaService{
		uc: uc,
	}
}

func (x *ResourceQuotaService) convertDTOToHard(resourceQuota *pb.HardResource) corev1.ResourceList {
	hard := corev1.ResourceList{}
	// 计算资源配额
	if resourceQuota.CpuLimit != "" {
		hard[corev1.ResourceLimitsCPU] = resource.MustParse(normalizeNumberString(resourceQuota.CpuLimit))
	}
	if resourceQuota.MemLimit != "" {
		hard[corev1.ResourceLimitsMemory] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(resourceQuota.MemLimit)))
	}
	if resourceQuota.CpuRequest != "" {
		hard[corev1.ResourceRequestsCPU] = resource.MustParse(normalizeNumberString(resourceQuota.CpuRequest))
	}
	if resourceQuota.MemRequest != "" {
		hard[corev1.ResourceRequestsMemory] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(resourceQuota.MemRequest)))
	}
	// 存储资源配额
	if resourceQuota.MaxResourceStorage != "" {
		hard[corev1.ResourceRequestsStorage] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(resourceQuota.MaxResourceStorage)))
	}
	if resourceQuota.MaxPersistentVolumeClaim != "" {
		hard[corev1.ResourcePersistentVolumeClaims] = resource.MustParse(normalizeNumberString(resourceQuota.MaxPersistentVolumeClaim))
	}

	// 其他资源配额
	if resourceQuota.MaxPod != "" {
		hard[corev1.ResourcePods] = resource.MustParse(normalizeNumberString(resourceQuota.MaxPod))
	}
	if resourceQuota.MaxConfigmap != "" {
		hard[corev1.ResourceConfigMaps] = resource.MustParse(normalizeNumberString(resourceQuota.MaxConfigmap))
	}
	if resourceQuota.MaxSecret != "" {
		hard[corev1.ResourceSecrets] = resource.MustParse(normalizeNumberString(resourceQuota.MaxSecret))
	}
	if resourceQuota.MaxService != "" {
		hard[corev1.ResourceServices] = resource.MustParse(normalizeNumberString(resourceQuota.MaxService))
	}
	return hard
}

func (x *ResourceQuotaService) convertResourceList2DTO(resourceList corev1.ResourceList) *pb.HardResource {
	// 返回string类型
	quantityToString := func(name corev1.ResourceName) string {
		// 判断key是否存在
		quantityValue, ok := resourceList[name]
		if !ok {
			return ""
		}
		return quantityValue.String()
	}
	if resourceList == nil {
		return &pb.HardResource{}
	}
	return &pb.HardResource{
		// 计算资源配额
		CpuLimit:   resourceToString(resourceList, corev1.ResourceLimitsCPU),
		MemLimit:   resourceToString(resourceList, corev1.ResourceLimitsMemory),
		MemRequest: resourceToString(resourceList, corev1.ResourceRequestsMemory),
		CpuRequest: resourceToString(resourceList, corev1.ResourceRequestsCPU),
		// 存储资源配额
		MaxResourceStorage:       resourceToString(resourceList, corev1.ResourceRequestsStorage),
		MaxPersistentVolumeClaim: resourceToString(resourceList, corev1.ResourcePersistentVolumeClaims),
		// 其他资源配额
		MaxPod:       quantityToString(corev1.ResourcePods),
		MaxConfigmap: quantityToString(corev1.ResourceConfigMaps),
		MaxSecret:    quantityToString(corev1.ResourceSecrets),
		MaxService:   quantityToString(corev1.ResourceServices),
	}
}

func (x *ResourceQuotaService) convertResourceQuota2DTO(resourceQuota *corev1.ResourceQuota) *pb.ResourceQuotaItem {

	if resourceQuota == nil {
		return &pb.ResourceQuotaItem{}

	}
	var (
		scopes  []string
		yamlStr string
	)
	for _, scope := range resourceQuota.Spec.Scopes {
		scopes = append(scopes, string(scope))
	}
	if resourceQuota.APIVersion == "" {
		resourceQuota.APIVersion = "v1"
	}
	if resourceQuota.Kind == "" {
		resourceQuota.Kind = "ResourceQuota"
	}
	yamlStr, _ = utils.ConvertResourceTOYaml(resourceQuota)

	return &pb.ResourceQuotaItem{
		Name:         resourceQuota.Name,
		Namespace:    resourceQuota.Namespace,
		HardResource: x.convertResourceList2DTO(resourceQuota.Spec.Hard),
		UsedResource: x.convertResourceList2DTO(resourceQuota.Status.Used),
		CreateTime:   uint64(resourceQuota.CreationTimestamp.UnixNano() / 1e6),
		Scopes:       scopes,
		Yaml:         yamlStr,
	}
}

func (x *ResourceQuotaService) ListResourceQuota(ctx context.Context, req *pb.ListResourceQuotaRequest) (*pb.ListResourceQuotaResponse, error) {
	resourceQuotaList, total, err := x.uc.ListResourceQuota(ctx, &biz.ListResourceQuotaRequest{
		ResourceQuotaCommonParams: biz.ResourceQuotaCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
		},
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ResourceQuotaItem, 0, len(resourceQuotaList))
	for _, resourceQuota := range resourceQuotaList {
		list = append(list, x.convertResourceQuota2DTO(resourceQuota))
	}
	return &pb.ListResourceQuotaResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *ResourceQuotaService) CreateResourceQuota(ctx context.Context, req *pb.CreateResourceQuotaRequest) (*pb.CreateResourceQuotaResponse, error) {
	Hard := x.convertDTOToHard(req.HardResource)
	err := x.uc.CreateResourceQuota(ctx, &biz.CreateResourceQuotaRequest{
		ResourceQuotaCommonParams: biz.ResourceQuotaCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Hard: Hard,
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResourceQuotaResponse{}, nil
}

func (x *ResourceQuotaService) UpdateResourceQuota(ctx context.Context, req *pb.CreateResourceQuotaRequest) (*pb.CreateResourceQuotaResponse, error) {
	Hard := x.convertDTOToHard(req.HardResource)
	err := x.uc.UpdateResourceQuota(ctx, &biz.UpdateResourceQuotaRequest{
		ResourceQuotaCommonParams: biz.ResourceQuotaCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Hard: Hard,
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResourceQuotaResponse{}, nil
}

func (x *ResourceQuotaService) CreateOrUpdateResourceQuota(ctx context.Context, req *pb.CreateResourceQuotaRequest) (*pb.CreateResourceQuotaResponse, error) {
	Hard := x.convertDTOToHard(req.HardResource)
	err := x.uc.CreateOrUpdateResourceQuota(ctx, &biz.CreateResourceQuotaRequest{
		ResourceQuotaCommonParams: biz.ResourceQuotaCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Hard: Hard,
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResourceQuotaResponse{}, nil
}

func (x *ResourceQuotaService) DeleteResourceQuota(ctx context.Context, req *pb.DeleteResourceQuotaRequest) (*pb.DeleteResourceQuotaResponse, error) {
	err := x.uc.DeleteResourceQuota(ctx, &biz.DeleteResourceQuotaRequest{
		ResourceQuotaCommonParams: biz.ResourceQuotaCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResourceQuotaResponse{}, nil
}
