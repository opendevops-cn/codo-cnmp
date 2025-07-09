package biz

import (
	ezrolloutv1 "codo-cnmp/common/ezrollout/v1"
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"regexp"
)

const (
	EzRolloutResource = "ezrollouts"
	EzRolloutKind     = "EzRollout"
)

var ezRolloutGVR = schema.GroupVersionResource{
	Group:    ezrolloutv1.GroupVersion.Group,
	Version:  ezrolloutv1.GroupVersion.Version,
	Resource: EzRolloutResource,
}

type EzRolloutCommonParams struct {
	// 集群名称
	ClusterName string `json:"clusterName"`
	// 命名空间
	Namespace string `json:"namespace"`
	// 名称
	Name string `json:"name"`
}

type ListEzRolloutRequest struct {
	EzRolloutCommonParams
	Keyword  string `json:"keyword"`
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
	ListAll  bool   `json:"listAll"`
}

type CreateEzRolloutRequest struct {
	EzRolloutCommonParams
	Selector         *metav1.LabelSelector        `json:"selector"`
	Labels           map[string]string            `json:"labels"`
	Annotations      map[string]string            `json:"annotations"`
	ScaleUpMetrics   []autoscalingv2.MetricSpec   // 扩容指标
	ScaleDownMetrics []autoscalingv2.MetricSpec   // 缩容指标
	EzOnlineScaler   *ezrolloutv1.EzOnlineScaler  // 线上版本伸缩
	EzOfflineScaler  *ezrolloutv1.EzOfflineScaler // 离线版本伸缩
	OnlineVersion    string                       // 线上版本名称
}

type DeleteEzRolloutRequest struct {
	EzRolloutCommonParams
}

type GetEzRolloutRequest struct {
	EzRolloutCommonParams
}

type CreateOrUpdateEzRolloutByYamlRequest struct {
	ClusterName string `json:"clusterName"`
	Yaml        string
}

type IEzRollout interface {
	// ListEzRollouts 列出EzRollout资源
	ListEzRollouts(ctx context.Context, req *ListEzRolloutRequest) ([]*ezrolloutv1.EzRollout, uint32, error)
	// CreateEzRollout 创建EzRollout资源
	CreateEzRollout(ctx context.Context, req *CreateEzRolloutRequest) error
	// UpdateEzRollout 更新EzRollout资源
	UpdateEzRollout(ctx context.Context, req *CreateEzRolloutRequest) error
	// DeleteEzRollout 删除EzRollout资源
	DeleteEzRollout(ctx context.Context, req *DeleteEzRolloutRequest) error
	// GetEzRollout 获取EzRollout资源详情
	GetEzRollout(ctx context.Context, req *GetEzRolloutRequest) (*ezrolloutv1.EzRollout, error)
	// CreateOrUpdateEzRolloutByYaml 创建或更新EzRollout资源
	CreateOrUpdateEzRolloutByYaml(ctx context.Context, req *CreateOrUpdateEzRolloutByYamlRequest) error
}

type EzRolloutUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *EzRolloutUseCase) filterObjectsByKeyword(objects *ezrolloutv1.EzRolloutList, keyword string) *ezrolloutv1.EzRolloutList {
	result := &ezrolloutv1.EzRolloutList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, object := range objects.Items {
		if utils.MatchString(pattern, object.Name) || utils.MatchLabels(pattern, object.Labels) || utils.MatchSelector(pattern, object.Spec.Selector) {
			result.Items = append(result.Items, object)
		}
	}
	return result
}

func (x *EzRolloutUseCase) ListEzRollouts(ctx context.Context, req *ListEzRolloutRequest) ([]*ezrolloutv1.EzRollout, uint32, error) {
	clientSet, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}
	var (
		allFilteredObjects = make([]*ezrolloutv1.EzRollout, 0)
		continueToken      = ""
		limit              = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		unstructObj, err := clientSet.Resource(ezRolloutGVR).Namespace(req.Namespace).List(ctx, ListOptions)
		if err != nil {
			if k8serrors.IsUnauthorized(err) {
				x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
				return allFilteredObjects, 0, err
			}
			if k8serrors.IsForbidden(err) {
				x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
				return allFilteredObjects, 0, err
			}
			if k8serrors.IsNotFound(err) {
				x.log.WithContext(ctx).Errorf("命名空间: %s ezrollout未找到", req.Namespace)
				return allFilteredObjects, 0, err
			}
			return allFilteredObjects, 0, fmt.Errorf("查询ezrollout失败: %w", err)
		}
		ezRolloutList := &ezrolloutv1.EzRolloutList{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), ezRolloutList)
		if err != nil {
			return allFilteredObjects, 0, fmt.Errorf("解析ezrollout失败: %w", err)
		}

		filteredGameServerSets := x.filterObjectsByKeyword(ezRolloutList, req.Keyword)
		for _, object := range filteredGameServerSets.Items {
			allFilteredObjects = append(allFilteredObjects, &object)
		}

		if ezRolloutList.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = ezRolloutList.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredObjects, uint32(len(allFilteredObjects)), nil
	}
	if len(allFilteredObjects) == 0 {
		return allFilteredObjects, 0, nil
	}
	// 否则分页返回结果
	paginatedObjects, total := utils.K8sPaginate(allFilteredObjects, req.Page, req.PageSize)
	return paginatedObjects, total, nil
}

func (x *EzRolloutUseCase) CreateEzRollout(ctx context.Context, req *CreateEzRolloutRequest) error {
	clientSet, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	ezRollout := &ezrolloutv1.EzRollout{
		TypeMeta: metav1.TypeMeta{
			Kind:       EzRolloutKind,
			APIVersion: ezrolloutv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: ezrolloutv1.EzRolloutSpec{
			Selector:         req.Selector,
			OnlineVersion:    req.OnlineVersion,
			ScaleUpMetrics:   req.ScaleUpMetrics,
			ScaleDownMetrics: req.ScaleDownMetrics,
			OnlineScaler:     req.EzOnlineScaler,
			OfflineScaler:    req.EzOfflineScaler,
		},
		Status: ezrolloutv1.EzRolloutStatus{},
	}
	unstructObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ezRollout)
	if err != nil {
		return fmt.Errorf("解析ezrollout失败: %w", err)
	}
	_, err = clientSet.Resource(ezRolloutGVR).Namespace(req.Namespace).Create(ctx, &unstructured.Unstructured{Object: unstructObj}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建ezrollout失败: %w", err)
	}
	return nil
}

func (x *EzRolloutUseCase) UpdateEzRollout(ctx context.Context, req *CreateEzRolloutRequest) error {
	clientSet, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("获取集群客户端失败: %w", err)
	}

	// 检查资源是否存在
	existing, err := clientSet.Resource(ezRolloutGVR).Namespace(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("要更新的资源不存在: %w", err)
		}
		return fmt.Errorf("获取现有资源失败: %w", err)
	}

	// 准备更新的资源
	ezRollout := &ezrolloutv1.EzRollout{
		TypeMeta: metav1.TypeMeta{
			Kind:       EzRolloutKind,
			APIVersion: ezrolloutv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            req.Name,
			Namespace:       req.Namespace,
			Labels:          req.Labels,
			Annotations:     req.Annotations,
			ResourceVersion: existing.GetResourceVersion(),
		},
		Spec: ezrolloutv1.EzRolloutSpec{
			Selector:         req.Selector,
			OnlineVersion:    req.OnlineVersion,
			ScaleUpMetrics:   req.ScaleUpMetrics,
			ScaleDownMetrics: req.ScaleDownMetrics,
			OnlineScaler:     req.EzOnlineScaler,
			OfflineScaler:    req.EzOfflineScaler,
		},
	}

	// 转换为非结构化对象
	unstructObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ezRollout)
	if err != nil {
		return fmt.Errorf("转换对象失败: %w", err)
	}

	// 执行更新操作
	_, err = clientSet.Resource(ezRolloutGVR).Namespace(req.Namespace).Update(
		ctx,
		&unstructured.Unstructured{Object: unstructObj},
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("更新资源失败: %w", err)
	}

	return nil
}

func (x *EzRolloutUseCase) DeleteEzRollout(ctx context.Context, req *DeleteEzRolloutRequest) error {
	clientSet, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.Resource(ezRolloutGVR).Namespace(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除ezrollout失败: %w", err)
	}
	return nil
}

func (x *EzRolloutUseCase) GetEzRollout(ctx context.Context, req *GetEzRolloutRequest) (*ezrolloutv1.EzRollout, error) {
	clientSet, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	unstructObj, err := clientSet.Resource(ezRolloutGVR).Namespace(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	ezRollout := &ezrolloutv1.EzRollout{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), ezRollout)
	if err != nil {
		return nil, err
	}
	return ezRollout, nil
}

func (x *EzRolloutUseCase) CreateOrUpdateEzRolloutByYaml(ctx context.Context, req *CreateOrUpdateEzRolloutByYamlRequest) error {
	clientSet, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	scheme := runtime.NewScheme()
	if err := autoscalingv2.AddToScheme(scheme); err != nil {
		return fmt.Errorf("注册 autoscaling/v2 失败: %w", err)
	}
	if err := ezrolloutv1.AddToScheme(scheme); err != nil {
		return fmt.Errorf("注册ezrollout资源失败: %w", err)
	}
	// 解析YAML为非结构化对象
	// obj := &unstructured.Unstructured{}
	// if err := yaml.Unmarshal([]byte(req.Yaml), &obj.Object); err != nil {
	// 	return fmt.Errorf("解析YAML失败: %w", err)
	// }
	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	obj, groupVersionKind, err := decode.Decode([]byte(req.Yaml), nil, nil)
	if err != nil {
		return fmt.Errorf("解析YAML失败: %w", err)
	}

	// 获取元数据
	ezrollout, ok := obj.(*ezrolloutv1.EzRollout)
	if !ok {
		return fmt.Errorf("无效的资源类型，期望 EzRollout，实际: %v", groupVersionKind)
	}

	// 转换为非结构化对象
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ezrollout)
	if err != nil {
		return fmt.Errorf("转换对象失败: %w", err)
	}

	namespace := ezrollout.Namespace
	name := ezrollout.Name

	// 尝试获取现有资源
	_, err = clientSet.Resource(ezRolloutGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 资源不存在，创建新资源
			_, err = clientSet.Resource(ezRolloutGVR).Namespace(namespace).Create(
				ctx,
				&unstructured.Unstructured{Object: unstructuredObj},
				metav1.CreateOptions{},
			)
			if err != nil {
				return fmt.Errorf("创建资源失败: %w", err)
			}
			return nil
		}
		return fmt.Errorf("检查资源是否存在失败: %w", err)
	}

	// 资源存在，更新资源
	_, err = clientSet.Resource(ezRolloutGVR).Namespace(namespace).Update(
		ctx,
		&unstructured.Unstructured{Object: unstructuredObj},
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("更新资源失败: %w", err)
	}

	return nil

}

func NewEzRolloutUseCase(cluster IClusterUseCase, logger log.Logger) *EzRolloutUseCase {
	return &EzRolloutUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/ezrollout")),
	}
}

func NewIEzRolloutUseCase(x *EzRolloutUseCase) IEzRollout {
	return x
}
