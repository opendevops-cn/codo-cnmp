package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"context"
	"fmt"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"sort"
)

type HpaService struct {
	pb.UnimplementedHPAServer
	uc *biz.HpaUseCase
	uf *biz.UserFollowUseCase
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *HpaService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Hpa,
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
func (x *HpaService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.HpaItem) error {
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

func (x *HpaService) GetHpaDetail(ctx context.Context, request *pb.HpaDetailRequest) (*pb.HpaDetailResponse, error) {
	hpa, err := x.uc.GetHpa(ctx, &biz.GetHpaRequest{
		HpaCommonParams: biz.HpaCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name}})
	if err != nil {
		return nil, err
	}
	return &pb.HpaDetailResponse{
		Detail: x.convertHpa2DTO(hpa),
	}, nil
}

func NewHpaService(uc *biz.HpaUseCase, uf *biz.UserFollowUseCase) *HpaService {
	return &HpaService{
		uc: uc,
		uf: uf,
	}
}

func (x *HpaService) convertMetric2DTO(metric []autoscalingv2.MetricSpec) map[string]string {
	result := make(map[string]string)
	for _, m := range metric {
		switch m.Type {
		case autoscalingv2.ResourceMetricSourceType:
			result["type"] = "resource"
			if m.Resource.Name == "cpu" {
				result["cpu"] = fmt.Sprintf("%d%%", *m.Resource.Target.AverageUtilization)
				result["current_cpu"] = fmt.Sprintf("%d%%", *m.Resource.Target.AverageUtilization)
			} else if m.Resource.Name == "memory" {
				memoryValue := m.Resource.Target.AverageValue.String()
				result["memory"] = fmt.Sprintf("%d%%", *m.Resource.Target.AverageUtilization)
				result["target"] = memoryValue
			}
		case autoscalingv2.PodsMetricSourceType:
			result["type"] = "pods"
			result["metric"] = m.Pods.Metric.Name
			result["target"] = m.Pods.Target.Value.String()

		case autoscalingv2.ObjectMetricSourceType:
			result["type"] = "object"
			result["metric"] = m.Object.Metric.Name
			result["target"] = m.Object.Target.Value.String()

		case autoscalingv2.ExternalMetricSourceType:
			result["type"] = "external"
			result["metric"] = m.External.Metric.Name
			result["target"] = m.External.Target.Value.String()
		default:

		}

	}

	return result

}

func (x *HpaService) convertCurrentMetric2DTO(status []autoscalingv2.MetricStatus) map[string]string {
	result := make(map[string]string)
	for _, m := range status {
		switch m.Type {
		case autoscalingv2.ResourceMetricSourceType:
			result["type"] = "resource"
			if m.Resource.Name == "cpu" {
				result["current_cpu"] = fmt.Sprintf("%d%%", *m.Resource.Current.AverageUtilization)
			} else if m.Resource.Name == "memory" {
				result["current_memory"] = fmt.Sprintf("%d%%", *m.Resource.Current.AverageUtilization)
			}
		default:

		}

	}

	return result

}

func (x *HpaService) convertHpa2DTO(hpa *autoscalingv2.HorizontalPodAutoscaler) *pb.HpaItem {
	targetMetric := x.convertMetric2DTO(hpa.Spec.Metrics)
	currentMetrics := x.convertCurrentMetric2DTO(hpa.Status.CurrentMetrics)
	if hpa.APIVersion == "" {
		hpa.APIVersion = "autoscaling/v2"
	}
	if hpa.Kind == "" {
		hpa.Kind = "HorizontalPodAutoscaler"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(hpa)
	return &pb.HpaItem{
		Name:                     hpa.Name,
		Namespace:                hpa.Namespace,
		MinReplicas:              uint32(*hpa.Spec.MinReplicas),
		MaxReplicas:              uint32(hpa.Spec.MaxReplicas),
		CurrentReplicas:          uint32(hpa.Status.CurrentReplicas),
		Workload:                 hpa.Spec.ScaleTargetRef.Name,
		WorkloadType:             hpa.Spec.ScaleTargetRef.Kind,
		TargetCpuUtilization:     targetMetric["cpu"],
		TargetMemoryUtilization:  targetMetric["memory"],
		CurrentCpuUtilization:    currentMetrics["current_cpu"],
		CurrentMemoryUtilization: currentMetrics["current_memory"],
		Labels:                   hpa.Labels,
		Annotations:              hpa.Annotations,
		CreateTime:               uint64(hpa.CreationTimestamp.UnixNano() / 1e6),
		UpdateTime:               uint64(hpa.CreationTimestamp.UnixNano() / 1e6),
		Yaml:                     yamlStr,
	}
}

func (x *HpaService) ListHpa(ctx context.Context, req *pb.ListHpaRequest) (*pb.ListHpaResponse, error) {
	hpaList, total, err := x.uc.ListHpa(ctx, &biz.ListHpaRequest{
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
	list := make([]*pb.HpaItem, 0, len(hpaList))
	for _, hpa := range hpaList {
		list = append(list, x.convertHpa2DTO(hpa))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListHpaResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *HpaService) DeleteHpa(ctx context.Context, req *pb.DeleteHpaRequest) (*pb.DeleteHpaResponse, error) {
	err := x.uc.DeleteHpa(ctx, &biz.DeleteHpaRequest{
		HpaCommonParams: biz.HpaCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})

	if err != nil {
		return nil, err
	}
	return &pb.DeleteHpaResponse{}, nil
}

func (x *HpaService) CreateOrUpdateHpaByYaml(ctx context.Context, req *pb.CreateOrUpdateHpaByYamlRequest) (*pb.CreateOrUpdateHpaByYamlResponse, error) {
	err := x.uc.CreateOrUpdateHpaByYaml(ctx, &biz.CreateOrUpdateHpaByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateHpaByYamlResponse{}, nil
}
