package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"codo-cnmp/common/utils"
	"codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type SvcCommonParams struct {
	ClusterName string
	Namespace   string
	Name        string
}

type ListSvcRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	ListAll     bool
	Page        uint32
	PageSize    uint32
	SvcType     pb.SvcType
}

type GetSvcRequest struct {
	SvcCommonParams
}

type DeleteSvcRequest struct {
	SvcCommonParams
}

type CreateSvcRequest struct {
	SvcCommonParams
	Labels                   map[string]string
	Annotations              map[string]string
	Selector                 map[string]string
	PublishNotReadyAddresses bool
	ExternalIPs              []string
	SessionAffinity          string
	SessionAffinitySeconds   uint32
	Headless                 bool
	Ports                    []corev1.ServicePort
	SvcType                  pb.SvcType
	ExternalName             string
}

type ISvcUseCase interface {
	ListSvc(ctx context.Context, req *ListSvcRequest) ([]*corev1.Service, uint32, error)
	GetSvc(ctx context.Context, req *GetSvcRequest) (*corev1.Service, error)
	DeleteSvc(ctx context.Context, req *DeleteSvcRequest) error
	UpdateSvc(ctx context.Context, req *CreateSvcRequest) error
	CreateSvc(ctx context.Context, req *CreateSvcRequest) error
}

type SvcUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *SvcUseCase) UpdateSvc(ctx context.Context, req *CreateSvcRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("获取集群客户端失败: %w", err)
	}

	// 检查 Service 是否存在
	svc, err := clientSet.CoreV1().Services(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取 Service 失败: %v", err)
		return fmt.Errorf("获取 Service 失败: %w", err)
	}

	// 构建补丁数据
	patchOperations := x.buildServiceJsonPatch(req, svc)

	// 转换为 JSON
	patchBytes, err := json.Marshal(patchOperations)
	if err != nil {
		return fmt.Errorf("序列化补丁数据失败: %w", err)
	}

	// 使用 StrategicMergePatch 更新 Service
	_, err = clientSet.CoreV1().Services(req.Namespace).Patch(
		ctx,
		req.Name,
		types.JSONPatchType,
		patchBytes,
		metav1.PatchOptions{},
	)
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新Svc失败: %v", err)
		return fmt.Errorf("更新 Service 失败: %w", err)
	}

	return nil
}

func (x *SvcUseCase) buildServiceJsonPatch(req *CreateSvcRequest, svc *corev1.Service) []map[string]interface{} {
	var operations []map[string]interface{}

	// 处理labels
	if req.Labels != nil {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/metadata/labels",
			"value": req.Labels,
		})
	}

	// 处理annotations
	if req.Annotations != nil {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/metadata/annotations",
			"value": req.Annotations,
		})
	}

	// 处理selector
	if req.Selector != nil {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/selector",
			"value": req.Selector,
		})
	}

	// 处理ports
	if req.Ports != nil {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/ports",
			"value": req.Ports,
		})
	}

	// 处理服务类型
	if corev1.ServiceType(req.SvcType) != "" {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/type",
			"value": corev1.ServiceType(req.SvcType.String()),
		})
	}

	// 处理publishNotReadyAddresses
	operations = append(operations, map[string]interface{}{
		"op":    "replace",
		"path":  "/spec/publishNotReadyAddresses",
		"value": req.PublishNotReadyAddresses,
	})

	// 处理externalIPs
	operations = append(operations, map[string]interface{}{
		"op":    "replace",
		"path":  "/spec/externalIPs",
		"value": req.ExternalIPs,
	})

	// 处理sessionAffinity
	if req.SessionAffinity != "" {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/sessionAffinity",
			"value": req.SessionAffinity,
		})

		if req.SessionAffinitySeconds > 0 && req.SessionAffinity == "ClientIP" {
			operations = append(operations, map[string]interface{}{
				"op":   "replace",
				"path": "/spec/sessionAffinityConfig",
				"value": map[string]interface{}{
					"clientIP": map[string]interface{}{
						"timeoutSeconds": req.SessionAffinitySeconds,
					},
				},
			})
		}
	}

	// 处理Headless
	if req.Headless {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/clusterIP",
			"value": "None",
		})
	}

	// 处理ExternalName
	if req.ExternalName != "" {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/externalName",
			"value": req.ExternalName,
		})
	}

	return operations
}

func (x *SvcUseCase) CreateSvc(ctx context.Context, req *CreateSvcRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("创建k8s client失败: %w", err)
	}
	var timeoutSeconds *int32
	if req.SessionAffinitySeconds > 0 {
		timeoutSeconds = new(int32)
		*timeoutSeconds = int32(req.SessionAffinitySeconds)
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports:                    req.Ports,
			Selector:                 req.Selector,
			Type:                     corev1.ServiceType(req.SvcType.String()),
			ExternalIPs:              req.ExternalIPs,
			PublishNotReadyAddresses: req.PublishNotReadyAddresses,
			SessionAffinityConfig:    nil,
		},
	}
	if req.SessionAffinity != "" {
		svc.Spec.SessionAffinity = corev1.ServiceAffinity(req.SessionAffinity)
	}
	if req.SessionAffinity == "ClientIP" {
		svc.Spec.SessionAffinityConfig = &corev1.SessionAffinityConfig{
			ClientIP: &corev1.ClientIPConfig{
				TimeoutSeconds: timeoutSeconds,
			},
		}
	}
	if req.Headless {
		svc.Spec.ClusterIP = "None"
	}
	if req.ExternalName != "" {
		svc.Spec.ExternalName = req.ExternalName
	}
	_, err = clientSet.CoreV1().Services(req.Namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建Svc失败: %v", err)
		return fmt.Errorf("创建Svc失败: %w", err)
	}
	return nil
}

func NewSvcUseCase(cluster IClusterUseCase, logger log.Logger) *SvcUseCase {
	return &SvcUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/svc")),
	}
}

func (x *SvcUseCase) ListSvc(ctx context.Context, req *ListSvcRequest) ([]*corev1.Service, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, fmt.Errorf("创建k8s client失败: %w", err)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var svcType string
	switch req.SvcType {
	case pb.SvcType_SVC_TYPE_UNSPECIFIED:
		svcType = ""
	default:
		svcType = req.SvcType.String()
	}
	var (
		allFilteredSvcList = make([]*corev1.Service, 0)
		continueToken      = ""
		limit              = int64(req.PageSize)
	)
	for {
		svcListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		svcList, err := clientSet.CoreV1().Services(req.Namespace).List(ctx, svcListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取Svc列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取Svc列表失败: %w", err)
		}
		filteredSvcList := x.filterSvcByKeyword(svcList, req.Keyword, svcType)
		for _, svc := range filteredSvcList.Items {
			allFilteredSvcList = append(allFilteredSvcList, &svc)
		}
		if svcList.Continue == "" {
			break
		}
		continueToken = svcList.Continue
	}
	if req.ListAll {
		return allFilteredSvcList, uint32(len(allFilteredSvcList)), nil
	}
	if len(allFilteredSvcList) == 0 {
		return allFilteredSvcList, 0, nil
	}
	// 否则分页返回结果
	paginatedSvcs, total := utils.K8sPaginate(allFilteredSvcList, req.Page, req.PageSize)
	return paginatedSvcs, total, nil
}

func (x *SvcUseCase) GetSvc(ctx context.Context, req *GetSvcRequest) (*corev1.Service, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	svc, err := clientSet.CoreV1().Services(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取Svc失败: %v", err)
		return nil, fmt.Errorf("获取Svc失败: %w", err)
	}
	return svc, nil
}

func (x *SvcUseCase) DeleteSvc(ctx context.Context, req *DeleteSvcRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().Services(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除Svc失败: %v", err)
		return fmt.Errorf("删除Svc失败: %w", err)
	}
	return nil
}

// matchServiceFields 检查服务的各个字段是否匹配搜索条件
func matchServiceFields(svc *corev1.Service, pattern string) bool {
	// 1. 检查服务名称
	if utils.MatchString(pattern, svc.Name) {
		return true
	}
	// 3. 检查 ClusterIP
	if utils.MatchString(pattern, svc.Spec.ClusterIP) {
		return true
	}

	// 4. 检查端口信息
	for _, port := range svc.Spec.Ports {
		// 检查端口名称
		if port.Name != "" && utils.MatchString(pattern, port.Name) {
			return true
		}

		// 检查端口号（转换为字符串进行匹配）
		portStr := fmt.Sprintf("%d", port.Port)
		if utils.MatchString(pattern, portStr) {
			return true
		}

		// 检查目标端口号
		targetPortStr := port.TargetPort.String()
		if targetPortStr != "" && utils.MatchString(pattern, targetPortStr) {
			return true
		}

		// 检查协议
		if string(port.Protocol) != "" && utils.MatchString(pattern, string(port.Protocol)) {
			return true
		}
	}

	//// 5. 检查外部IP（如果有）
	//for _, externalIP := range svc.Spec.ExternalIPs {
	//	if utils.MatchString(pattern, externalIP) {
	//		return true
	//	}
	//}
	//// 6. 检查 LoadBalancer Ingress IP/Hostname
	//for _, ingress := range svc.Status.LoadBalancer.Ingress {
	//	if ingress.IP != "" && utils.MatchString(pattern, ingress.IP) {
	//		return true
	//	}
	//	if ingress.Hostname != "" && utils.MatchString(pattern, ingress.Hostname) {
	//		return true
	//	}
	//}
	//
	//// 7. 检查 ExternalName（如果是 ExternalName 类型）
	//if svc.Spec.Type == corev1.ServiceTypeExternalName && utils.MatchString(pattern, svc.Spec.ExternalName) {
	//	return true
	//}

	return false
}

func (x *SvcUseCase) filterSvcByKeyword(list *corev1.ServiceList, keyword string, svcType string) *corev1.ServiceList {
	if keyword == "" && svcType == "" {
		return list
	}
	result := &corev1.ServiceList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, svc := range list.Items {
		if svcType == "" {
			if matchServiceFields(&svc, pattern) {
				result.Items = append(result.Items, svc)
			}
		} else {
			if svc.Spec.Type == corev1.ServiceType(svcType) && matchServiceFields(&svc, pattern) {
				result.Items = append(result.Items, svc)
			}
		}
	}
	return result
}

func (x *SvcUseCase) GetSvcReferences(ctx context.Context, svc *corev1.Service, req *SvcCommonParams) (int32, []map[string]string, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return 0, nil, err
	}

	labelSelector := labels.SelectorFromSet(svc.Spec.Selector).String()
	if labelSelector == "" {
		return 0, nil, nil
	}

	refMap := make(map[string]struct{})
	podCount := 0
	continueToken := ""

	for {
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{
			Limit:         100,
			Continue:      continueToken,
			LabelSelector: labelSelector,
		})
		if err != nil {
			return 0, nil, err
		}

		// 处理当前页的 pods
		for _, pod := range pods.Items {
			topOwner, err := x.findTopOwner(ctx, clientSet, pod.Namespace, pod.OwnerReferences)
			if err != nil {
				return 0, nil, err
			}
			podCount++
			if topOwner != nil {
				refMap[fmt.Sprintf("%s/%s", topOwner.Kind, topOwner.Name)] = struct{}{}
			} else {
				refMap[fmt.Sprintf("Pod/%s", pod.Name)] = struct{}{}
			}
		}
		continueToken = pods.Continue
		if continueToken == "" {
			break
		}
	}

	references := make([]string, 0, len(refMap))
	for ref := range refMap {
		references = append(references, ref)
	}

	// 返回结果
	result := make([]map[string]string, len(references))
	for i, ref := range references {
		kind, name, _ := strings.Cut(ref, "/")
		result[i] = map[string]string{
			"kind": kind,
			"name": name,
		}
	}
	return int32(podCount), result, nil
}

func (x *SvcUseCase) findTopOwner(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, owners []metav1.OwnerReference) (*metav1.OwnerReference, error) {
	if len(owners) == 0 {
		return nil, nil
	}
	owner := owners[0]
	switch owner.Kind {
	case "ReplicaSet":
		rs, err := clientSet.AppsV1().ReplicaSets(namespace).Get(ctx, owner.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if len(rs.OwnerReferences) > 0 {
			return x.findTopOwner(ctx, clientSet, namespace, rs.OwnerReferences)
		}
	case "StatefulSet", "DaemonSet", "Deployment", "Job", "CronJob", "CloneSet":
		return &owner, nil
	}

	return &owner, nil
}

func NewISvcUseCase(x *SvcUseCase) ISvcUseCase {
	return x
}
