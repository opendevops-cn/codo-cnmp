package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"regexp"
)

type ResourceQuotaCommonParams struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

type ListResourceQuotaRequest struct {
	ResourceQuotaCommonParams
	Page     uint32
	PageSize uint32
	ListAll  bool
	Keyword  string
}

type GetResourceQuotaRequest struct {
	ResourceQuotaCommonParams
}

// CreateResourceQuotaRequest is the request struct for creating a limit range
type CreateResourceQuotaRequest struct {
	ResourceQuotaCommonParams
	Hard  corev1.ResourceList         `json:"hard"`
	Scope []corev1.ResourceQuotaScope `json:"scope"`
}

type UpdateResourceQuotaRequest struct {
	ResourceQuotaCommonParams
	Hard  corev1.ResourceList       `json:"hard"`
	Scope corev1.ResourceQuotaScope `json:"scope"`
}

type DeleteResourceQuotaRequest struct {
	ResourceQuotaCommonParams
}

type IResourceQuotaUseCase interface {
	GetResourceQuota(ctx context.Context, req *GetResourceQuotaRequest) (*corev1.ResourceQuota, error)
	ListResourceQuota(ctx context.Context, req *ListResourceQuotaRequest) ([]*corev1.ResourceQuota, uint32, error)
	CreateResourceQuota(ctx context.Context, req *CreateResourceQuotaRequest) error
	UpdateResourceQuota(ctx context.Context, req *UpdateResourceQuotaRequest) error
	DeleteResourceQuota(ctx context.Context, req *DeleteResourceQuotaRequest) error
	CreateOrUpdateResourceQuota(ctx context.Context, req *CreateResourceQuotaRequest) error
}

type ResourceQuotaUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *ResourceQuotaUseCase) CreateOrUpdateResourceQuota(ctx context.Context, req *CreateResourceQuotaRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	_, err = clientSet.CoreV1().ResourceQuotas(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 资源配额不存在，创建
			return x.CreateResourceQuota(ctx, req)
		}
		return fmt.Errorf("查询resourceQuota失败 %s.%s: %w", req.Namespace, req.Name, err)
	}
	// 存在则更新
	return x.UpdateResourceQuota(ctx, &UpdateResourceQuotaRequest{
		ResourceQuotaCommonParams: req.ResourceQuotaCommonParams,
		Hard:                      req.Hard,
	})
}

func (x *ResourceQuotaUseCase) DeleteResourceQuota(ctx context.Context, req *DeleteResourceQuotaRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().ResourceQuotas(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除resourceQuota失败 %s.%s: %w", req.Namespace, req.Name, err)
	}
	return nil
}

func (x *ResourceQuotaUseCase) GetResourceQuota(ctx context.Context, req *GetResourceQuotaRequest) (*corev1.ResourceQuota, error) {
	cluster, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	resourceQuota, err := cluster.CoreV1().ResourceQuotas(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询limitRange失败 %s: %w", req.Namespace, err)
	}
	return resourceQuota, nil
}

func (x *ResourceQuotaUseCase) ListResourceQuota(ctx context.Context, req *ListResourceQuotaRequest) ([]*corev1.ResourceQuota, uint32, error) {
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
		allFilteredResourceQuota = make([]*corev1.ResourceQuota, 0)
		continueToken            = ""
		limit                    = int64(req.PageSize)
	)
	for {
		resourceQuotaListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		resourceQuotaList, err := clientSet.CoreV1().ResourceQuotas(req.Namespace).List(ctx, resourceQuotaListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取LimitRange列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取LimitRange列表失败: %w", err)
		}
		filteredConfigMapList := x.filterResourceQuotaByKeyword(resourceQuotaList, req.Keyword)
		for _, configMap := range filteredConfigMapList.Items {
			allFilteredResourceQuota = append(allFilteredResourceQuota, &configMap)
		}
		if resourceQuotaList.Continue == "" {
			break
		}
		continueToken = resourceQuotaList.Continue
	}
	if len(allFilteredResourceQuota) > 0 {
		gvks, _, err := scheme.Scheme.ObjectKinds(allFilteredResourceQuota[0])
		if err != nil {
			return nil, 0, fmt.Errorf("查询limitRange失败: %w", err)
		}
		if len(gvks) > 0 {
			for _, deployment := range allFilteredResourceQuota {
				deployment.Kind = gvks[0].Kind
				deployment.APIVersion = gvks[0].GroupVersion().String()
			}
		}
	}
	if req.ListAll {
		return allFilteredResourceQuota, uint32(len(allFilteredResourceQuota)), nil
	}
	if len(allFilteredResourceQuota) == 0 {
		return allFilteredResourceQuota, 0, nil
	}
	// 否则分页返回结果
	paginatedConfigMaps, total := utils.K8sPaginate(allFilteredResourceQuota, req.Page, req.PageSize)
	return paginatedConfigMaps, total, nil

}

func (x *ResourceQuotaUseCase) filterResourceQuotaByKeyword(resourceQuotaList *corev1.ResourceQuotaList, keyword string) *corev1.ResourceQuotaList {
	if keyword == "" {
		return resourceQuotaList
	}
	result := &corev1.ResourceQuotaList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, resourceQuota := range resourceQuotaList.Items {
		if utils.MatchString(pattern, resourceQuota.Name) {
			result.Items = append(result.Items, resourceQuota)
		}
	}
	return result
}

func (x *ResourceQuotaUseCase) CreateResourceQuota(ctx context.Context, req *CreateResourceQuotaRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	resourceQuota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: req.Hard,
		},
	}
	// 根据n是否已经存在
	if _, err := clientSet.CoreV1().ResourceQuotas(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{}); err == nil {
		x.log.WithContext(ctx).Errorf("同名ResourceQuota已存在: %v", err)
		return fmt.Errorf("创建resourceQuota失败: %w", err)
	}
	// Dry-run 校验
	_, err = clientSet.CoreV1().ResourceQuotas(req.Namespace).Create(ctx, resourceQuota, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("dry-run校验失败: %w", err)
	}
	// 非 dry-run 模式下创建 limitRange
	_, err = clientSet.CoreV1().ResourceQuotas(req.Namespace).Create(ctx, resourceQuota, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建resourceQuota失败 %s: %w", req.Namespace, err)
	}
	return nil
}

func (x *ResourceQuotaUseCase) UpdateResourceQuota(ctx context.Context, req *UpdateResourceQuotaRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	resourceQuota, err := clientSet.CoreV1().ResourceQuotas(req.Namespace).Get(context.Background(), req.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("查询limitRange失败 %s: %w", req.Namespace, err)
	}
	resourceQuota.Spec.Hard = req.Hard
	// dry-run 校验
	_, err = clientSet.CoreV1().ResourceQuotas(req.Namespace).Update(ctx, resourceQuota, metav1.UpdateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("dry-run校验失败: %w", err)
	}
	_, err = clientSet.CoreV1().ResourceQuotas(req.Namespace).Update(ctx, resourceQuota, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新limitRange失败 %s: %w", req.Namespace, err)
	}
	return nil
}

func NewResourceQuotaUseCase(cluster IClusterUseCase, logger log.Logger) *ResourceQuotaUseCase {
	return &ResourceQuotaUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/resourceQuota")),
	}
}

func NewIResourceQuotaUseCase(x *ResourceQuotaUseCase) IResourceQuotaUseCase {
	return x
}
