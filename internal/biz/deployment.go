package biz

import (
	"context"
	"fmt"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeploymentCommonParams 公共参数
type DeploymentCommonParams struct {
	ClusterName string // 集群名称
	Namespace   string // 命名空间
	Name        string // deployment 名称
}

type ListDeploymentRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

type CreateOrUpdateDeploymentByYamlRequest struct {
	ClusterName string // 集群名称
	Yaml        string // yaml 内容
}

type RollbackDeploymentRequest struct {
	DeploymentCommonParams
	Revision uint32 //   // 回滚版本
}

type ScaleDeploymentRequest struct {
	DeploymentCommonParams
	Replicas uint32 // 副本数量
}

type DeploymentDetailResponse struct {
	Deployment  *appsv1.Deployment   // deployment 详情
	ReplicaSets []*appsv1.ReplicaSet // deployment 关联的 replicaset 列表
	Pods        []*corev1.Pod        // deployment 关联的 pod 列表
	Events      []corev1.Event       // deployment 关联的事件列表
}

type UpdateDeploymentStrategyRequest struct {
	DeploymentCommonParams
	UpdateStrategyType uint32 // 更新策略类型
	MaxSurge           string // 最大超载
	MaxUnavailable     string // 最大不可用
}

type IDeploymentUseCase interface {
	// ListDeployment 获取 deployment 列表
	ListDeployment(ctx context.Context, req *ListDeploymentRequest) ([]*appsv1.Deployment, uint32, error)
	// DeleteDeployment 删除 deployment
	DeleteDeployment(ctx context.Context, req *DeploymentCommonParams) error
	// CreateOrUpdateDeploymentByYaml 创建或更新 deployment
	CreateOrUpdateDeploymentByYaml(ctx context.Context, req *CreateOrUpdateDeploymentByYamlRequest) error
	// RestartDeployment 重启 deployment
	RestartDeployment(ctx context.Context, req *DeploymentCommonParams) error
	// RollbackDeployment 回滚 deployment
	RollbackDeployment(ctx context.Context, req *RollbackDeploymentRequest) error
	// ScaleDeployment 伸缩 deployment
	ScaleDeployment(ctx context.Context, req *ScaleDeploymentRequest) error
	// GetDeploymentYaml 获取 deployment yaml
	GetDeploymentYaml(ctx context.Context, req *DeploymentCommonParams) (string, error)
	// GetDeploymentHistory 获取 deployment 历史版本
	GetDeploymentHistory(ctx context.Context, req *DeploymentCommonParams) ([]*appsv1.ReplicaSet, error)
	// GetDeploymentDetail 获取 deployment 详情
	GetDeploymentDetail(ctx context.Context, req *DeploymentCommonParams) (DeploymentDetailResponse, error)
	// GetDeployment 获取 deployment
	GetDeployment(ctx context.Context, req *DeploymentCommonParams) (*appsv1.Deployment, error)
	// UpdateDeploymentStrategy 更新 deployment更新策略
	UpdateDeploymentStrategy(ctx context.Context, req *UpdateDeploymentStrategyRequest) (bool, error)
}

func NewIDeploymentUseCase(x *DeploymentUseCase) IDeploymentUseCase {
	return x
}

type DeploymentUseCase struct {
	cluster IClusterUseCase
	pod     IPodUseCase
	log     *log.Helper
}

func parseIntOrPercent(value string) (*intstr.IntOrString, error) {
	// 检查是否为百分比
	if strings.HasSuffix(value, "%") {
		// 是百分比，直接返回字符串形式
		return &intstr.IntOrString{
			Type:   intstr.String,
			StrVal: value,
		}, nil
	}

	// 尝试解析为整数
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("解析整数失败: %v", err)
	}

	// 是整数，返回整数形式
	return &intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: int32(intValue),
	}, nil
}

func (x *DeploymentUseCase) UpdateDeploymentStrategy(ctx context.Context, req *UpdateDeploymentStrategyRequest) (bool, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, err
	}
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return false, fmt.Errorf("查询 deployment 失败: %w", err)
	}
	var (
		maxSurge       *intstr.IntOrString
		maxUnavailable *intstr.IntOrString
	)
	if req.MaxSurge != "" {
		maxSurge, err = parseIntOrPercent(req.MaxSurge)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 deployment更新策略 maxSurge 失败: %v", err)
			return false, fmt.Errorf("解析 maxSurge 失败: %v", err)
		}
	}
	if req.MaxUnavailable != "" {
		maxUnavailable, err = parseIntOrPercent(req.MaxUnavailable)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 deployment更新策略 maxUnavailable 失败: %v", err)
			return false, fmt.Errorf("解析 maxUnavailable 失败: %v", err)
		}
	}

	// 更新策略
	switch req.UpdateStrategyType {
	case uint32(pb.UpdateDeploymentStrategyRequest_RollingUpdate):
		deployment.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxSurge:       maxSurge,
				MaxUnavailable: maxUnavailable,
			},
		}
	case uint32(pb.UpdateDeploymentStrategyRequest_Recreate):
		deployment.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
		}
	default:
		return false, fmt.Errorf("未知的策略类型: %d", req.UpdateStrategyType)
	}

	_, err = clientSet.AppsV1().Deployments(req.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新 deployment 策略失败: %v", err)
		return false, fmt.Errorf("更新 deployment 策略失败: %w", err)
	}
	return true, nil
}

func NewDeploymentUseCase(cluster IClusterUseCase, pod IPodUseCase, logger log.Logger) *DeploymentUseCase {
	return &DeploymentUseCase{
		cluster: cluster,
		pod:     pod,
		log:     log.NewHelper(log.With(logger, "module", "biz/deployment")),
	}
}

// ListDeployment 获取 deployment 列表
func (x *DeploymentUseCase) ListDeployment(ctx context.Context, req *ListDeploymentRequest) ([]*appsv1.Deployment, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
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
		allFilteredDeployments = make([]*appsv1.Deployment, 0)
		limit                  = int64(req.PageSize)
		continueToken          = ""
	)

	for {
		deploymentListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		deployments, err := clientSet.AppsV1().Deployments(req.Namespace).List(ctx, deploymentListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询deployments失败: %v", err)
			return nil, 0, fmt.Errorf("查询deployments失败: %w", err)
		}
		filteredDeployments := filterDeploymentsByKeyword(deployments, req.Keyword)
		for _, deployment := range filteredDeployments.Items {
			allFilteredDeployments = append(allFilteredDeployments, &deployment)
		}

		if deployments.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = deployments.Continue
	}
	if len(allFilteredDeployments) > 0 {
		gvks, _, err := scheme.Scheme.ObjectKinds(allFilteredDeployments[0])
		if err != nil {
			return nil, 0, fmt.Errorf("查询deployment失败: %w", err)
		}
		if len(gvks) > 0 {
			for _, deployment := range allFilteredDeployments {
				deployment.Kind = gvks[0].Kind
				deployment.APIVersion = gvks[0].GroupVersion().String()
			}
		}
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredDeployments, uint32(len(allFilteredDeployments)), nil
	}
	if len(allFilteredDeployments) == 0 {
		return allFilteredDeployments, 0, nil
	}
	// 否则分页返回结果
	paginatedDeployments, total := utils.K8sPaginate(allFilteredDeployments, req.Page, req.PageSize)
	return paginatedDeployments, total, nil

}

// 检查关键词是否在 Pod 名称、状态、命名空间、镜像中
func filterDeploymentsByKeyword(deployments *appsv1.DeploymentList, keyword string) *appsv1.DeploymentList {
	result := &appsv1.DeploymentList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, deployment := range deployments.Items {
		if utils.MatchString(pattern, deployment.Name) ||
			utils.MatchLabels(pattern, deployment.Labels) ||
			utils.MatchContainerImages(pattern, deployment.Spec.Template.Spec.Containers) {
			result.Items = append(result.Items, deployment)
		}
	}
	return result
}

// CreateOrUpdateDeploymentByYaml 创建或更新 deployment
// 注意：yaml 内容必须包含 namespace 字段，否则会创建失败
func (x *DeploymentUseCase) CreateOrUpdateDeploymentByYaml(ctx context.Context, req *CreateOrUpdateDeploymentByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	deployment := &appsv1.Deployment{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&deployment); err != nil {
		x.log.WithContext(ctx).Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "解析YAML失败, 请检查格式是否正确")
	}
	// 获取当前的deployment对象
	_, err = clientSet.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 创建deployment
			_, err = clientSet.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("创建deployment失败: %w", err)
			}
			return nil
		}
		if k8serrors.IsForbidden(err) {
			return fmt.Errorf("查询deployment失败,: %w", err)
		}
		return fmt.Errorf("查询deployment失败: %w", err)
	}

	// 更新deployment
	_, err = clientSet.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新deployment失败: %v", err)
		return fmt.Errorf("更新deployment失败: %w", err)
	}
	return nil

}

// DeleteDeployment 删除 deployment
func (x *DeploymentUseCase) DeleteDeployment(ctx context.Context, req *DeploymentCommonParams) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	// todo 默认删除策略是立即删除，后续可以添加删除策略
	err = clientSet.AppsV1().Deployments(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除deployment失败: %v", err)
		return fmt.Errorf("删除deployment失败: %w", err)
	}
	return nil
}

// RestartDeployment 重启 deployment
func (x *DeploymentUseCase) RestartDeployment(ctx context.Context, req *DeploymentCommonParams) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(context.Background(), req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return fmt.Errorf("查询 deployment 失败: %w", err)
	}
	// 更新 Deployment 的 spec.template.metadata.annotations 触发restart
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// 更新 Deployment
	_, err = clientSet.AppsV1().Deployments(req.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("重启 deployment 失败: %v", err)
		return fmt.Errorf("重启 deployment 失败: %w", err)
	}
	return nil
}

// RollbackDeployment 回滚 deployment
func (x *DeploymentUseCase) RollbackDeployment(ctx context.Context, req *RollbackDeploymentRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(context.Background(), req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return fmt.Errorf("查询 deployment 失败: %w", err)
	}
	// 获取与 Deployment 关联的 ReplicaSets
	rsList, err := clientSet.AppsV1().ReplicaSets(req.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Deployment关联的ReplicaSets失败: %v", err)
		return fmt.Errorf("查询 ReplicaSets 失败: %w", err)
	}

	var replicaSet *appsv1.ReplicaSet
	// 找到指定版本的 ReplicaSet
	for _, rs := range rsList.Items {
		if rs.Annotations != nil && rs.Annotations["deployment.kubernetes.io/revision"] == strconv.Itoa(int(req.Revision)) {
			replicaSet = &rs
			break
		}
	}
	if replicaSet == nil {
		return fmt.Errorf("未找到指定版本的 ReplicaSet")
	}
	// 更新 Deployment 的 spec.template 触发回滚
	deployment.Spec.Template = replicaSet.Spec.Template
	if _, err = clientSet.AppsV1().Deployments(req.Namespace).Update(ctx, deployment, metav1.UpdateOptions{}); err != nil {
		x.log.WithContext(ctx).Errorf("回滚 deployment 失败: %v", err)
		return fmt.Errorf("回滚 deployment 失败: %w", err)
	}
	return nil

}

// ScaleDeployment 缩放 deployment
func (x *DeploymentUseCase) ScaleDeployment(ctx context.Context, req *ScaleDeploymentRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(context.Background(), req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return fmt.Errorf("查询 deployment 失败: %w", err)
	}
	// 更新 Deployment 的 spec.replicas 触发缩放
	replicasCount := int32(req.Replicas)
	deployment.Spec.Replicas = &replicasCount
	// 更新 Deployment
	_, err = clientSet.AppsV1().Deployments(req.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("缩放 deployment 失败: %v", err)
		return fmt.Errorf("缩放 deployment 失败: %w", err)
	}
	return nil
}

func (x *DeploymentUseCase) GetDeploymentYaml(ctx context.Context, req *DeploymentCommonParams) (string, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return "", err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("查询 deployment 失败: %w", err)
	}
	if deployment.APIVersion == "" {
		deployment.APIVersion = "apps/v1"
	}
	if deployment.Kind == "" {
		deployment.Kind = "Deployment"
	}
	return ConvertResourceTOYaml(deployment)
}

func ConvertResourceTOYaml(obj runtime.Object) (string, error) {
	encoder := json.NewSerializerWithOptions(
		json.DefaultMetaFactory, nil, nil,
		json.SerializerOptions{Yaml: true, Pretty: true, Strict: true},
	)
	var jsonBuffer strings.Builder
	if err := encoder.Encode(obj, &jsonBuffer); err != nil {
		return "", fmt.Errorf("序列化对象为 YAML 失败: %v", err)
	}
	return jsonBuffer.String(), nil
}

// GetDeploymentHistory 获取 deployment 历史版本
func (x *DeploymentUseCase) GetDeploymentHistory(ctx context.Context, req *DeploymentCommonParams) ([]*appsv1.ReplicaSet, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return nil, fmt.Errorf("查询 deployment 失败: %w", err)
	}
	// 获取与 Deployment 关联的 ReplicaSets
	rsList, err := clientSet.AppsV1().ReplicaSets(req.Namespace).List(ctx, metav1.ListOptions{
		//LabelSelector: "deployment=" + deployment.Name,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询Deployment关联的ReplicaSets失败: %v", err)
		return nil, fmt.Errorf("查询 ReplicaSets 失败: %w", err)
	}
	var history []*appsv1.ReplicaSet
	for _, rs := range rsList.Items {
		if ownerRef := metav1.GetControllerOf(&rs); ownerRef != nil && ownerRef.UID == deployment.UID {
			history = append(history, &rs)
		}
	}
	// 按照创建时间排序(倒序)
	sort.Slice(history, func(i, j int) bool {
		return history[i].CreationTimestamp.Time.After(history[j].CreationTimestamp.Time)
	})
	return history, nil
}

// GetDeploymentDetail 获取 deployment 详情
func (x *DeploymentUseCase) GetDeploymentDetail(ctx context.Context, req *DeploymentCommonParams) (DeploymentDetailResponse, error) {
	var deploymentDetail DeploymentDetailResponse
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return deploymentDetail, err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return deploymentDetail, fmt.Errorf("查询 deployment 失败: %w", err)
	}
	// 获取与 Deployment 关联的 ReplicaSets
	//rsList, err := uc.GetDeploymentHistory(ctx, req)
	//if err != nil {
	//	return deploymentDetail, fmt.Errorf("查询 ReplicaSets 失败: %v", err)
	//}
	// 获取与 Deployment 关联的 Pods
	//podList, err := uc.pod.LzistPodsByDeployment(ctx, &QueryPodItem{
	//	ClusterName:    req.ClusterName,
	//	Namespace:      req.Namespace,
	//	DeploymentName: deployment.Name,
	//	ListAll:        true,
	//	PageSize:       10,
	//})
	//if err != nil {
	//	return deploymentDetail, fmt.Errorf("查询 Pods 失败: %v", err)
	//}
	return DeploymentDetailResponse{
		Deployment: deployment,
		//ReplicaSets: rsList,
		//Pods: podList,
	}, nil
}

// GetDeployment 获取 deployment
func (x *DeploymentUseCase) GetDeployment(ctx context.Context, req *DeploymentCommonParams) (*appsv1.Deployment, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	// 获取现有的 Deployment
	deployment, err := clientSet.AppsV1().Deployments(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询 deployment 失败: %v", err)
		return nil, fmt.Errorf("查询 deployment 失败: %w", err)
	}
	return deployment, nil
}
