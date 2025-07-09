package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type PodService struct {
	pb.UnimplementedPodServer
	uc   *biz.PodUseCase
	uf   *biz.UserFollowUseCase
	node *biz.NodeUseCase
}

func NewPodService(uc *biz.PodUseCase, uf *biz.UserFollowUseCase, node *biz.NodeUseCase) *PodService {
	return &PodService{
		uc:   uc,
		uf:   uf,
		node: node,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *PodService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Pod,
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
func (x *PodService) setFollowedStatusForPodList(ctx context.Context, clusterName string, items []*pb.PodItem) error {
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

func (x *PodService) setFollowedStatusForNamespacePodList(ctx context.Context, clusterName string, items []*pb.NamespacePodItem) error {
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

func (x *PodService) getNodeResourceUsage(ctx context.Context, clusterName string, Name string, podResourceUsage *pb.PodResourceUsage) (*pb.PodResourceUsage, error) {
	nodeItem, err := x.node.GetNodeForAPI(ctx, &biz.ListNodeRequest{
		ClusterName: clusterName,
		Name:        Name,
	})
	if err != nil {
		return nil, err
	}
	totalCpu := nodeItem.Allocatable.Cpu().MilliValue()
	totalMemory := nodeItem.Allocatable.Memory().Value()
	totalEphemeralStorage := nodeItem.Allocatable.StorageEphemeral().Value()
	podResourceUsage.TotalCpuRequest = uint32(totalCpu)
	podResourceUsage.TotalCpuLimit = uint32(totalCpu)
	podResourceUsage.TotalMemoryLimit = uint32(totalMemory)
	podResourceUsage.TotalMemoryRequest = uint32(totalMemory)
	podResourceUsage.TotalEphemeralStorageRequest = uint32(totalEphemeralStorage)
	podResourceUsage.TotalEphemeralStorageLimit = uint32(totalEphemeralStorage)
	return podResourceUsage, nil
}

func (x *PodService) convertDO2DTO(clusterName string, pod *corev1.Pod) *pb.PodItem {
	containers := make([]*corev1.Container, 0, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		containers = append(containers, &container)
	}
	// workload 工作负载类型
	var readyCount uint32
	containerStatuses := pod.Status.ContainerStatuses
	for _, containerStatus := range containerStatuses {
		if containerStatus.Ready {
			readyCount++
		}
	}
	//工作负载类型
	var (
		workload     string
		restartCount uint32
	)
	ownerReferences := pod.OwnerReferences
	if len(ownerReferences) > 0 {
		workload = ownerReferences[0].Kind
		if workload == "ReplicaSet" {
			// 重写为 Deployment
			workload = "Deployment"
		}
	} else {
		workload = ""
	}
	if pod.Status.ContainerStatuses != nil && len(pod.Status.ContainerStatuses) > 0 {
		restartCount = uint32(pod.Status.ContainerStatuses[0].RestartCount)
	}
	dto := &pb.PodItem{
		Name:          pod.Name,
		Namespace:     pod.Namespace,
		Status:        string(pod.Status.Phase),
		Containers:    containers,
		CreateTime:    uint64(pod.CreationTimestamp.UnixNano() / 1e6),
		RestartCount:  restartCount,
		Workload:      workload,
		ReadyCount:    readyCount,
		ResourceUsage: x.convertPodResourceUsage(pod),
	}
	ctx := context.Background()
	newResourceUsage, err := x.getNodeResourceUsage(ctx, clusterName, pod.Spec.NodeName, dto.ResourceUsage)
	if err != nil {
		return nil
	}
	dto.ResourceUsage = newResourceUsage
	return dto
}
func (x *PodService) convertPodConditions(conditions []corev1.PodCondition) []*pb.PodHealth {
	podConditions := make([]*pb.PodHealth, 0, len(conditions))
	for _, condition := range conditions {
		var (
			lastTransitionTime uint64
		)
		var lastProbeTime *uint64
		if !condition.LastProbeTime.IsZero() {
			tmp := uint64(condition.LastProbeTime.UnixNano() / 1e6)
			lastProbeTime = &tmp
		}
		if !condition.LastTransitionTime.IsZero() {
			lastTransitionTime = uint64(condition.LastTransitionTime.UnixNano() / 1e6)
		}
		podConditions = append(podConditions, &pb.PodHealth{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastProbeTime:      lastProbeTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return podConditions
}

func (x *PodService) convertDO2NamespaceItem(pod *corev1.Pod) *pb.NamespacePodItem {
	containers := make([]*corev1.Container, 0, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		containers = append(containers, &container)
	}
	if pod.APIVersion == "" {
		pod.APIVersion = "v1"
	}
	if pod.Kind == "" {
		pod.Kind = "Pod"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(pod)
	var restartCount uint32
	if pod.Status.ContainerStatuses != nil {
		restartCount = uint32(pod.Status.ContainerStatuses[0].RestartCount)
	}
	status, reason := GetPodStatus(pod)
	readyCount := GetPodReadyCount(pod)
	ownerReferences := pod.OwnerReferences
	var workload string
	if len(ownerReferences) > 0 {
		workload = ownerReferences[0].Kind
		if workload == "ReplicaSet" {
			// 重写为 Deployment
			workload = "Deployment"
		}
	} else {
		workload = ""
	}

	dto := &pb.NamespacePodItem{
		Name:           pod.Name,
		Namespace:      pod.Namespace,
		Status:         status,
		Containers:     containers,
		CreateTime:     uint64(pod.CreationTimestamp.UnixNano() / 1e6),
		RestartCount:   restartCount,
		ReadyCount:     readyCount,
		PodIp:          pod.Status.PodIP,
		NodeIp:         pod.Status.HostIP,
		Labels:         pod.Labels,
		Yaml:           yamlStr,
		ContainerCount: uint32(len(containers)),
		Reason:         reason,
		Annotations:    pod.Annotations,
		Conditions:     x.convertPodConditions(pod.Status.Conditions),
		Workload:       workload,
	}
	return dto
}

// GetPodResourceUsage 获取 Pod 的资源使用情况
func (x *PodService) convertPodResourceUsage(pod *corev1.Pod) *pb.PodResourceUsage {
	var (
		CPURequest, CPULimit, MemoryRequest, MemoryLimit, EphemeralStorageRequest, EphemeralStorageLimit int64
	)

	// 遍历所有容器
	for _, container := range pod.Spec.Containers {
		resources := container.Resources

		// CPU Request
		if cpuReq, ok := resources.Requests[corev1.ResourceCPU]; ok {
			CPURequest += cpuReq.MilliValue()
		}
		// CPU Limit
		if cpuLimit, ok := resources.Limits[corev1.ResourceCPU]; ok {
			CPULimit += cpuLimit.MilliValue()
		}
		// Memory Request 单位为bytes
		if memReq, ok := resources.Requests[corev1.ResourceMemory]; ok {
			MemoryRequest += memReq.Value()
		}
		// Memory Limit 单位为bytes
		if memLimit, ok := resources.Limits[corev1.ResourceMemory]; ok {
			MemoryLimit += memLimit.Value()
		}
		// Ephemeral Storage Request 单位为bytes
		if storageReq, ok := resources.Requests[corev1.ResourceEphemeralStorage]; ok {
			EphemeralStorageRequest += storageReq.Value()
		}
		// Ephemeral Storage Limit 单位为bytes
		if storageLimit, ok := resources.Limits[corev1.ResourceEphemeralStorage]; ok {
			EphemeralStorageLimit += storageLimit.Value()
		}
	}
	return &pb.PodResourceUsage{
		CpuRequest:              uint32(CPURequest),
		CpuLimit:                uint32(CPULimit),
		MemoryRequest:           uint32(MemoryRequest),
		MemoryLimit:             uint32(MemoryLimit),
		EphemeralStorageRequest: uint32(EphemeralStorageRequest),
		EphemeralStorageLimit:   uint32(EphemeralStorageLimit),
	}
}

// GetPodStatus 获取详细的 Pod 状态
func GetPodStatus(pod *corev1.Pod) (string, string) {
	// 如果 Pod 未就绪，查找具体原因
	if pod.Status.Phase != corev1.PodRunning {
		return string(pod.Status.Phase), ""
	}

	// 检查 conditions
	//for _, condition := range pod.Status.Conditions {
	//	if condition.Status != corev1.ConditionTrue {
	//		return fmt.Sprintf("%s:%s", condition.SecretType, condition.Reason)
	//	}
	//}

	// 检查容器状态
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			if containerStatus.State.Waiting != nil {
				return containerStatus.State.Waiting.Reason, containerStatus.State.Waiting.Message
			}
			if containerStatus.State.Terminated != nil {
				return containerStatus.State.Terminated.Reason, containerStatus.State.Terminated.Message
			}
		}
	}

	return string(corev1.PodRunning), ""
}

// GetPodReadyCount 查询pod 就绪容器数
func GetPodReadyCount(pod *corev1.Pod) uint32 {
	var readyCount uint32
	if pod.Status.Phase != corev1.PodRunning {
		return readyCount
	}
	if pod.Status.ContainerStatuses == nil {
		return readyCount
	}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			readyCount++
		}
	}
	return readyCount
}

func (x *PodService) convertControllerPod2DTO(pod corev1.Pod) *pb.ControllerPodItem {
	podConditions := make([]*pb.PodCondition, 0, len(pod.Status.Conditions))
	for _, condition := range pod.Status.Conditions {
		podConditions = append(podConditions, &pb.PodCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastProbeTime:      uint64(condition.LastProbeTime.UnixNano() / 1e6),
			LastTransitionTime: uint64(condition.LastTransitionTime.UnixNano() / 1e6),
			Reason:             &condition.Reason,
			Message:            &condition.Message,
		})
	}

	podContainers := make([]*pb.PodContainerDetail, 0, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		detail := &pb.PodContainerDetail{
			ContainerName: container.Name,
			Image:         container.Image,
		}
		var status corev1.ContainerStatus
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if container.Name == containerStatus.Name {
				if &containerStatus != nil {
					status = containerStatus
				}
			}
		}
		detail.Status = x.convertContainerStatus2DTO(status)
		if detail.Status != nil {
			detail.RestartCount = detail.Status.RestartCount
			detail.ContainerId = detail.Status.ContainerID
		}
		var resource pb.ResourceRequirements
		if container.Resources.Requests != nil {
			var requests = make(map[string]string)
			for resourceName, quantity := range container.Resources.Requests {
				requests[string(resourceName)] = quantity.String() // 转换 resource.Quantity 为字符串
			}
			resource.Requests = requests
		}
		if container.Resources.Limits != nil {
			var limits = make(map[string]string)
			for resourceName, quantity := range container.Resources.Limits {
				limits[string(resourceName)] = quantity.String()
			}
			resource.Limits = limits
		}
		detail.Resources = &resource

		for _, specContainer := range pod.Spec.Containers {
			if specContainer.Name == container.Name {
				// 镜像拉取策略
				detail.ImagePullPolicy = string(specContainer.ImagePullPolicy)

				// 端口
				ports := make([]*corev1.ContainerPort, 0, len(specContainer.Ports))
				for _, port := range specContainer.Ports {
					ports = append(ports, &port)
				}
				detail.Ports = ports
				// 环境变量
				env := make([]*corev1.EnvVar, 0, len(specContainer.Env))
				for _, e := range specContainer.Env {
					env = append(env, &e)
				}
				detail.Env = env
				// 卷
				volume := make([]*pb.ContainerVolumeMount, 0, len(specContainer.VolumeMounts))
				for _, v := range specContainer.VolumeMounts {
					volumeBytes, _ := yaml.Marshal(v)
					volume = append(volume, &pb.ContainerVolumeMount{
						ReadOnly: v.ReadOnly,
						Content:  string(volumeBytes),
					})
				}
				detail.VolumeMounts = volume
				// 存活探针
				if container.LivenessProbe != nil {
					LivenessProbeBytes, _ := yaml.Marshal(container.LivenessProbe)
					detail.LivenessProbe = string(LivenessProbeBytes)
				}
				// 就绪探针
				if container.ReadinessProbe != nil {
					ReadinessProbeBytes, _ := yaml.Marshal(container.ReadinessProbe)
					detail.ReadinessProbe = string(ReadinessProbeBytes)
				}
				// 启动探针
				if container.StartupProbe != nil {
					startProbeBytes, _ := yaml.Marshal(container.StartupProbe)
					detail.StartupProbe = string(startProbeBytes)
				}
				break
			}
		}
		podContainers = append(podContainers, detail)
	}
	var yamlStr string
	if pod.APIVersion == "" {
		pod.APIVersion = "v1"
	}
	if pod.Kind == "" {
		pod.Kind = "Pod"
	}
	yamlStr, _ = utils.ConvertResourceTOYaml(&pod)

	var isReady bool
	for _, condition := range podConditions {
		if condition.Type == "Ready" {
			if condition.Status == "True" {
				isReady = true
			} else {
				isReady = false
			}
		}
	}
	status, _ := GetPodStatus(&pod)
	return &pb.ControllerPodItem{
		Name:          pod.Name,
		PodIp:         pod.Status.PodIP,
		Status:        status,
		NodeIp:        pod.Status.HostIP,
		CreateTime:    uint64(pod.CreationTimestamp.UnixNano() / 1e6),
		PodConditions: podConditions,
		Containers:    podContainers,
		NodeName:      pod.Spec.NodeName,
		Yaml:          yamlStr,
		IsReady:       isReady,
		Namespace:     pod.Namespace,
	}
}

// convertContainerStatus2DTO 转换容器状态为 DTO
func (x *PodService) convertContainerStatus2DTO(status corev1.ContainerStatus) *pb.ContainerStatus {
	result := &pb.ContainerStatus{
		// 基本信息
		Name:         status.Name,
		Ready:        status.Ready,
		RestartCount: uint32(status.RestartCount),
		Image:        status.Image,
		ImageID:      status.ImageID,
		ContainerID:  status.ContainerID,
		Started:      status.Started,

		// 容器状态转换
		State:     convertContainerState(status.State),
		LastState: convertContainerState(status.LastTerminationState),

		// 资源配置
		Resources: convertResourceRequirements(status.Resources),
	}
	return result
}

func convertContainerState(state corev1.ContainerState) *pb.ContainerState {
	if state.Running != nil {
		return &pb.ContainerState{
			Running: &pb.ContainerStateRunning{
				StartedAt: uint64(state.Running.StartedAt.UnixNano() / 1e6),
			},
		}
	}
	if state.Terminated != nil {
		return &pb.ContainerState{
			Terminated: &pb.ContainerStateTerminated{
				ExitCode:   uint32(state.Terminated.ExitCode),
				Reason:     state.Terminated.Reason,
				Message:    state.Terminated.Message,
				StartedAt:  uint64(state.Terminated.StartedAt.UnixNano() / 1e6),
				FinishedAt: uint64(state.Terminated.FinishedAt.UnixNano() / 1e6),
			},
		}
	}
	if state.Waiting != nil {
		return &pb.ContainerState{
			Waiting: &pb.ContainerStateWaiting{
				Reason:  state.Waiting.Reason,
				Message: state.Waiting.Message,
			},
		}
	}
	return nil
}

// convertResourceRequirements 优化后的资源配置转换
func convertResourceRequirements(rr *corev1.ResourceRequirements) *pb.ResourceRequirements {
	if rr == nil || rr.Requests == nil && rr.Limits == nil {
		return nil
	}

	resource := &pb.ResourceRequirements{}

	// 转换资源请求
	if rr.Requests != nil {
		resource.Requests = convertResourceList(rr.Requests)
	}

	// 转换资源限制
	if rr.Limits != nil {
		resource.Limits = convertResourceList(rr.Limits)
	}

	return resource
}

// convertResourceList 转换资源列表
func convertResourceList(resources corev1.ResourceList) map[string]string {
	if len(resources) == 0 {
		return nil
	}

	result := make(map[string]string)
	for resourceName, quantity := range resources {
		if !quantity.IsZero() {
			result[string(resourceName)] = quantity.String()
		}
	}
	return result
}

func (x *PodService) ListPod(ctx context.Context, req *pb.ListPodRequest) (*pb.ListPodResponse, error) {
	pods, total, err := x.uc.ListPods(ctx, &biz.ListPodRequest{
		ClusterName: req.ClusterName,
		NodeName:    req.NodeName,
		Namespace:   req.Namespace,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ListAll:     utils.IntToBool(int(req.ListAll)),
		KeyWord:     req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.PodItem, 0, len(pods))
	for _, pod := range pods {
		list = append(list, x.convertDO2DTO(req.ClusterName, pod))
	}
	if err := x.setFollowedStatusForPodList(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListPodResponse{
		List:  list,
		Total: int32(total),
	}, nil
}

func (x *PodService) convertMetrics2DTO(metrics *biz.SidecarMetricResultList) *pb.SidecarMetricResultList {
	items := make([]*pb.SidecarMetric, 0, len(metrics.Items))
	for _, metric := range metrics.Items {
		dataPoints := make([]*pb.DataPoint, 0, len(metric.DataPoints))
		metricPoints := make([]*pb.MetricPoint, 0, len(metric.DataPoints))
		uids := make([]string, 0, len(metric.UIDs))
		for _, dataPoint := range metric.DataPoints {
			dataPoints = append(dataPoints, &pb.DataPoint{
				X: dataPoint.X,
				Y: dataPoint.Y,
			})
		}
		for _, metricPoint := range metric.MetricPoints {
			metricPoints = append(metricPoints, &pb.MetricPoint{
				Timestamp: metricPoint.Timestamp.Unix(),
				Value:     metricPoint.Value,
			})
		}
		for _, uid := range metric.UIDs {
			uids = append(uids, string(uid))
		}
		items = append(items, &pb.SidecarMetric{
			MetricName:   metric.MetricName,
			DataPoints:   dataPoints,
			MetricPoints: metricPoints,
			Uids:         uids,
		})
	}
	return &pb.SidecarMetricResultList{
		Items: items,
	}
}

func (x *PodService) GetPodCpuMetrics(ctx context.Context, req *pb.GetPodMetricsRequest) (*pb.SidecarMetricResultList, error) {
	metrics, err := x.uc.GetPodMetricsFromScraper(ctx, &biz.MetricsRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.PodName,
		},
		MetricName: "cpu",
		Time:       uint64(req.Time),
	})
	if err != nil {
		return &pb.SidecarMetricResultList{Items: make([]*pb.SidecarMetric, 0)}, err
	}

	return x.convertMetrics2DTO(&metrics), nil
}

func (x *PodService) GetPodMemoryMetrics(ctx context.Context, req *pb.GetPodMetricsRequest) (*pb.SidecarMetricResultList, error) {
	metrics, err := x.uc.GetPodMetricsFromScraper(ctx, &biz.MetricsRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.PodName,
		},
		MetricName: "memory",
		Time:       uint64(req.Time),
	})
	if err != nil {
		return &pb.SidecarMetricResultList{Items: make([]*pb.SidecarMetric, 0)}, err
	}

	return x.convertMetrics2DTO(&metrics), nil
}

func (x *PodService) GetPodContainerMetrics(ctx context.Context, req *pb.GetPodContainerMetricsRequest) (*pb.GetPodContainerMetricsResponse, error) {
	metrics, err := x.uc.GetPodContainerMetrics(ctx, &biz.PodContainerMetricsRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.PodName,
		},
		ContainerName: req.ContainerName,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetPodContainerMetricsResponse{
		Metrics: metrics,
	}, nil
}

func (x *PodService) DeletePod(ctx context.Context, req *pb.DeletePodRequest) (*pb.DeletePodResponse, error) {
	result, err := x.uc.DeletePod(ctx, &biz.DeletePodRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.PodName,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeletePodResponse{
		Result: result,
	}, nil
}

func (x *PodService) ListControllerPod(ctx context.Context, req *pb.ListControllerPodRequest) (*pb.ListControllerPodResponse, error) {
	pods, err := x.uc.ListControllerPod(ctx, req)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ControllerPodItem, 0, len(pods))
	for _, pod := range pods {
		item := x.convertControllerPod2DTO(*pod)
		metrics, _ := x.uc.GetPodContainerMetrics(
			ctx, &biz.PodContainerMetricsRequest{
				PodCommonParams: biz.PodCommonParams{
					ClusterName: req.ClusterName,
					Namespace:   req.Namespace,
					PodName:     pod.Name}})

		if metrics != nil {
			for _, m := range metrics.Containers {
				for _, c := range item.Containers {
					if m.Name == c.ContainerName {
						usage := &pb.ContainerUsage{}
						usage.Cpu = uint32(m.Usage.Cpu().MilliValue())
						usage.CpuUnit = "m"
						usage.Memory = uint32(m.Usage.Memory().Value() / 1024)
						usage.MemoryUnit = "ki"
						c.Usage = usage
					}
				}
			}
		}
		list = append(list, item)
	}
	return &pb.ListControllerPodResponse{
		List:  list,
		Total: uint32(len(list)),
	}, nil
}

func (x *PodService) DownloadPodLogs(ctx context.Context, req *pb.DownloadPodLogsRequest) (*pb.DownloadPodLogsResponse, error) {
	request := &biz.DownloadPodLogsRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.PodName,
		},
		ContainerName: req.ContainerName,
	}
	if req.StartTime != 0 {
		SinceTime := time.Unix(req.StartTime, 0)
		request.SinceTime = SinceTime
	}
	data, err := x.uc.DownloadPodLogs(ctx, request)

	if err != nil {
		return nil, err
	}
	return &pb.DownloadPodLogsResponse{
		Data: data,
	}, nil
}

func (x *PodService) ListPodByNamespace(ctx context.Context, req *pb.ListPodByNamespaceRequest) (*pb.ListPodByNamespaceResponse, error) {
	pods, total, err := x.uc.ListPodByNamespace(ctx, &biz.ListPodByNamespaceRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
		KeyWord:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.NamespacePodItem, 0, len(pods))
	for _, pod := range pods {
		list = append(list, x.convertDO2NamespaceItem(pod))
	}
	if err := x.setFollowedStatusForNamespacePodList(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListPodByNamespaceResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *PodService) GetNamespacePodDetail(ctx context.Context, req *pb.GetNamespacePodDetailRequest) (*pb.GetNamespacePodDetailResponse, error) {
	pod, err := x.uc.GetNamespacePodDetail(ctx, &biz.GetNamespacePodDetailRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	item := x.convertDO2NamespaceItem(pod)
	return &pb.GetNamespacePodDetailResponse{
		Detail: item,
	}, nil
}

func (x *PodService) CreateOrUpdatePodByYaml(ctx context.Context, req *pb.CreateOrUpdatePodByYamlRequest) (*pb.CreateOrUpdatePodByYamlResponse, error) {
	success, err := x.uc.CreateOrUpdatePodByYaml(ctx, &biz.CreateOrUpdatePodByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdatePodByYamlResponse{
		Success: success,
	}, nil
}

func (x *PodService) BatchDeletePod(ctx context.Context, req *pb.BatchDeletePodsRequest) (*pb.BatchDeletePodsResponse, error) {
	err := x.uc.BatchDeletePods(ctx, &biz.BatchDeletePodsRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
		},
		PodNames: req.PodNames,
	})
	if err != nil {
		return nil, err
	}
	return &pb.BatchDeletePodsResponse{}, nil
}

func (x *PodService) EvictPod(ctx context.Context, req *pb.EvictPodRequest) (*pb.EvictPodResponse, error) {
	// 1. 暂停调度
	_, err := x.node.CordonNode(ctx, &biz.HandleNodeRequest{
		ClusterName: req.ClusterName,
		Name:        req.NodeName,
		Operation:   pb.NodeOperation_NodeCordon,
	})
	if err != nil {
		return nil, err
	}
	// 2. 驱逐 Pod
	success, err := x.uc.EvictPod(ctx, &biz.EvictPodRequest{
		ClusterName:        req.ClusterName,
		Namespace:          req.Namespace,
		PodName:            req.PodName,
		NodeName:           req.NodeName,
		Force:              req.Force,
		GracePeriodSeconds: int64(req.GracePeriodSeconds),
	})
	if err != nil {
		return nil, err
	}
	return &pb.EvictPodResponse{
		Success: success,
	}, nil
}
