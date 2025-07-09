package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
	"strings"
)

type LimitRangeService struct {
	pb.UnimplementedLimitRangeServer
	uc *biz.LimitRangeUseCase
}

func (x *LimitRangeService) GetLimitRangeDetail(ctx context.Context, request *pb.LimitRangeDetailRequest) (*pb.LimitRangeDetailResponse, error) {
	limitRange, err := x.uc.GetLimitRange(ctx, &biz.GetLimitRangeRequest{
		LimitRangeCommonParams: biz.LimitRangeCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name}})
	if err != nil {
		return nil, err
	}
	return &pb.LimitRangeDetailResponse{
		Detail: x.convertLimitRange2DTO(limitRange),
	}, nil
}

func NewLimitRangeService(uc *biz.LimitRangeUseCase) *LimitRangeService {
	return &LimitRangeService{
		uc: uc,
	}
}

func resourceToString(resourceList corev1.ResourceList, name corev1.ResourceName) string {
	quantity, ok := resourceList[name]
	if !ok {
		return ""
	}
	switch name {
	case corev1.ResourceCPU, corev1.ResourceLimitsCPU, corev1.ResourceRequestsCPU:
		return quantityToCoresString(quantity)
	case corev1.ResourceMemory, corev1.ResourceEphemeralStorage, corev1.ResourceStorage, corev1.ResourceRequestsMemory,
		corev1.ResourceRequestsStorage, corev1.ResourceLimitsMemory, corev1.ResourceRequestsEphemeralStorage:
		return quantityToGiString(quantity)
	default:
		return quantity.String()
	}
}

func quantityToCoresString(q resource.Quantity) string {
	// 将毫核值转换为核数
	if q.IsZero() {
		return "0"
	}
	cores := float64(q.MilliValue()) / 1000.0
	return fmt.Sprintf("%.2f", cores) // 格式化为最多2位有效数字
}

func quantityToGiString(q resource.Quantity) string {
	// 将字节值转换为兆字节
	if q.IsZero() {
		return "0"
	}
	memory := float64(q.Value()) / 1024.0 / 1024.0 / 1024.0
	return fmt.Sprintf("%.2f", memory) // 格式化为最多2位有效数字
}

func float64StringToInt64(str string) int64 {
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return int64(floatValue)
}

func normalizeNumberString(str string) string {
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return str // 如果解析失败，返回原始字符串
	}
	// 检查是否为整数
	if floatValue == float64(int64(floatValue)) {
		// 如果是整数，返回整数形式
		return strconv.FormatInt(int64(floatValue), 10)
	}
	// 如果是小数，返回原始小数形式
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", floatValue), "0"), ".")
}

// 提取 Container 限制
func extractContainerLimitRange(limit corev1.LimitRangeItem) *pb.ContainerLimitRange {
	var containerLimits pb.ContainerLimitRange
	if limit.Default != nil {
		containerLimits.DefaultCpu = quantityToCoresString(*limit.Default.Cpu())
		containerLimits.DefaultMem = quantityToGiString(*limit.Default.Memory())
	}
	if limit.DefaultRequest != nil {
		containerLimits.DefaultCpuRequest = quantityToCoresString(*limit.DefaultRequest.Cpu())
		containerLimits.DefaultMemRequest = quantityToGiString(*limit.DefaultRequest.Memory())
	}
	if limit.Max != nil {
		containerLimits.MaxCpu = quantityToCoresString(*limit.Max.Cpu())
		containerLimits.MaxMem = quantityToGiString(*limit.Max.Memory())
	}
	if limit.Min != nil {
		containerLimits.MinCpu = quantityToCoresString(*limit.Min.Cpu())
		containerLimits.MinMem = quantityToGiString(*limit.Min.Memory())
	}
	return &containerLimits
}

// 提取 Pod 限制
func extractPodLimitRange(limit corev1.LimitRangeItem) *pb.PodLimitRange {
	var podLimits pb.PodLimitRange
	if limit.Max != nil {
		podLimits.MaxCpu = quantityToCoresString(*limit.Max.Cpu())
		podLimits.MaxMem = quantityToGiString(*limit.Max.Memory())
	}
	return &podLimits
}

// 提取 PVC 限制
func extractPVCLimitRange(limit corev1.LimitRangeItem) *pb.PersistentVolumeClaimLimitRange {
	var pvcLimits pb.PersistentVolumeClaimLimitRange
	if limit.Min != nil {
		pvcLimits.MinResourceStorageRequest = quantityToGiString(*limit.Min.Storage())
	}
	if limit.Max != nil {
		pvcLimits.MaxResourceStorageRequest = quantityToGiString(*limit.Max.Storage())
	}
	return &pvcLimits
}

func convertDTOToLimitRange(dto *pb.CreateLimitRangeRequest) []corev1.LimitRangeItem {
	var limits []corev1.LimitRangeItem

	// 处理 Container 限制
	if dto.ContainerLimitRange != nil {
		limits = append(limits, convertContainerDTOToLimitRangeItem(dto.ContainerLimitRange))
	}

	// 处理 Pod 限制
	if dto.PodLimitRange != nil {
		limits = append(limits, convertPodDTOToLimitRangeItem(dto.PodLimitRange))
	}

	// 处理 PVC 限制
	if dto.PersistentVolumeClaimLimitRange != nil {
		limits = append(limits, convertPVCDTOToLimitRangeItem(dto.PersistentVolumeClaimLimitRange))
	}

	return limits
}

// 处理 Container 限制
func convertContainerDTOToLimitRangeItem(containerDTO *pb.ContainerLimitRange) corev1.LimitRangeItem {
	item := corev1.LimitRangeItem{
		Type: corev1.LimitTypeContainer,
	}
	item.Default = make(corev1.ResourceList)
	item.DefaultRequest = make(corev1.ResourceList)
	item.Max = make(corev1.ResourceList)
	item.Min = make(corev1.ResourceList)
	if containerDTO.DefaultCpu != "" {
		item.Default[corev1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%s", normalizeNumberString(containerDTO.DefaultCpu)))
	}
	if containerDTO.DefaultMem != "" {
		item.Default[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(containerDTO.DefaultMem)))
	}
	if containerDTO.DefaultCpuRequest != "" {
		item.DefaultRequest[corev1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%s", normalizeNumberString(containerDTO.DefaultCpuRequest)))
	}
	if containerDTO.DefaultMemRequest != "" {
		item.DefaultRequest[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(containerDTO.DefaultMemRequest)))
	}
	if containerDTO.MaxCpu != "" {
		item.Max[corev1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%s", normalizeNumberString(containerDTO.MaxCpu)))
	}
	if containerDTO.MaxMem != "" {
		item.Max[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(containerDTO.MaxMem)))
	}
	if containerDTO.MinCpu != "" {
		item.Min[corev1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%s", normalizeNumberString(containerDTO.MinCpu)))
	}
	if containerDTO.MinMem != "" {
		item.Min[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%sGi", normalizeNumberString(containerDTO.MinMem)))
	}
	return item
}

// 处理 Pod 限制
func convertPodDTOToLimitRangeItem(podDTO *pb.PodLimitRange) corev1.LimitRangeItem {
	item := corev1.LimitRangeItem{
		Type: corev1.LimitTypePod,
	}
	item.Max = make(corev1.ResourceList)
	if podDTO.MaxCpu != "" {
		item.Max[corev1.ResourceCPU] = resource.MustParse(podDTO.MaxCpu)
	}
	if podDTO.MaxMem != "" {
		item.Max[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%sGi", podDTO.MaxMem))
	}
	return item
}

// 处理 PVC 限制
func convertPVCDTOToLimitRangeItem(pvcDTO *pb.PersistentVolumeClaimLimitRange) corev1.LimitRangeItem {
	item := corev1.LimitRangeItem{
		Type: corev1.LimitTypePersistentVolumeClaim,
	}
	item.Min = make(corev1.ResourceList)
	item.Max = make(corev1.ResourceList)
	if pvcDTO.MinResourceStorageRequest != "" {
		item.Min[corev1.ResourceStorage] = resource.MustParse(fmt.Sprintf("%sGi", pvcDTO.MinResourceStorageRequest))
	}
	if pvcDTO.MaxResourceStorageRequest != "" {
		item.Max[corev1.ResourceStorage] = resource.MustParse(fmt.Sprintf("%sGi", pvcDTO.MaxResourceStorageRequest))
	}
	return item
}

func (x *LimitRangeService) convertLimitRange2DTO(limitRange *corev1.LimitRange) *pb.LimitRangeItem {
	if limitRange == nil {
		return &pb.LimitRangeItem{}
	}
	var containerLimits pb.ContainerLimitRange
	var podLimits pb.PodLimitRange
	var pvcLimits pb.PersistentVolumeClaimLimitRange
	for _, limit := range limitRange.Spec.Limits {
		switch limit.Type {
		case corev1.LimitTypeContainer:
			containerLimits = *extractContainerLimitRange(limit)
		case corev1.LimitTypePod:
			podLimits = *extractPodLimitRange(limit)
		case corev1.LimitTypePersistentVolumeClaim:
			pvcLimits = *extractPVCLimitRange(limit)
		default:
			// ignore other limit type
		}
	}
	if limitRange.APIVersion == "" {
		limitRange.APIVersion = "v1"
	}
	if limitRange.Kind == "" {
		limitRange.Kind = "LimitRange"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(limitRange)
	return &pb.LimitRangeItem{
		Name:                            limitRange.Name,
		Namespace:                       limitRange.Namespace,
		CreateTime:                      uint64(limitRange.CreationTimestamp.UnixNano() / 1e6),
		ContainerLimitRange:             &containerLimits,
		PodLimitRange:                   &podLimits,
		PersistentVolumeClaimLimitRange: &pvcLimits,
		Yaml:                            yamlStr,
	}
}

func (x *LimitRangeService) ListLimitRange(ctx context.Context, req *pb.ListLimitRangeRequest) (*pb.ListLimitRangeResponse, error) {
	limitRangeList, total, err := x.uc.ListLimitRange(ctx, &biz.ListLimitRangeRequest{
		LimitRangeCommonParams: biz.LimitRangeCommonParams{
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
	list := make([]*pb.LimitRangeItem, 0, len(limitRangeList))
	for _, limitRange := range limitRangeList {
		list = append(list, x.convertLimitRange2DTO(limitRange))
	}
	return &pb.ListLimitRangeResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *LimitRangeService) CreateLimitRange(ctx context.Context, req *pb.CreateLimitRangeRequest) (*pb.CreateLimitRangeResponse, error) {
	Limits := convertDTOToLimitRange(req)
	err := x.uc.CreateLimitRange(ctx, &biz.CreateLimitRangeRequest{
		LimitRangeCommonParams: biz.LimitRangeCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Limits: Limits,
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateLimitRangeResponse{}, nil
}

func (x *LimitRangeService) UpdateLimitRange(ctx context.Context, req *pb.CreateLimitRangeRequest) (*pb.CreateLimitRangeResponse, error) {
	Limits := convertDTOToLimitRange(req)
	err := x.uc.UpdateLimitRange(ctx, &biz.UpdateLimitRangeRequest{
		LimitRangeCommonParams: biz.LimitRangeCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Limits: Limits,
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateLimitRangeResponse{}, nil
}

func (x *LimitRangeService) CreateOrUpdateLimitRange(ctx context.Context, req *pb.CreateLimitRangeRequest) (*pb.CreateLimitRangeResponse, error) {
	Limits := convertDTOToLimitRange(req)
	err := x.uc.CreateOrUpdateLimitRange(ctx, &biz.CreateLimitRangeRequest{
		LimitRangeCommonParams: biz.LimitRangeCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Limits: Limits,
	},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateLimitRangeResponse{}, nil
}

func (x *LimitRangeService) DeleteLimitRange(ctx context.Context, req *pb.DeleteLimitRangeRequest) (*pb.DeleteLimitRangeResponse, error) {
	err := x.uc.DeleteLimitRange(ctx, &biz.DeleteLimitRangeRequest{
		LimitRangeCommonParams: biz.LimitRangeCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteLimitRangeResponse{}, nil

}
