package biz

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"codo-cnmp/common/utils"
	"codo-cnmp/pb"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EventCommonParams 公共参数
type EventCommonParams struct {
	ClusterName    string // 集群名称
	Namespace      string // 命名空间
	ControllerName string // 控制器名称
	ControllerType string // 控制器类型
	PodName        string // Pod名称
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime time.Time // 开始时间
	EndTime   time.Time // 结束时间
}

// ListEventRequest 事件列表请求参数
type ListEventRequest struct {
	EventCommonParams
	TimeRange
	Page     uint32
	PageSize uint32
	Keyword  string
	ListAll  bool
}

type IEventUseCase interface {
	ListEvent(ctx context.Context, req *ListEventRequest) ([]*corev1.Event, uint32, error)
}

func NewIEventUseCase(x *EventUseCase) IEventUseCase {
	return x
}

type EventUseCase struct {
	Cluster IClusterUseCase
	Pod     IPodUseCase
	log     *log.Helper
}

func NewEventUseCase(cluster IClusterUseCase, pod IPodUseCase, logger log.Logger) *EventUseCase {
	return &EventUseCase{
		Cluster: cluster,
		Pod:     pod,
		log:     log.NewHelper(log.With(logger, "module", "biz/event")),
	}
}

// BuildEventListOptions 构建事件列表参数
func (x *EventUseCase) BuildEventListOptions(ctx context.Context, kind string, name string, limit uint32, continueToken string) *metav1.ListOptions {
	eventListOptions := metav1.ListOptions{
		Limit:    int64(limit),
		Continue: continueToken,
	}
	if name != "" && kind != "" {
		eventListOptions.FieldSelector = fmt.Sprintf("involvedObject.kind=%s,involvedObject.name=%s", kind, name)
	}
	return &eventListOptions
}

func (x *EventUseCase) FetchEvents(ctx context.Context, clientSet *kubernetes.Clientset, kind string, name string, req *ListEventRequest) ([]*corev1.Event, error) {
	var (
		events        []*corev1.Event
		continueToken string
	)

	eventListOptions := x.BuildEventListOptions(ctx, kind, name, req.PageSize, continueToken)

	for {
		// 获取事件列表
		EventList, err := clientSet.CoreV1().Events(req.Namespace).List(ctx, *eventListOptions)
		if err != nil {
			x.log.Warnf("获取事件列表失败: %v", err)
			return nil, fmt.Errorf("获取事件列表失败: %w", err)
		}
		// 过滤事件
		for _, event := range EventList.Items {
			// 过滤时间范围
			//if !req.TimeRange.StartTime.IsZero() && !event.EventTime.IsZero() && event.EventTime.Time.Before(req.TimeRange.StartTime) {
			//	continue
			//}
			//if !req.TimeRange.EndTime.IsZero() && !event.EventTime.IsZero() && event.EventTime.Time.After(req.TimeRange.EndTime) {
			//	continue
			//}
			events = append(events, &event)
		}
		// 如果没有下一页则退出
		if EventList.Continue == "" {
			break
		}
		// 设置下一页
		eventListOptions.Continue = EventList.Continue

	}
	return events, nil
}

// filterEvents 过滤事件列表
func filterEvents(events []*corev1.Event, keyword string) []*corev1.Event {
	// 模糊搜索事件等级、资源类型、资源名称、内容、消息
	var filteredEvents []*corev1.Event
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, event := range events {
		if utils.MatchString(pattern, event.Reason) ||
			utils.MatchString(pattern, event.Message) ||
			utils.MatchString(pattern, event.Type) ||
			utils.MatchString(pattern, event.InvolvedObject.Kind) ||
			utils.MatchString(pattern, event.InvolvedObject.Name) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents
}

func (x *EventUseCase) GetLabelSelector(ctx context.Context, req *ListEventRequest) (string, error) {
	var labelSelector string
	switch req.ControllerType {
	case pb.ListEventRequest_Deployment.String():
		clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}

		deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		labelSelector = metav1.FormatLabelSelector(deployment.Spec.Selector)
	case pb.ListEventRequest_StatefulSet.String():
		clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}
		statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		labelSelector = metav1.FormatLabelSelector(statefulSet.Spec.Selector)
	case pb.ListEventRequest_Pod.String():

		labelSelector = fmt.Sprintf("pod=%s", req.PodName)
	case pb.ListEventRequest_DaemonSet.String():
		clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}
		daemonSet, err := clientSet.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		labelSelector = metav1.FormatLabelSelector(daemonSet.Spec.Selector)
	case pb.ListEventRequest_Job.String():
	case pb.ListEventRequest_GameServerSet.String():
		clientSet, err := x.Cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}
		gameServerSet, err := clientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		labelSelector = metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: gameServerSet.Spec.GameServerTemplate.Labels,
		})
	case pb.ListEventRequest_CloneSet.String():
		clientSet, err := x.Cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}
		cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		labelSelector = metav1.FormatLabelSelector(cloneSet.Spec.Selector)
	case pb.ListControllerPodRequest_Hpa.String():
		clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}
		hpa, err := clientSet.AutoscalingV2beta2().HorizontalPodAutoscalers(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		switch hpa.Spec.ScaleTargetRef.Kind {
		case "Deployment":
		}
	case pb.ListEventRequest_Service.String():
		clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
		if err != nil {
			return "", err
		}
		service, err := clientSet.CoreV1().Services(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		if service.Spec.Selector == nil {
			return "", fmt.Errorf("service %s has no selector", req.ControllerName)
		}
		labelSelector = metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: service.Spec.Selector})
	case pb.ListEventRequest_Ingress.String():
		return "", nil
	default:
		return "", fmt.Errorf("未知的控制器类型: %s", req.ControllerType)
	}
	return labelSelector, nil
}

// isPodNameMatchController 检查 Pod 名称是否匹配控制器命名规则
func isPodNameMatchController(podName, controllerName, controllerType string) bool {
	// 基本前缀检查
	if !strings.HasPrefix(podName, controllerName+"-") {
		return false
	}

	// 特殊控制器的命名规则检查
	switch controllerType {
	case pb.ListEventRequest_StatefulSet.String():
		// StatefulSet Pod 名称格式: <statefulset-name>-<ordinal>
		return regexp.MustCompile(fmt.Sprintf("^%s-[0-9]+$", regexp.QuoteMeta(controllerName))).MatchString(podName)
	case pb.ListEventRequest_Deployment.String(),
		pb.ListEventRequest_DaemonSet.String(),
		pb.ListEventRequest_Job.String(),
		pb.ListEventRequest_CloneSet.String(),
		pb.ListControllerPodRequest_Hpa.String(),
		pb.ListEventRequest_GameServerSet.String():
		// 这些控制器都使用 <controller-name>-<hash> 格式
		return true
	default:
		return false
	}
}

// ListEventsForControllerV2 获取控制器相关事件,
func (x *EventUseCase) ListEventsForControllerV2(ctx context.Context, clientSet *kubernetes.Clientset, req *ListEventRequest, podsList []*corev1.Pod) ([]*corev1.Event, error) {

	// 获取命名空间下所有最近事件
	allEvents, err := x.FetchEvents(ctx, clientSet, "", "", req)
	if err != nil {
		return nil, err
	}
	currentPodNames := make(map[string]struct{})
	if len(podsList) > 0 {
		for _, pod := range podsList {
			currentPodNames[pod.Name] = struct{}{}
		}
	}

	// 过滤相关事件
	var filteredEvents []*corev1.Event
	for _, event := range allEvents {
		involved := event.InvolvedObject

		// 1. 检查控制器本身的事件
		if involved.Kind == req.ControllerType && involved.Name == req.ControllerName {
			filteredEvents = append(filteredEvents, event)
			continue
		}

		// 2. 检查 Pod 事件
		if involved.Kind == "Pod" {
			// 检查当前存在的 Pod
			if _, exists := currentPodNames[involved.Name]; exists {
				filteredEvents = append(filteredEvents, event)
				continue
			}

			// 检查已删除的 Pod
			if isPodNameMatchController(involved.Name, req.ControllerName, req.ControllerType) {
				filteredEvents = append(filteredEvents, event)
			}
		}

		// 3. 检查中间资源事件
		switch req.ControllerType {
		case pb.ListEventRequest_Deployment.String():
			if involved.Kind == "ReplicaSet" {
				rs, err := clientSet.AppsV1().ReplicaSets(req.Namespace).Get(ctx, involved.Name, metav1.GetOptions{})
				if err == nil {
					for _, owner := range rs.OwnerReferences {
						if owner.Kind == "Deployment" && owner.Name == req.ControllerName {
							filteredEvents = append(filteredEvents, event)
							break
						}
					}
				}
			}
		case pb.ListEventRequest_StatefulSet.String():
			if involved.Kind == "PersistentVolumeClaim" && strings.HasPrefix(involved.Name, req.ControllerName+"-") {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}
	return filteredEvents, nil
}

func (x *EventUseCase) ListEventsForController(ctx context.Context, clientSet *kubernetes.Clientset, req *ListEventRequest, podsList []*corev1.Pod) ([]*corev1.Event, error) {
	var (
		events []*corev1.Event
	)
	eg := errgroup.WithContext(ctx)
	eventsCh := make(chan []*corev1.Event, len(podsList)+1) // +1 是为了包括控制器事件

	// Fetch controller events
	eg.Go(func(ctx context.Context) error {
		controllerEvents, err := x.FetchEvents(ctx, clientSet, req.ControllerType, req.ControllerName, req)
		if err != nil {
			return err
		}
		eventsCh <- controllerEvents
		return nil
	})

	// Fetch pod events
	for _, pod := range podsList {
		p := pod
		eg.Go(func(ctx context.Context) error {
			podReq := *req
			podReq.PodName = p.Name
			podEvents, err := x.FetchEvents(ctx, clientSet, "Pod", p.Name, &podReq)
			if err != nil {
				return err
			}
			eventsCh <- podEvents
			return nil
		})
	}
	// 等待所有事件获取完成
	go func() {
		eg.Wait()
		close(eventsCh)
	}()

	for ev := range eventsCh {
		events = append(events, ev...)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return events, nil
}

// getEventTime 获取事件的排序时间
func getEventTime(event *corev1.Event) time.Time {
	// 优先使用 EventTime（v1.19+的新字段）
	if !event.EventTime.IsZero() {
		return event.EventTime.Time
	}
	// 其次使用 LastTimestamp
	if !event.LastTimestamp.IsZero() {
		return event.LastTimestamp.Time
	}
	// 最后使用 FirstTimestamp
	return event.FirstTimestamp.Time
}

func (x *EventUseCase) ListEvent(ctx context.Context, req *ListEventRequest) ([]*corev1.Event, uint32, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	var events []*corev1.Event
	switch req.ControllerType {
	case pb.ListEventRequest_Pod.String():
		events, err = x.FetchEvents(ctx, clientSet, req.ControllerType, req.ControllerName, req)
		if err != nil {
			return nil, 0, err
		}
	case pb.ListEventRequest_Ingress.String():
		events, err = x.ListEventsForControllerV2(ctx, clientSet, req, nil)
		if err != nil {
			return nil, 0, err
		}
	case pb.ListEventRequest_Deployment.String(),
		pb.ListEventRequest_StatefulSet.String(),
		pb.ListEventRequest_DaemonSet.String(),
		pb.ListControllerPodRequest_Hpa.String(),
		pb.ListEventRequest_Job.String(),
		pb.ListControllerPodRequest_GameServerSet.String(),
		pb.ListControllerPodRequest_CloneSet.String(),
		pb.ListControllerPodRequest_Service.String():
		podsList, err := x.Pod.ListControllerPod(ctx, &pb.ListControllerPodRequest{
			Namespace:      req.Namespace,
			ClusterName:    req.ClusterName,
			ControllerType: pb.ListControllerPodRequest_ControllerType(pb.ListControllerPodRequest_ControllerType_value[req.ControllerType]),
			ControllerName: req.ControllerName,
			ListAll:        1,
		})
		if err != nil {
			return nil, 0, err
		}

		if req.TimeRange.StartTime.IsZero() {
			req.TimeRange.StartTime = time.Now().Add(-time.Hour)
		}
		if req.TimeRange.EndTime.IsZero() {
			req.TimeRange.EndTime = time.Now()
		}

		events, err = x.ListEventsForControllerV2(ctx, clientSet, req, podsList)
		if err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, fmt.Errorf("未知的控制器类型: %s", req.ControllerType)

	}
	filteredEvents := filterEvents(events, req.Keyword)

	//sort.Slice(filteredEvents, func(i, j int) bool {
	//	return filteredEvents[i].LastTimestamp.Time.After(filteredEvents[j].LastTimestamp.Time)
	//})

	if req.ListAll {
		return filteredEvents, uint32(len(filteredEvents)), nil
	}

	paginatedEvents, total := utils.K8sPaginate(filteredEvents, req.Page, req.PageSize)
	return paginatedEvents, total, nil
}
