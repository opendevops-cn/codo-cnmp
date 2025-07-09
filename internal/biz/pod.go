package biz

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"codo-cnmp/common/xerr"
	"github.com/go-kratos/kratos/v2/errors"
	gamekruiseiov1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	e "github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/yaml"

	"strings"

	"codo-cnmp/common/utils"
	"codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	schema "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type PodCommonParams struct {
	ClusterName string `json:"cluster_name"` // 集群名称
	Namespace   string `json:"namespace"`    // 命名空间
	PodName     string `json:"pod_name"`     // Pod 名称
}

// ListPodRequest 查询Pod列表请求参数
type ListPodRequest struct {
	ClusterName    string // 集群名称
	NodeName       string // 节点名称
	DeploymentName string // Deployment名称
	Namespace      string // 命名空间
	KeyWord        string // 关键词
	Page           uint32 // 页码
	PageSize       uint32 // 页大小
	ListAll        bool   // 是否查询所有Pod
}

// EvictPodRequest 驱逐Pod请求参数
type EvictPodRequest struct {
	ClusterName        string // 集群名称
	Namespace          string // 命名空间
	PodName            string // Pod名称
	NodeName           string // 节点名称
	GracePeriodSeconds int64  // 宽限期
	Force              bool   // 是否强制删除
}

// MetricsRequest 从 metrics-scraper 获取 Pod 监控数据请求参数
type MetricsRequest struct {
	PodCommonParams
	MetricName string // 指标名称
	Time       uint64 // 时间戳
}

// PodContainerMetricsRequest 获取 Pod 容器监控数据请求参数
type PodContainerMetricsRequest struct {
	PodCommonParams
	ContainerName string // 容器名称
}

// DeletePodRequest 删除Pod请求参数
type DeletePodRequest struct {
	PodCommonParams
}

// BatchDeletePodsRequest 批量删除Pods请求参数
type BatchDeletePodsRequest struct {
	PodCommonParams
	PodNames []string // Pod名称列表
}

// TailPodLogRequest 追踪Pod日志请求参数
type TailPodLogRequest struct {
	PodCommonParams
	ContainerName string `json:"container_name"` // 容器名称
	TailLines     uint32 `json:"tail_lines"`     // 追踪行数
}

// ExecPodRequest 进入Pod请求参数
type ExecPodRequest struct {
	PodCommonParams
	ContainerName string `json:"container_name"` // 容器名称
	Command       string `json:"command"`        // 命令
	Shell         string `json:"shell"`          // tty shell类型
}

// DownloadPodLogsRequest 下载Pod日志请求参数
type DownloadPodLogsRequest struct {
	PodCommonParams
	ContainerName string    `json:"container_name"` // 容器名称
	SinceTime     time.Time `json:"since_time"`     // 开始时间
}

// ListPodsByCloneSetRequest 查询CloneSet关联的Pods列表请求参数
type ListPodsByCloneSetRequest struct {
	PodCommonParams
	CloneSetName string // CloneSet名称
}

// ListControllerPodRequest 查询控制器关联的Pods列表请求参数
type ListControllerPodRequest struct {
	PodCommonParams
	ControllerType pb.ListControllerPodRequest_ControllerType // 控制器类型
	ControllerName string                                     // 控制器名称
	Page           uint32                                     // 页码
	PageSize       uint32                                     // 页大小
	ListAll        bool                                       // 是否查询所有Pod
	KeyWord        string                                     // 关键词
}

// ListPodByNamespaceRequest 查询指定命名空间下的Pod列表请求参数
type ListPodByNamespaceRequest struct {
	PodCommonParams
	Page     uint32 // 页码
	PageSize uint32 // 页大小
	ListAll  bool   // 是否查询所有Pod
	KeyWord  string // 关键词
}

// GetNamespacePodDetailRequest 获取命名空间pod详情请求参数
type GetNamespacePodDetailRequest struct {
	PodCommonParams
}

// CreateOrUpdatePodByYamlRequest 通过YAML创建或更新Pod请求参数
type CreateOrUpdatePodByYamlRequest struct {
	ClusterName string
	Yaml        string // YAML内容
}

type IPodUseCase interface {
	// ListPods ListPods 查询Pods列表
	ListPods(ctx context.Context, req *ListPodRequest) ([]*corev1.Pod, uint32, error)
	// EvictPod EvictPod 驱逐Pod
	EvictPod(ctx context.Context, req *EvictPodRequest) (bool, error)
	// ListPodsByDeployment ListPodsByDeployment 查询Deployment关联的Pods列表
	ListPodsByDeployment(ctx context.Context, req *ListPodRequest) ([]*corev1.Pod, error)
	// GetPodMetricsFromScraper 从 metrics-scraper 获取 Pod 监控数据
	GetPodMetricsFromScraper(ctx context.Context, req *MetricsRequest) (SidecarMetricResultList, error)
	// GetPodContainerMetrics 获取 Pod 容器监控数据
	GetPodContainerMetrics(ctx context.Context, req *PodContainerMetricsRequest) (*v1beta1.PodMetrics, error)
	// DeletePod 删除Pod
	DeletePod(ctx context.Context, req *DeletePodRequest) (bool, error)
	// BatchDeletePods 批量删除Pods
	BatchDeletePods(ctx context.Context, req *BatchDeletePodsRequest) error
	// TailPodLogs 追踪Pod日志
	TailPodLogs(ctx context.Context, output chan string, req *TailPodLogRequest) error
	// ListControllerPod 查询控制器关联的Pods列表
	ListControllerPod(ctx context.Context, req *pb.ListControllerPodRequest) ([]*corev1.Pod, error)
	// ExecPod 进入Pod
	ExecPod(ctx context.Context, req *ExecPodRequest) (remotecommand.Executor, error)
	// DownloadPodLogs 下载Pod日志
	DownloadPodLogs(ctx context.Context, req *DownloadPodLogsRequest) (string, error)
	// ListPodsByCloneSet 查询CloneSet关联的Pods列表
	ListPodsByCloneSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error)
	// ListPodsByGameServerSet 查询GameServerSet关联的Pods列表
	ListPodsByGameServerSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error)
	// ListPodsByStatefulSet 查询StatefulSet关联的Pods列表
	ListPodsByStatefulSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error)
	// ListPodsByDaemonSet 查询DaemonSet关联的Pods列表
	ListPodsByDaemonSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error)
	// ListPodByNamespace 查询指定命名空间下的Pod列表
	ListPodByNamespace(ctx context.Context, req *ListPodByNamespaceRequest) ([]*corev1.Pod, uint32, error)
	// GetNamespacePodDetail 获取命名空间pod详情
	GetNamespacePodDetail(ctx context.Context, req *GetNamespacePodDetailRequest) (*corev1.Pod, error)
	// GetWorkloadByPod 获取Pod关联的工作负载
	GetWorkloadByPod(ctx context.Context, pod *corev1.Pod, req *PodCommonParams) (string, string, error)
	// CreateOrUpdatePodByYaml 通过YAML创建或更新Pod
	CreateOrUpdatePodByYaml(ctx context.Context, req *CreateOrUpdatePodByYamlRequest) (bool, error)
	// ListPodsBySideCarSet 查询SidecarSet关联的Pods列表
	ListPodsBySideCarSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error)
	// GetPod 查询Pod
	GetPod(ctx context.Context, req *ListControllerPodRequest) (*corev1.Pod, error)
	// 查询服务关联的Pods列表
	ListPodByService(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error)
}

func (x *PodUseCase) DeletePod(ctx context.Context, req *DeletePodRequest) (bool, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建k8s client失败: %v", err)
		return false, fmt.Errorf("创建k8s client失败: %w", err)
	}
	// 使用后台删除策略删除Pod
	deletePolicy := metav1.DeletePropagationBackground
	err = clientSet.CoreV1().Pods(req.Namespace).Delete(ctx, req.PodName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除Pod失败: %v", err)
		return false, fmt.Errorf("删除Pod失败: %w", err)
	}
	return true, nil
}

// ListControllerPod 查询控制器关联的Pods列表
func (x *PodUseCase) ListControllerPod(ctx context.Context, req *pb.ListControllerPodRequest) ([]*corev1.Pod, error) {
	switch req.ControllerType {
	case pb.ListControllerPodRequest_Deployment:
		return x.ListPodsByDeployment(ctx, &ListPodRequest{
			ClusterName:    req.ClusterName,
			Namespace:      req.Namespace,
			DeploymentName: req.ControllerName,
			Page:           req.Page,
			PageSize:       req.PageSize,
			ListAll:        utils.IntToBool(int(req.ListAll)),
		})
	case pb.ListControllerPodRequest_CloneSet:
		return x.ListPodsByCloneSet(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
	case pb.ListControllerPodRequest_GameServerSet:
		return x.ListPodsByGameServerSet(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
	case pb.ListControllerPodRequest_StatefulSet:
		return x.ListPodsByStatefulSet(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
	case pb.ListControllerPodRequest_DaemonSet:
		return x.ListPodsByDaemonSet(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
	case pb.ListControllerPodRequest_Job:
		return x.ListPodsByJob(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
	case pb.ListControllerPodRequest_Hpa:
		return []*corev1.Pod{}, nil
	case pb.ListControllerPodRequest_SideCarSet:
		return x.ListPodsBySideCarSet(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
	case pb.ListControllerPodRequest_Pod:
		pod, err := x.GetPod(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})
		return []*corev1.Pod{pod}, err

	case pb.ListControllerPodRequest_Service:
		return x.ListPodByService(ctx, &ListControllerPodRequest{
			PodCommonParams: PodCommonParams{
				ClusterName: req.ClusterName,
				Namespace:   req.Namespace,
			},
			ControllerName: req.ControllerName,
		})

	default:
		return nil, fmt.Errorf("未知的控制器类型: %v", req.ControllerType)
	}
}

type PodUseCase struct {
	Cluster IClusterUseCase
	log     *log.Helper
}

// ListPodByService implements IPodUseCase.
func (x *PodUseCase) ListPodByService(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	service, err := clientSet.CoreV1().Services(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Service失败: %v", err)
		return nil, fmt.Errorf("查询Service失败: %w", err)
	}
	// selector := labels.Set(service.Spec.Selector).AsSelector()
	labelSelector := labels.SelectorFromSet(service.Spec.Selector).String()
	if labelSelector == "" {
		return nil, nil
	}
	// TODO: 分页查询
	pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Pods失败: %v", err)
		return nil, fmt.Errorf("查询Pods失败: %w", err)
	}
	var result []*corev1.Pod
	for _, pod := range pods.Items {
		result = append(result, &pod)
	}
	return result, nil
}

func (x *PodUseCase) GetPod(ctx context.Context, req *ListControllerPodRequest) (*corev1.Pod, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	pod, err := clientSet.CoreV1().Pods(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Pod失败: %v", err)
		return nil, fmt.Errorf("查询Pod失败: %w", err)
	}
	return pod, nil
}

func (x *PodUseCase) ListPodsBySideCarSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	kruiseClientSet, err := x.Cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	sideCarSets, err := kruiseClientSet.AppsV1alpha1().SidecarSets().Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询SideCarSets失败: %v", err)
		return nil, fmt.Errorf("查询SideCarSets失败: %w", err)
	}

	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	selector, err := metav1.LabelSelectorAsSelector(sideCarSets.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("解析LabelSelector失败: %v", err)
	}
	var (
		limit         = int64(req.PageSize)
		continueToken = ""
		result        []*corev1.Pod
	)
	for {
		podListOptions := metav1.ListOptions{
			LabelSelector: selector.String(),
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(metav1.NamespaceAll).List(ctx, podListOptions)
		if err != nil {
			return nil, fmt.Errorf("查询pods失败: %w", err)
		}
		for _, pod := range pods.Items {
			result = append(result, &pod)
		}
		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}
	return result, nil

}

func (x *PodUseCase) BatchDeletePods(ctx context.Context, req *BatchDeletePodsRequest) error {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建k8s client失败: %v", err)
		return fmt.Errorf("创建k8s client失败: %w", err)
	}
	// 使用后台删除策略删除Pod
	deletePolicy := metav1.DeletePropagationBackground
	for _, podName := range req.PodNames {
		err = clientSet.CoreV1().Pods(req.Namespace).Delete(ctx, podName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if err != nil {
			x.log.WithContext(ctx).Errorf("删除Pod失败: %v", err)
			return fmt.Errorf("删除Pod失败: %w", err)
		}
	}
	return nil
}

func (x *PodUseCase) CreateOrUpdatePodByYaml(ctx context.Context, req *CreateOrUpdatePodByYamlRequest) (bool, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, err
	}
	pod := &corev1.Pod{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&pod); err != nil {
		return false, fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
	}
	// 获取当前的Pod对象
	_, err = clientSet.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 创建Pod
			_, err = clientSet.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
			if err != nil {
				return false, fmt.Errorf("创建Pod失败: %w", err)
			}
			return true, nil
		}
		return false, fmt.Errorf("查询Pod失败: %w", err)
	}
	// 更新pod
	_, err = clientSet.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新pod失败: %v", err)
		return false, fmt.Errorf("更新pod失败: %w", err)
	}
	return true, nil
}

func (x *PodUseCase) GetWorkloadByPod(ctx context.Context, pod *corev1.Pod, req *PodCommonParams) (string, string, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return "", "", err
	}
	// 取第一个 OwnerReference
	ownerRefs := pod.OwnerReferences
	if len(ownerRefs) == 0 {
		return "", "", nil
	}
	ownerRef := ownerRefs[0]

	switch ownerRef.Kind {
	case "ReplicaSet":
		// 获取 ReplicaSet，检查是否关联 Deployment
		rs, err := clientSet.AppsV1().ReplicaSets(pod.Namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
		if err != nil {
			return "", "", fmt.Errorf("查询ReplicaSet失败: %w", err)
		}
		if len(rs.OwnerReferences) > 0 && rs.OwnerReferences[0].Kind == "Deployment" {
			return rs.OwnerReferences[0].Kind, rs.OwnerReferences[0].Name, nil
		}
		return ownerRef.Kind, ownerRef.Name, nil

	case "StatefulSet":
		// 检查labels是否包含“game.kruise.io”，, 如果包含，则认为是GameServerSet，否则认为是StatefulSet
		pattern := "(?i).*" + regexp.QuoteMeta("game.kruise.io") + ".*"
		ok := utils.MatchLabels(pattern, pod.Labels)
		if ok {
			return "GameServerSet", ownerRef.Name, nil
		}
		return "StatefulSet", ownerRef.Name, nil
	case "DaemonSet", "ControllerRevision", "CloneSet":
		// 直接返回 StatefulSet, DaemonSet, 或 ControllerRevision
		return ownerRef.Kind, ownerRef.Name, nil

	default:
		return "", "", fmt.Errorf("未知的控制器类型: %v", ownerRef.Kind)
	}
}

func (x *PodUseCase) GetNamespacePodDetail(ctx context.Context, req *GetNamespacePodDetailRequest) (*corev1.Pod, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	pod, err := clientSet.CoreV1().Pods(req.Namespace).Get(ctx, req.PodName, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Pod失败: %v", err)
		return nil, fmt.Errorf("查询Pod失败: %w", err)
	}
	return pod, nil
}

func (x *PodUseCase) ListPodByNamespace(ctx context.Context, req *ListPodByNamespaceRequest) ([]*corev1.Pod, uint32, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var (
		allFilteredPods = make([]*corev1.Pod, 0)
		continueToken   string
		limit           = int64(req.PageSize)
	)

	for {
		podListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询pods失败: %v", err)
			return nil, 0, fmt.Errorf("查询pods失败: %w", err)
		}
		filteredPods := filterPods(pods, req.KeyWord)
		allFilteredPods = append(allFilteredPods, filteredPods...)

		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredPods, 0, nil
	}
	if len(allFilteredPods) == 0 {
		return nil, 0, nil
	}
	paginatedPods, total := utils.K8sPaginate(allFilteredPods, req.Page, req.PageSize)
	return paginatedPods, total, nil

}

func NewPodUseCase(cluster IClusterUseCase, logger log.Logger) *PodUseCase {
	return &PodUseCase{
		Cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/pod")),
	}
}

func NewIPodUseCase(x *PodUseCase) IPodUseCase {
	return x
}

func (x *PodUseCase) ListPods(ctx context.Context, req *ListPodRequest) ([]*corev1.Pod, uint32, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var (
		allFilteredPods = make([]*corev1.Pod, 0)
		continueToken   string
		limit           = int64(req.PageSize)
		fieldSelector   = fmt.Sprintf("spec.nodeName=%s", req.NodeName)
	)

	for {
		podListOptions := metav1.ListOptions{
			FieldSelector: fieldSelector,
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询pods失败: %v", err)
			return nil, 0, fmt.Errorf("查询pods失败: %w", err)
		}
		filteredPods := filterPods(pods, req.KeyWord)
		allFilteredPods = append(allFilteredPods, filteredPods...)

		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredPods, uint32(len(allFilteredPods)), nil
	}
	if len(allFilteredPods) == 0 {
		return allFilteredPods, 0, nil
	}
	paginatedPods, total := utils.K8sPaginate(allFilteredPods, req.Page, req.PageSize)
	return paginatedPods, total, nil
}

// 检查关键词是否在 Pod 名称、状态、命名空间、镜像、工作负载名称中
func filterPods(pods *corev1.PodList, keyword string) []*corev1.Pod {
	var result []*corev1.Pod
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, pod := range pods.Items {
		if utils.MatchString(pattern, pod.Name) ||
			utils.MatchString(pattern, string(pod.Status.Phase)) ||
			utils.MatchString(pattern, pod.Namespace) ||
			utils.MatchWorkloadName(pattern, pod.GetOwnerReferences()) ||
			utils.MatchContainerImages(pattern, pod.Spec.Containers) {
			result = append(result, &pod)
		}
	}
	return result
}

func (x *PodUseCase) EvictPod(ctx context.Context, req *EvictPodRequest) (bool, error) {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, fmt.Errorf("创建k8s client失败: %w", err)
	}
	// 设置默认的宽限期和删除策略
	var gracePeriod int64 = 30 // 默认 30s
	var deletePolicy metav1.DeletionPropagation = metav1.DeletePropagationBackground

	// 如果用户指定了宽限期，使用用户配置
	if req.GracePeriodSeconds >= 0 {
		gracePeriod = req.GracePeriodSeconds
	}

	// 如果强制删除，则设置宽限期为 0
	if req.Force {
		//gracePeriod = 0
		deletePolicy = metav1.DeletePropagationForeground
	}

	eviction := &policyv1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1",
			Kind:       "Eviction",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.PodName,
			Namespace: req.Namespace,
		},
		DeleteOptions: &metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriod,
			PropagationPolicy:  &deletePolicy,
		},
	}
	err = clientSet.CoreV1().Pods(req.Namespace).EvictV1(ctx, eviction)
	if err != nil {
		x.log.WithContext(ctx).Errorf("驱逐Pod失败: %v", err)
		return false, fmt.Errorf("驱逐Pod失败: %w", err)
	}
	return true, nil
}

// ListPodsByDeployment 查询Deployment关联的Pods列表
func (x *PodUseCase) ListPodsByDeployment(ctx context.Context, req *ListPodRequest) ([]*corev1.Pod, error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var (
		limit         = int64(req.PageSize)
		continueToken = ""
		result        []*corev1.Pod
	)
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.DeploymentName, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Deployment失败: %v", err)
		return nil, fmt.Errorf("查询Deployment失败: %w", err)
	}
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("解析LabelSelector失败: %v", err)
	}
	for {
		podListOptions := metav1.ListOptions{
			LabelSelector: selector.String(),
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询pods失败: %v", err)
			return nil, fmt.Errorf("查询pods失败: %w", err)
		}
		for _, pod := range pods.Items {
			result = append(result, &pod)
		}
		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}
	return result, nil
}

type PodMetricsResponse struct {
	// Define the structure based on the expected response from metrics-scraper
	CPUUsage    string `json:"cpu_usage"`
	MemoryUsage string `json:"memory_usage"`
}

type DataPoints []DataPoint

type DataPoint struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     uint64    `json:"value"`
}

// SidecarMetric is a format of data used by our sidecar. This is also the format of data that is being sent by backend API.
type SidecarMetric struct {
	// DataPoints is a list of X, Y int64 data points, sorted by X.
	DataPoints `json:"dataPoints"`
	// MetricPoints is a list of value, timestamp metrics used for sparklines on a pod list page.
	MetricPoints []MetricPoint `json:"metricPoints"`
	// MetricName is the name of metric stored in this struct.
	MetricName string `json:"metricName"`
	// Label stores information about identity of resources (UIDS) described by this metric.
	UIDs []types.UID `json:"uids"`
}

type SidecarMetricResultList struct {
	Items []SidecarMetric `json:"items"`
}

func (x *PodUseCase) GetPodMetricsFromScraper(ctx context.Context, req *MetricsRequest) (SidecarMetricResultList, error) {
	result := SidecarMetricResultList{}
	suffix := fmt.Sprintf("/api/v1/dashboard/namespaces/%s/pod-list/%s/metrics/%s/%d", req.Namespace, req.PodName, req.MetricName, req.Time)
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建k8s client失败: %v", err)
		return result, err
	}
	// 构造请求 URL [默认 metrics-scraper 服务暴露在 kube-system 命名空间的 8000 端口]
	url := clientSet.CoreV1().RESTClient().Get().
		Namespace("kube-system").
		//Namespace("kuboard").
		Resource("services").
		Name("metrics-scraper:8000").
		SubResource("proxy").
		Suffix(suffix). // API 路径
		URL().String()
	// 使用 RESTClient 执行 HTTP 请求
	resp, err := clientSet.RESTClient().Get().RequestURI(url).DoRaw(ctx)
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return result, e.Wrapf(xerr.NewErrCodeMsg(xerr.ErrNotAllowed, "没有权限查询指标数据"), "没有权限查询指标数据")
		}
		if k8serrors.IsNotFound(err) {
			return result, e.Wrapf(xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请在集群中部署metrics-scraper组件"), "metrics-scraper组件未找到")
		}
		x.log.WithContext(ctx).Errorf("查询指标数据失败: %v", err)
		return result, e.Wrapf(xerr.NewErrCodeMsg(xerr.ErrResourceGetError, "查询指标数据失败"), "查询指标数据失败: %v", err)

	}
	if err := json.Unmarshal(resp, &result); err != nil {
		x.log.WithContext(ctx).Errorf("解析指标数据失败: %v", err)
		return result, err
	}
	return result, nil
}

func (x *PodUseCase) GetPodContainerMetrics(ctx context.Context, req *PodContainerMetricsRequest) (*v1beta1.PodMetrics, error) {
	var result *v1beta1.PodMetrics
	clientSet, err := x.Cluster.GetMetricsClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return result, fmt.Errorf("创建k8s client失败: %w", err)
	}
	// 查询 Pod Container 监控数据
	// 1. 查询 Pod 监控数据
	metrics, err := clientSet.MetricsV1beta1().PodMetricses(req.Namespace).Get(ctx, req.PodName, metav1.GetOptions{})
	if err != nil {
		return result, fmt.Errorf("查询Pod监控数据失败: %w", err)
	}
	return metrics, nil
}

func (x *PodUseCase) TailPodLogs(ctx context.Context, output chan string, req *TailPodLogRequest) error {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	tailLines := int64(req.TailLines)
	podLogOptions := corev1.PodLogOptions{
		Container: req.ContainerName,
		TailLines: &tailLines,
		Follow:    true,
		Previous:  false,
	}
	reqLogs := clientSet.CoreV1().Pods(req.Namespace).GetLogs(req.PodName, &podLogOptions)
	podLogs, err := reqLogs.Stream(ctx)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取Pod日志失败: %v", err)
		return fmt.Errorf("获取Pod日志失败: %w", err)
	}
	defer podLogs.Close()
	reader := bufio.NewReader(podLogs)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				return nil
			} else if err != nil {
				return fmt.Errorf("读取Pod日志失败: %v", err)
			}
			select {
			case <-ctx.Done():
				return nil
			case output <- line:
			}
		}
	}
}

func (x *PodUseCase) ExecPod(ctx context.Context, req *ExecPodRequest) (remotecommand.Executor, error) {
	// 获取集群的认证配置信息
	cfg, err := x.Cluster.GetClientConfigByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}

	// 构建 Kubernetes 的 REST 配置信息
	config, err := cfg.BuildRestConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	// 初始化 Pod Exec 的请求，这里只进行连接，不执行任何命令
	execReq := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(req.PodName).
		Namespace(req.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   []string{req.Shell},
			Container: req.ContainerName,
			Stdin:     true,
			Stdout:    true,
			//Stderr:    true,
			TTY: true,
		}, schema.ParameterCodec)

	// 创建 SPDY 执行器（用于后续的命令执行）
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", execReq.URL())
	if err != nil {
		return nil, err
	}
	return executor, err
}

func (x *PodUseCase) DownloadPodLogs(ctx context.Context, req *DownloadPodLogsRequest) (string, error) {
	var (
		logData string
	)
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return logData, err
	}
	podLogOptions := corev1.PodLogOptions{
		Container: req.ContainerName,
		SinceTime: &metav1.Time{Time: req.SinceTime},
	}
	reqLogs := clientSet.CoreV1().Pods(req.Namespace).GetLogs(req.PodName, &podLogOptions)
	podLogs, err := reqLogs.Stream(ctx)
	if err != nil {
		return logData, fmt.Errorf("获取Pod日志失败: %w", err)
	}
	defer podLogs.Close()
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		logData += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return logData, fmt.Errorf("读取Pod日志失败: %v", err)
	}
	// 对日志内容进行 base64 编码
	encodedLog := base64.StdEncoding.EncodeToString([]byte(logData))
	return encodedLog, nil
}

//func (x *PodUseCase) ExecCommand(ctx context.Context, req *ExecCommandRequest) error {
//	// 获取集群的认证配置信息
//	cfg, err := x.Cluster.GetClientConfigByClusterName(ctx, req.ClusterName)
//	if err != nil {
//		return fmt.Errorf("获取集群认证信息失败: %v", err)
//	}
//
//	// 构建 Kubernetes 的 REST 配置信息
//	config, err := cfg.BuildRestConfig()
//	if err != nil {
//		return fmt.Errorf("构建k8s REST config失败: %v", err)
//	}
//
//	clientSet, err := x.GetClientSet(ctx, req.ClusterName)
//	if err != nil {
//		return fmt.Errorf("创建k8s客户端失败: %v", err)
//	}
//	// 初始化 Pod Exec 的请求，这里只进行连接，不执行任何命令
//	execReq := clientSet.CoreV1().RESTClient().Post().
//		Resource("pods").
//		Name(req.PodName).
//		Namespace(req.Namespace).
//		SubResource("exec").
//		VersionedParams(&corev1.PodExecOptions{
//			Command:   []string{req.Shell},
//			Container: req.ContainerName,
//			Stdin:     true,
//			Stdout:    true,
//			//Stderr:    true,
//			TTY:       true,
//		}, schema.ParameterCodec)
//
//	// 创建 SPDY 执行器（用于后续的命令执行）
//	executor, err := remotecommand.NewSPDYExecutor(config, "POST", execReq.URL())
//	if err != nil {
//		return fmt.Errorf("创建SPDYExecutor失败: %v", err)
//	}
//
//	streamOptions := remotecommand.StreamOptions{
//		Stdin:  req.Input,  // 输入流
//		Stdout: req.Output, // 输出流（返回给客户端）
//		//TerminalSizeQueue: req.TerminalSession, // 终端大小队列
//		Tty: true, // TTY 模式
//	}
//	return executor.StreamWithContext(ctx, streamOptions)
//}

func (x *PodUseCase) ListPodsByCloneSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	kruiseClientSet, err := x.Cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	cloneSet, err := kruiseClientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return nil, fmt.Errorf("没有权限访问CloneSet: %w", err)
		}
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("CloneSet不存在: %w", err)
		}
		return nil, fmt.Errorf("查询CloneSet失败: %w", err)
	}
	selector, err := metav1.LabelSelectorAsSelector(cloneSet.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("解析LabelSelector失败: %v", err)
	}
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	var (
		limit         = int64(req.PageSize)
		continueToken = ""
		result        []*corev1.Pod
	)
	for {
		podListOptions := metav1.ListOptions{
			LabelSelector: selector.String(),
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			return nil, fmt.Errorf("查询pods失败: %w", err)
		}
		for _, pod := range pods.Items {
			result = append(result, &pod)
		}
		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}
	return result, nil
}

func (x *PodUseCase) ListPodsByGameServerSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	gamekruiseClientSet, err := x.Cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	gameServerSet, err := gamekruiseClientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return nil, fmt.Errorf("没有权限访问GameServerSet: %w", err)
		}
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("GameServerSet不存在: %w", err)
		}
		return nil, fmt.Errorf("查询GameServerSet失败: %w", err)
	}
	Selector := labels.SelectorFromSet(map[string]string{
		gamekruiseiov1alpha1.GameServerOwnerGssKey: gameServerSet.Name,
	}).String()
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	var result []*corev1.Pod
	pod, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: Selector,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询pod失败: %v", err)
		return nil, fmt.Errorf("查询pod失败: %w", err)
	}
	for _, pod := range pod.Items {
		result = append(result, &pod)
	}
	return result, nil
}

func (x *PodUseCase) ListPodsByStatefulSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return nil, fmt.Errorf("没有权限访问StatefulSet: %v", err)
		}
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("StatefulSet不存在: %v", err)
		}
		x.log.WithContext(ctx).Errorf("查询Stateful失败: %v", err)
		return nil, fmt.Errorf("查询Stateful失败: %w", err)

	}
	selector, err := metav1.LabelSelectorAsSelector(statefulSet.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("解析LabelSelector失败: %v", err)
	}
	var (
		limit         = int64(req.PageSize)
		continueToken = ""
		result        []*corev1.Pod
	)
	for {
		podListOptions := metav1.ListOptions{
			LabelSelector: selector.String(),
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询pods失败: %v", err)
			return nil, fmt.Errorf("查询pods失败: %w", err)
		}
		for _, pod := range pods.Items {
			result = append(result, &pod)
		}
		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}
	return result, nil
}

func (x *PodUseCase) ListPodsByDaemonSet(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	daemonSet, err := clientSet.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		if errors.IsForbidden(err) {
			return nil, fmt.Errorf("没有权限访问DaemonSet: %v", err)
		}
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("DaemonSet不存在: %v", err)
		}
		x.log.WithContext(ctx).Errorf("查询DaemonSet失败: %v", err)
		return nil, fmt.Errorf("查询DaemonSet失败: %w", err)

	}

	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("解析LabelSelector失败: %v", err)
	}
	var (
		limit         = int64(req.PageSize)
		continueToken = ""
		result        []*corev1.Pod
	)
	for {
		podListOptions := metav1.ListOptions{
			LabelSelector: selector.String(),
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询pods失败: %v", err)
			return nil, fmt.Errorf("查询pods失败: %w", err)
		}
		for _, pod := range pods.Items {
			result = append(result, &pod)
		}
		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}
	return result, nil
}

func (x *PodUseCase) ListPodsByJob(ctx context.Context, req *ListControllerPodRequest) ([]*corev1.Pod, error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	job, err := clientSet.BatchV1().Jobs(req.Namespace).Get(ctx, req.ControllerName, metav1.GetOptions{})
	if err != nil {
		if errors.IsForbidden(err) {
			return nil, fmt.Errorf("没有权限访问Job: %v", err)
		}
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("job不存在: %v", err)
		}
		x.log.WithContext(ctx).Errorf("查询Job失败: %v", err)
		return nil, fmt.Errorf("查询Job失败: %w", err)
	}

	selector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("解析LabelSelector失败: %v", err)
	}
	var (
		limit         = int64(req.PageSize)
		continueToken = ""
		result        []*corev1.Pod
	)
	for {
		podListOptions := metav1.ListOptions{
			LabelSelector: selector.String(),
			Limit:         limit,
			Continue:      continueToken,
		}
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, podListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询pods失败: %v", err)
			return nil, fmt.Errorf("查询pods失败: %w", err)
		}
		for _, pod := range pods.Items {
			result = append(result, &pod)
		}
		if pods.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = pods.Continue
	}
	return result, nil
}
