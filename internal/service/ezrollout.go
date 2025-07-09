package service

import (
	"context"
	"fmt"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sort"
	"strconv"
	"time"

	ezRolloutv1 "codo-cnmp/common/ezrollout/v1"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MetricOnlineNumber = "online_number"
	MetricEntityCount  = "entity_count"
	MetricCPUUsage     = "virtual_cpu"
	MetricMemoryRSS    = "virtual_memory"
)

var supportedMetrics = map[string]struct{}{
	MetricOnlineNumber: {},
	MetricEntityCount:  {},
	MetricCPUUsage:     {},
	MetricMemoryRSS:    {},
}

type EzRolloutService struct {
	pb.UnimplementedEzRolloutServer
	uc *biz.EzRolloutUseCase
	uf *biz.UserFollowUseCase
}

// getUserFollowMap 获取用户关注的版本伸缩列表
func (x *EzRolloutService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_EzRollout,
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
func (x *EzRolloutService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.EzRolloutInfo) error {
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

func NewEzRolloutService(uc *biz.EzRolloutUseCase, uf *biz.UserFollowUseCase) *EzRolloutService {
	return &EzRolloutService{uc: uc, uf: uf}
}

// isMetricSupported 检查指标是否支持
func isMetricSupported(metric string) bool {
	_, ok := supportedMetrics[metric]
	return ok
}

// FormatMillisToLocal 将毫秒时间戳转换为本地时间字符串
func FormatMillisToLocal(millis int64) string {
	return time.UnixMilli(millis).Local().Format(time.DateTime)
}

// convertSelector 转换标签选择器
func (x *EzRolloutService) convertSelector(selector *metav1.LabelSelector) map[string]string {
	if selector == nil {
		return nil
	}
	return selector.MatchLabels
}

// convertHpaRules 转换HPA规则
func (x *EzRolloutService) convertHpaRules(rules *autoscalingv2.HPAScalingRules) *pb.HPAScalingRules {
	if rules == nil {
		return nil
	}

	policies := make([]*pb.HPAScalingPolicy, 0, len(rules.Policies))
	for _, policy := range rules.Policies {
		policies = append(policies, &pb.HPAScalingPolicy{
			Type:          string(policy.Type),
			Value:         uint32(policy.Value),
			PeriodSeconds: uint32(policy.PeriodSeconds),
		})
	}

	return &pb.HPAScalingRules{
		StabilizationWindowSeconds: uint32(*rules.StabilizationWindowSeconds),
		SelectPolicy:               string(*rules.SelectPolicy),
		Policies:                   policies,
	}
}

// convertMetrics 转换指标配置
func (x *EzRolloutService) convertMetrics(metrics []autoscalingv2.MetricSpec) []*pb.MetricSpecV2 {
	if metrics == nil {
		return nil
	}

	result := make([]*pb.MetricSpecV2, 0, len(metrics))
	for _, m := range metrics {
		if m.Type != autoscalingv2.PodsMetricSourceType || m.Pods == nil || m.Pods.Metric.Name == "" {
			continue
		}

		var value uint32
		switch m.Pods.Target.Type {
		case autoscalingv2.ValueMetricType:
			if m.Pods.Target.Value != nil {
				value = quantityToNumber(m.Pods.Target.Value, m.Pods.Metric.Name)
			}
		case autoscalingv2.AverageValueMetricType:
			if m.Pods.Target.AverageValue != nil {
				value = quantityToNumber(m.Pods.Target.AverageValue, m.Pods.Metric.Name)
			}
		}

		result = append(result, &pb.MetricSpecV2{
			Name:         m.Pods.Metric.Name,
			Value:        value,
			CurrentValue: 0, // todo: 后续补充

		})
	}
	return result
}

// quantityToNumber 将 Quantity 转换为数字
func quantityToNumber(q *resource.Quantity, metricName string) uint32 {
	if q == nil {
		return 0
	}
	n := resource.MustParse(q.String())
	return uint32(n.Value())
}

// convertMetricsV22DO 转换指标配置到DO
func (x *EzRolloutService) convertMetricsV22DO(metrics []*pb.MetricSpecV2, isScaleUp bool) []autoscalingv2.MetricSpec {
	if metrics == nil {
		return nil
	}

	targetType := autoscalingv2.ValueMetricType
	if isScaleUp {
		targetType = autoscalingv2.AverageValueMetricType
	}

	result := make([]autoscalingv2.MetricSpec, 0, len(metrics))
	for _, metric := range metrics {
		if !isMetricSupported(metric.Name) {
			continue
		}
		q := resource.NewQuantity(int64(metric.Value), resource.DecimalSI)
		spec := autoscalingv2.MetricSpec{
			Type: autoscalingv2.PodsMetricSourceType,
			Pods: &autoscalingv2.PodsMetricSource{
				Metric: autoscalingv2.MetricIdentifier{Name: metric.Name},
				Target: autoscalingv2.MetricTarget{Type: targetType},
			},
		}

		if targetType == autoscalingv2.AverageValueMetricType {
			spec.Pods.Target.AverageValue = q
		} else {
			spec.Pods.Target.Value = q
		}

		result = append(result, spec)
	}
	return result
}

// convertOnlineScaler2DO 转换线上扩缩容配置到DO
func (x *EzRolloutService) convertOnlineScaler2DO(scaleUp, scaleDown *pb.HPAScalingRules, minReplicas, maxReplicas uint32) *ezRolloutv1.EzOnlineScaler {
	convertRules := func(rules *pb.HPAScalingRules) *autoscalingv2.HPAScalingRules {
		if rules == nil {
			return nil
		}

		stabilizationWindowSeconds := int32(rules.StabilizationWindowSeconds)
		selectPolicy := autoscalingv2.ScalingPolicySelect(rules.SelectPolicy)

		policies := make([]autoscalingv2.HPAScalingPolicy, 0, len(rules.Policies))
		for _, policy := range rules.Policies {
			policies = append(policies, autoscalingv2.HPAScalingPolicy{
				Type:          autoscalingv2.HPAScalingPolicyType(policy.Type),
				Value:         int32(policy.Value),
				PeriodSeconds: int32(policy.PeriodSeconds),
			})
		}

		return &autoscalingv2.HPAScalingRules{
			StabilizationWindowSeconds: &stabilizationWindowSeconds,
			SelectPolicy:               &selectPolicy,
			Policies:                   policies,
		}
	}

	minRs := int32(minReplicas)
	maxRs := int32(maxReplicas)

	return &ezRolloutv1.EzOnlineScaler{
		MinReplicas: &minRs,
		MaxReplicas: &maxRs,
		Behavior: &autoscalingv2.HorizontalPodAutoscalerBehavior{
			ScaleUp:   convertRules(scaleUp),
			ScaleDown: convertRules(scaleDown),
		},
	}
}

// 方法1：使用 strconv.ParseUint
func stringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// convertDO2DTO 转换DO到DTO
func (x *EzRolloutService) convertDO2DTO(ezRollout *ezRolloutv1.EzRollout) *pb.EzRolloutInfo {
	if ezRollout == nil {
		return nil
	}

	yamlStr, err := utils.ConvertResourceTOYaml(ezRollout)
	if err != nil {
		return nil
	}

	var minReplicas, maxReplicas uint32
	if ezRollout.Spec.OnlineScaler != nil {
		if ezRollout.Spec.OnlineScaler.MinReplicas != nil {
			minReplicas = uint32(*ezRollout.Spec.OnlineScaler.MinReplicas)
		}
		if ezRollout.Spec.OnlineScaler.MaxReplicas != nil {
			maxReplicas = uint32(*ezRollout.Spec.OnlineScaler.MaxReplicas)
		}
	}

	return &pb.EzRolloutInfo{
		Namespace:          ezRollout.Namespace,
		Name:               ezRollout.Name,
		OnlineVersion:      ezRollout.Spec.OnlineVersion,
		CurrentReplicas:    0,
		OfflineDeadline:    uint64(ezRollout.Spec.OfflineScaler.Deadline),
		Selector:           x.convertSelector(ezRollout.Spec.Selector),
		Labels:             ezRollout.Labels,
		Annotations:        ezRollout.Annotations,
		MinReplicas:        minReplicas,
		MaxReplicas:        maxReplicas,
		ScaleUpMetrics:     x.convertMetrics(ezRollout.Spec.ScaleUpMetrics),
		ScaleDownMetrics:   x.convertMetrics(ezRollout.Spec.ScaleDownMetrics),
		CreateTime:         uint64(ezRollout.CreationTimestamp.UnixNano() / 1e6),
		ScaleUp:            x.convertHpaRules(ezRollout.Spec.OnlineScaler.Behavior.ScaleUp),
		ScaleDown:          x.convertHpaRules(ezRollout.Spec.OnlineScaler.Behavior.ScaleDown),
		Ready:              ezRollout.Status.Ready,
		LatestErrorMessage: ezRollout.Status.LatestError.Message,
		LatestErrorTime:    uint64(ezRollout.Status.LatestError.Timestamp),
		Yaml:               yamlStr,
		EnableScaleUp:      ezRollout.Spec.OfflineScaler.EnableScaleUp,
	}
}

// ListEzRollout 获取EzRollout列表
func (x *EzRolloutService) ListEzRollout(ctx context.Context, req *pb.ListEzRolloutRequest) (*pb.ListEzRolloutResponse, error) {
	ezRollouts, total, err := x.uc.ListEzRollouts(ctx, &biz.ListEzRolloutRequest{
		EzRolloutCommonParams: biz.EzRolloutCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return &pb.ListEzRolloutResponse{
				List:  []*pb.EzRolloutInfo{},
				Total: 0,
			}, nil
		}
		return nil, err
	}

	list := make([]*pb.EzRolloutInfo, 0, len(ezRollouts))
	for _, ezRollout := range ezRollouts {
		if dto := x.convertDO2DTO(ezRollout); dto != nil {
			list = append(list, dto)
		}
	}

	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})

	return &pb.ListEzRolloutResponse{
		List:  list,
		Total: total,
	}, nil
}

// CreateEzRollout 创建EzRollout
func (x *EzRolloutService) CreateEzRollout(ctx context.Context, req *pb.CreateEzRolloutRequest) (*pb.CreateEzRolloutResponse, error) {
	err := x.uc.CreateEzRollout(ctx, &biz.CreateEzRolloutRequest{
		EzRolloutCommonParams: biz.EzRolloutCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Labels:           req.Labels,
		Annotations:      req.Annotations,
		Selector:         &metav1.LabelSelector{MatchLabels: req.Selector},
		ScaleUpMetrics:   x.convertMetricsV22DO(req.ScaleUpMetrics, true),
		ScaleDownMetrics: x.convertMetricsV22DO(req.ScaleDownMetrics, false),
		OnlineVersion:    req.OnlineVersion,
		EzOnlineScaler:   x.convertOnlineScaler2DO(req.ScaleUp, req.ScaleDown, req.MinReplicas, req.MaxReplicas),
		EzOfflineScaler:  &ezRolloutv1.EzOfflineScaler{Deadline: int64(req.OfflineDeadline), EnableScaleUp: req.EnableScaleUp},
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateEzRolloutResponse{}, nil
}

// UpdateEzRollout 更新EzRollout
func (x *EzRolloutService) UpdateEzRollout(ctx context.Context, req *pb.UpdateEzRolloutRequest) (*pb.UpdateEzRolloutResponse, error) {
	err := x.uc.UpdateEzRollout(ctx, &biz.CreateEzRolloutRequest{
		EzRolloutCommonParams: biz.EzRolloutCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
		Labels:           req.Labels,
		Annotations:      req.Annotations,
		Selector:         &metav1.LabelSelector{MatchLabels: req.Selector},
		ScaleUpMetrics:   x.convertMetricsV22DO(req.ScaleUpMetrics, true),
		ScaleDownMetrics: x.convertMetricsV22DO(req.ScaleDownMetrics, false),
		OnlineVersion:    req.OnlineVersion,
		EzOnlineScaler:   x.convertOnlineScaler2DO(req.ScaleUp, req.ScaleDown, req.MinReplicas, req.MaxReplicas),
		EzOfflineScaler:  &ezRolloutv1.EzOfflineScaler{Deadline: int64(req.OfflineDeadline), EnableScaleUp: req.EnableScaleUp},
	})

	if err != nil {
		return nil, err
	}
	return &pb.UpdateEzRolloutResponse{}, nil
}

// DeleteEzRollout 删除EzRollout
func (x *EzRolloutService) DeleteEzRollout(ctx context.Context, req *pb.DeleteEzRolloutRequest) (*pb.DeleteEzRolloutResponse, error) {
	err := x.uc.DeleteEzRollout(ctx, &biz.DeleteEzRolloutRequest{
		EzRolloutCommonParams: biz.EzRolloutCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})

	if err != nil {
		return nil, err
	}
	return &pb.DeleteEzRolloutResponse{}, nil
}

// CreateOrUpdateEzRolloutByYaml 创建或更新EzRollout
func (x *EzRolloutService) CreateOrUpdateEzRolloutByYaml(ctx context.Context, req *pb.CreateOrUpdateEzRolloutByYamlRequest) (*pb.CreateOrUpdateEzRolloutByYamlResponse, error) {
	err := x.uc.CreateOrUpdateEzRolloutByYaml(ctx, &biz.CreateOrUpdateEzRolloutByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateEzRolloutByYamlResponse{}, nil
}

// GetEzRolloutDetail  获取EzRollout详情
func (x *EzRolloutService) GetEzRolloutDetail(ctx context.Context, req *pb.EzRolloutDetailRequest) (*pb.EzRolloutDetailResponse, error) {
	ezRollout, err := x.uc.GetEzRollout(ctx, &biz.GetEzRolloutRequest{
		EzRolloutCommonParams: biz.EzRolloutCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		}})
	if err != nil {
		return nil, err
	}

	return &pb.EzRolloutDetailResponse{
		Detail: x.convertDO2DTO(ezRollout),
	}, nil
}
