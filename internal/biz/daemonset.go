package biz

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"regexp"
	"strings"
	"time"
)

type DaemonSetCommonParams struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

// ListDaemonSetRequest 创建或更新 DaemonSet
type ListDaemonSetRequest struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Keyword     string `json:"keyword"`
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"page_size"`
	ListAll     bool   `json:"list_all"`
}

type CreateOrUpdateDaemonSetByYamlRequest struct {
	ClusterName string `json:"cluster_name"`
	Yaml        string `json:"yaml"`
}

type RollbackDaemonSetRequest struct {
	DaemonSetCommonParams
	Revision uint32 `json:"revision"`
}

type UpdateDaemonSetStrategyRequest struct {
	DaemonSetCommonParams
	UpdateStrategyType uint32 // 更新策略类型
	MaxSurge           string // 最大超载
	MaxUnavailable     string // 最大不可用
}

type IDaemonSetUseCase interface {
	// ListDaemonSet DaemonSet列表.
	ListDaemonSet(ctx context.Context, req *ListDaemonSetRequest) ([]*appsv1.DaemonSet, uint32, error)
	// DeleteDaemonSet 删除 DaemonSet
	DeleteDaemonSet(ctx context.Context, req *DaemonSetCommonParams) error
	// CreateOrUpdateDaemonSetByYaml 创建或更新 DaemonSet
	CreateOrUpdateDaemonSetByYaml(ctx context.Context, req *CreateOrUpdateDaemonSetByYamlRequest) error
	// RestartDaemonSet 重启 DaemonSet
	RestartDaemonSet(ctx context.Context, req *DaemonSetCommonParams) error
	// RollbackDaemonSet 回滚 DaemonSet
	RollbackDaemonSet(ctx context.Context, req *RollbackDaemonSetRequest) error
	// GetDaemonSetHistory 获取 DaemonSet 历史版本
	GetDaemonSetHistory(ctx context.Context, req *DaemonSetCommonParams) ([]*appsv1.ControllerRevision, error)
	// GetDaemonSet 获取 DaemonSet
	GetDaemonSet(ctx context.Context, req *DaemonSetCommonParams) (*appsv1.DaemonSet, error)
	// UpdateDaemonSetStrategy 更新 DaemonSet更新策略
	UpdateDaemonSetStrategy(ctx context.Context, req *UpdateDaemonSetStrategyRequest) (bool, error)
}

type DaemonSetUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

// 检查关键词是否在 Pod 名称、状态、命名空间、镜像中
func filterDaemonSetsByKeyword(daemonSets *appsv1.DaemonSetList, keyword string) *appsv1.DaemonSetList {
	result := &appsv1.DaemonSetList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, daemonSet := range daemonSets.Items {
		if utils.MatchString(pattern, daemonSet.Name) ||
			utils.MatchLabels(pattern, daemonSet.Labels) ||
			utils.MatchContainerImages(pattern, daemonSet.Spec.Template.Spec.Containers) {
			result.Items = append(result.Items, daemonSet)
		}
	}
	return result
}

func (x *DaemonSetUseCase) ListDaemonSet(ctx context.Context, req *ListDaemonSetRequest) ([]*appsv1.DaemonSet, uint32, error) {
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
		allFilteredDaemonSets = make([]*appsv1.DaemonSet, 0)
		limit                 = int64(req.PageSize)
		continueToken         = ""
	)

	for {
		daemonSetListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		daemonSets, err := clientSet.AppsV1().DaemonSets(req.Namespace).List(ctx, daemonSetListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询daemonsets失败: %v", err)
			return nil, 0, fmt.Errorf("查询daemonsets失败: %w", err)
		}
		filteredDaemonSets := filterDaemonSetsByKeyword(daemonSets, req.Keyword)
		for _, deployment := range filteredDaemonSets.Items {
			allFilteredDaemonSets = append(allFilteredDaemonSets, &deployment)
		}

		if daemonSets.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = daemonSets.Continue
	}
	if len(allFilteredDaemonSets) > 0 {
		gvks, _, err := scheme.Scheme.ObjectKinds(allFilteredDaemonSets[0])
		if err != nil {
			return nil, 0, fmt.Errorf("查询deployment失败: %w", err)
		}
		if len(gvks) > 0 {
			for _, deployment := range allFilteredDaemonSets {
				deployment.Kind = gvks[0].Kind
				deployment.APIVersion = gvks[0].GroupVersion().String()
			}
		}
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredDaemonSets, uint32(len(allFilteredDaemonSets)), nil
	}
	if len(allFilteredDaemonSets) == 0 {
		return allFilteredDaemonSets, 0, nil
	}
	// 否则分页返回结果
	paginatedDeployments, total := utils.K8sPaginate(allFilteredDaemonSets, req.Page, req.PageSize)
	return paginatedDeployments, total, nil
}

func (x *DaemonSetUseCase) DeleteDaemonSet(ctx context.Context, req *DaemonSetCommonParams) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	return clientSet.AppsV1().DaemonSets(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
}

func (x *DaemonSetUseCase) CreateOrUpdateDaemonSetByYaml(ctx context.Context, req *CreateOrUpdateDaemonSetByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	daemonSet := &appsv1.DaemonSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&daemonSet); err != nil {
		x.log.WithContext(ctx).Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "解析YAML失败, 请检查格式是否正确")
	}
	// 获取当前的daemonSet对象
	_, err = clientSet.AppsV1().DaemonSets(daemonSet.Namespace).Get(ctx, daemonSet.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 创建deployment
			_, err = clientSet.AppsV1().DaemonSets(daemonSet.Namespace).Create(ctx, daemonSet, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("创建daemonSet失败: %w", err)
			}
			return nil
		}
		if k8serrors.IsForbidden(err) {
			return fmt.Errorf("查询daemonSet失败,: %w", err)
		}
		return fmt.Errorf("查询daemonSet失败: %w", err)
	}

	// 更新daemonSet
	_, err = clientSet.AppsV1().DaemonSets(daemonSet.Namespace).Update(ctx, daemonSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新daemonSet失败: %v", err)
		return fmt.Errorf("更新daemonSet失败: %w", err)
	}
	return nil
}

func (x *DaemonSetUseCase) RestartDaemonSet(ctx context.Context, req *DaemonSetCommonParams) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	daemonSet, err := clientSet.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 创建deployment
			_, err = clientSet.AppsV1().DaemonSets(daemonSet.Namespace).Create(ctx, daemonSet, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("创建daemonSet失败: %w", err)
			}
			return nil
		}
		if k8serrors.IsForbidden(err) {
			return fmt.Errorf("查询daemonSet失败,: %w", err)
		}
		return fmt.Errorf("查询daemonSet失败: %w", err)
	}

	// 更新 DaemonSet 的 spec.template.metadata.annotations 触发restart
	if daemonSet.Spec.Template.Annotations == nil {
		daemonSet.Spec.Template.Annotations = make(map[string]string)
	}
	daemonSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	// 更新daemonSet
	_, err = clientSet.AppsV1().DaemonSets(daemonSet.Namespace).Update(ctx, daemonSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("重启daemonSet失败: %v", err)
		return fmt.Errorf("重启daemonSet失败: %w", err)
	}
	return nil
}

func (x *DaemonSetUseCase) RollbackDaemonSet(ctx context.Context, req *RollbackDaemonSetRequest) error {
	//TODO implement me
	panic("implement me")
}

// GetDaemonSetHistory 查询 DaemonSet 的历史版本
func (x *DaemonSetUseCase) GetDaemonSetHistory(ctx context.Context, req *DaemonSetCommonParams) ([]*appsv1.ControllerRevision, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("获取clientSet失败: %w", err)
	}
	// 获取 DaemonSet
	daemonSet, err := clientSet.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("查询daemonSet失败: %w", err)
	}
	// 获取 DaemonSet 的历史版本
	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("获取selector失败: %w", err)
	}
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, fmt.Errorf("查询controllerRevisions失败: %w", err)
	}
	revisions := make([]*appsv1.ControllerRevision, 0)
	for _, revision := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&revision); ownerRef != nil && ownerRef.UID == daemonSet.UID {
			revisions = append(revisions, &revision)
		}
	}
	return revisions, nil
}

func (x *DaemonSetUseCase) GetDaemonSet(ctx context.Context, req *DaemonSetCommonParams) (*appsv1.DaemonSet, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	daemonSet, err := clientSet.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("查询daemonSet失败: %w", err)
	}
	return daemonSet, nil
}

func (x *DaemonSetUseCase) UpdateDaemonSetStrategy(ctx context.Context, req *UpdateDaemonSetStrategyRequest) (bool, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, fmt.Errorf("获取clientSet失败: %w", err)
	}
	daemonSet, err := clientSet.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("查询daemonSet失败: %w", err)
	}
	if daemonSet.Spec.UpdateStrategy.Type == appsv1.OnDeleteDaemonSetStrategyType {
		return true, nil
	}
	// 更新策略
	if &daemonSet.Spec.UpdateStrategy == nil {
		daemonSet.Spec.UpdateStrategy = appsv1.DaemonSetUpdateStrategy{}
	}
	var (
		maxUnavailable *intstr.IntOrString
	)

	if req.MaxUnavailable != "" {
		maxUnavailable, err = parseIntOrPercent(req.MaxUnavailable)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 deployment更新策略 maxUnavailable 失败: %v", err)
			return false, fmt.Errorf("解析 maxUnavailable 失败: %v", err)
		}
	}

	switch req.UpdateStrategyType {
	case uint32(pb.UpdateDaemonSetUpdateStrategyRequest_RollingUpdate):
		daemonSet.Spec.UpdateStrategy = appsv1.DaemonSetUpdateStrategy{
			Type: appsv1.RollingUpdateDaemonSetStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDaemonSet{
				MaxUnavailable: maxUnavailable,
			},
		}
	case uint32(pb.UpdateDaemonSetUpdateStrategyRequest_OnDelete):
		daemonSet.Spec.UpdateStrategy = appsv1.DaemonSetUpdateStrategy{
			Type: appsv1.OnDeleteDaemonSetStrategyType,
		}
	default:
		return false, fmt.Errorf("未知的更新策略类型: %d", req.UpdateStrategyType)
	}
	_, err = clientSet.AppsV1().DaemonSets(req.Namespace).Update(ctx, daemonSet, metav1.UpdateOptions{})
	if err != nil {
		return false, fmt.Errorf("更新daemonSet失败: %w", err)
	}
	return true, nil
}

func NewDaemonSetUseCase(logger log.Logger, cluster IClusterUseCase) *DaemonSetUseCase {
	return &DaemonSetUseCase{cluster: cluster, log: log.NewHelper(log.With(logger, "module", "biz/daemonset"))}
}

func NewIDaemonSetUseCase(x *DaemonSetUseCase) IDaemonSetUseCase {
	return x
}
