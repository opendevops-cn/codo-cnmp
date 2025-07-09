package biz

import (
	"codo-cnmp/common/xerr"
	"codo-cnmp/pb"
	"context"
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"regexp"
	"sort"
	"strings"
	"time"

	"codo-cnmp/common/utils"
	"github.com/go-kratos/kratos/v2/log"
	appspub "github.com/openkruise/kruise-api/apps/pub"
	kruiseappsv1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type CloneSetCommonParams struct {
	ClusterName string // 集群名称
	Namespace   string // 命名空间
	Name        string // 名称
}

type LisCloneSetRequest struct {
	CloneSetCommonParams
	Keyword  string
	Page     uint32
	PageSize uint32
	ListAll  bool
}

type CreateOrUpdateCloneSetByYamlRequest struct {
	ClusterName string
	Yaml        string
}

type DeleteCloneSetRequest struct {
	CloneSetCommonParams
}

type RestartCloneSetRequest struct {
	CloneSetCommonParams
}

type RollbackCloneSetRequest struct {
	CloneSetCommonParams
	Revision uint32 // 回滚到的版本
}

type ScaleCloneSetRequest struct {
	CloneSetCommonParams
	Replicas uint32 // 副本数量
}

type GetCloneSetDetailRequest struct {
	CloneSetCommonParams
}

type DeleteCloneSetPodsRequest struct {
	CloneSetCommonParams
	PodNames     []string // Pod名称列表
	DeletePolicy uint32   // 删除策略
}

type UpdateScaleStrategyRequest struct {
	CloneSetCommonParams
	MinReadySeconds uint32 // 最小就绪时间
	MaxUnavailable  string // 最大不可用副本数
}

type UpdateUpgradeStrategyRequest struct {
	CloneSetCommonParams
	UpdateStrategyType uint32 // 更新策略类型
	GracePeriodSeconds uint32 //
	MaxUnavailable     string // 最大不可用副本数
	MaxSurge           string // 最大超载副本数
}

type ICloneSetUseCase interface {
	// ListCloneSets 查看CloneSet列表
	ListCloneSets(ctx context.Context, req *LisCloneSetRequest) ([]*kruiseappsv1alpha1.CloneSet, uint32, error)
	// CreateOrUpdateCloneSetByYaml 通过YAML创建或者更新CloneSet
	CreateOrUpdateCloneSetByYaml(ctx context.Context, req *CreateOrUpdateCloneSetByYamlRequest) error
	// DeleteCloneSet 删除CloneSet
	DeleteCloneSet(ctx context.Context, req *DeleteCloneSetRequest) error
	// RestartCloneSet 重启CloneSet
	RestartCloneSet(ctx context.Context, req *RestartCloneSetRequest) error
	// RollbackCloneSet 回滚CloneSet
	RollbackCloneSet(ctx context.Context, req *RollbackCloneSetRequest) error
	// ScaleCloneSet 伸缩CloneSet
	ScaleCloneSet(ctx context.Context, req *ScaleCloneSetRequest) error
	// GetCloneSetDetail 获取CloneSet详情
	GetCloneSetDetail(ctx context.Context, req *GetCloneSetDetailRequest) (*kruiseappsv1alpha1.CloneSet, error)
	// DeleteCloneSetPods 删除CloneSet的Pod
	DeleteCloneSetPods(ctx context.Context, req *DeleteCloneSetPodsRequest) error
	// UpdateScaleStrategy 流式扩容策略
	UpdateScaleStrategy(ctx context.Context, req *UpdateScaleStrategyRequest) error
	// UpdateUpgradeStrategy 修改升级策略
	UpdateUpgradeStrategy(ctx context.Context, req *UpdateUpgradeStrategyRequest) error
	// ListCloneSetControllerRevisions 获取CloneSet的历史版本
	ListCloneSetControllerRevisions(ctx context.Context, req *CloneSetCommonParams) ([]*appsv1.ControllerRevision, error)
}

func NewICloneSetUseCase(x *CloneSetUseCase) ICloneSetUseCase {
	return x
}

type CloneSetUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func NewCloneSetUseCase(cluster IClusterUseCase, logger log.Logger) *CloneSetUseCase {
	return &CloneSetUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/cloneSet")),
	}
}

func (x *CloneSetUseCase) ListCloneSets(ctx context.Context, req *LisCloneSetRequest) ([]*kruiseappsv1alpha1.CloneSet, uint32, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	continueToken := ""
	limit := int64(req.PageSize)
	var allFilteredCloneSets []*kruiseappsv1alpha1.CloneSet

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		cloneSets, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询CloneSet失败: %v", err)
			return nil, 0, fmt.Errorf("查询CloneSet失败: %w", err)
		}
		filteredCloneSets := filterCloneSetsByKeyword(cloneSets, req.Keyword)
		for _, deployment := range filteredCloneSets.Items {
			allFilteredCloneSets = append(allFilteredCloneSets, &deployment)
		}

		if cloneSets.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = cloneSets.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredCloneSets, uint32(len(allFilteredCloneSets)), nil
	} else {
		if len(allFilteredCloneSets) == 0 {
			return nil, 0, nil
		}
		// 否则分页返回结果
		paginatedPods, total := utils.K8sPaginate(allFilteredCloneSets, req.Page, req.PageSize)
		return paginatedPods, total, nil
	}
}

// 检查关键词是否在 Pod 名称、状态、命名空间、镜像中
func filterCloneSetsByKeyword(objects *kruiseappsv1alpha1.CloneSetList, keyword string) *kruiseappsv1alpha1.CloneSetList {
	result := &kruiseappsv1alpha1.CloneSetList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, cloneSet := range objects.Items {
		if utils.MatchString(pattern, cloneSet.Name) ||
			utils.MatchLabels(pattern, cloneSet.Labels) ||
			utils.MatchContainerImages(pattern, cloneSet.Spec.Template.Spec.Containers) {
			result.Items = append(result.Items, cloneSet)
		}
	}
	return result
}

func (x *CloneSetUseCase) CreateOrUpdateCloneSetByYaml(ctx context.Context, req *CreateOrUpdateCloneSetByYamlRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet := &kruiseappsv1alpha1.CloneSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&cloneSet); err != nil {
		x.log.WithContext(ctx).Errorf("解析CloneSet YAML失败: %v", err)
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "解析CloneSet YAML失败, 请检查格式是否正确")
	}
	// 获取当前的cloneSet对象
	_, err = clientSet.AppsV1alpha1().CloneSets(cloneSet.Namespace).Get(ctx, cloneSet.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// 创建
		_, err = clientSet.AppsV1alpha1().CloneSets(cloneSet.Namespace).Create(ctx, cloneSet, metav1.CreateOptions{})
		if err != nil {
			x.log.WithContext(ctx).Errorf("创建cloneSet失败: %v", err)
			return fmt.Errorf("创建失败: %w", err)
		}
		return nil
	}
	// 更新
	_, err = clientSet.AppsV1alpha1().CloneSets(cloneSet.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新cloneSet失败: %v", err)
		return fmt.Errorf("更新失败: %w", err)
	}
	return nil
}

// DeleteCloneSet 删除CloneSet
func (x *CloneSetUseCase) DeleteCloneSet(ctx context.Context, req *DeleteCloneSetRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取集群失败: %v", err)
		return err
	}
	err = clientSet.AppsV1alpha1().CloneSets(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除失败: %v", err)
		return fmt.Errorf("删除失败: %w", err)
	}
	return nil
}

// RestartCloneSet 重启CloneSet
func (x *CloneSetUseCase) RestartCloneSet(ctx context.Context, req *RestartCloneSetRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询CloneSet失败: %v", err)
		return fmt.Errorf("查询CloneSet失败: %w", err)
	}
	// 更新 cloneSet 的 spec.template.metadata.annotations 触发restart
	if cloneSet.Spec.Template.Annotations == nil {
		cloneSet.Spec.Template.Annotations = make(map[string]string)
	}
	cloneSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	// 更新 cloneSet
	_, err = clientSet.AppsV1alpha1().CloneSets(req.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("重启 cloneSet 失败: %v", err)
		return fmt.Errorf("重启 cloneSet 失败: %w", err)
	}
	return nil
}

// RollbackCloneSet 回滚CloneSet
func (x *CloneSetUseCase) RollbackCloneSet(ctx context.Context, req *RollbackCloneSetRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	// 获取 CloneSet 对象
	kruiseClientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet, err := kruiseClientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询CloneSet失败: %v", err)
		return fmt.Errorf("查询CloneSet失败: %w", err)
	}
	// 获取 CloneSet 对应的 ControllerRevision
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: cloneSet.Status.LabelSelector,
		Limit:         10,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询CloneSet对应的ControllerRevision失败: %v", err)
		return fmt.Errorf("查询CloneSet对应的ControllerRevision失败: %w", err)
	}
	// 获取对应版本的 ControllerRevision
	var revision *appsv1.ControllerRevision
	for _, rs := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&rs); ownerRef != nil && ownerRef.UID == cloneSet.UID && rs.Revision == int64(req.Revision) {
			revision = &rs
			break
		}
	}
	if revision == nil {
		x.log.WithContext(ctx).Errorf("未找到CloneSet对应版本的ControllerRevision")
		return fmt.Errorf("未找到CloneSet对应版本的ControllerRevision")
	}
	// 将 CloneSet 的 spec 更新为目标版本的配置
	newCloneSet := &kruiseappsv1alpha1.CloneSet{}
	err = json.Unmarshal(revision.Data.Raw, &newCloneSet)
	if err != nil {
		x.log.WithContext(ctx).Errorf("无法反序列化目标版本的 CloneSet spec: %v", err)
		return fmt.Errorf("无法反序列化目标版本的 CloneSet spec: %v", err)
	}
	cloneSet.Spec.Template = newCloneSet.Spec.Template
	// 更新 CloneSet
	_, err = kruiseClientSet.AppsV1alpha1().CloneSets(req.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("回滚CloneSet失败: %v", err)
		return fmt.Errorf("回滚CloneSet失败: %w", err)
	}
	return nil
}

// ScaleCloneSet 伸缩CloneSet
func (x *CloneSetUseCase) ScaleCloneSet(ctx context.Context, req *ScaleCloneSetRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询cloneSet失败: %v", err)
		return fmt.Errorf("查询cloneSet失败: %w", err)
	}
	replicas := int32(req.Replicas)
	cloneSet.Spec.Replicas = &replicas
	_, err = clientSet.AppsV1alpha1().CloneSets(req.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("伸缩cloneSet失败: %v", err)
		return fmt.Errorf("伸缩cloneSet失败: %w", err)
	}
	return nil
}

// GetCloneSetDetail 获取CloneSet详情
func (x *CloneSetUseCase) GetCloneSetDetail(ctx context.Context, req *GetCloneSetDetailRequest) (*kruiseappsv1alpha1.CloneSet, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询cloneSet失败: %v", err)
		return nil, fmt.Errorf("查询cloneSet失败: %w", err)
	}
	return cloneSet, nil
}

// DeleteCloneSetPods 删除CloneSet的Pod
func (x *CloneSetUseCase) DeleteCloneSetPods(ctx context.Context, req *DeleteCloneSetPodsRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询cloneSet失败: %v", err)
		return fmt.Errorf("查询CloneSet 失败: %w", err)
	}
	// 更新 CloneSet 的 scaleStrategy.podsToDelete
	cloneSet.Spec.ScaleStrategy = kruiseappsv1alpha1.CloneSetScaleStrategy{}
	cloneSet.Spec.ScaleStrategy.PodsToDelete = req.PodNames
	// 更新cloneSet replicas
	if req.DeletePolicy == uint32(pb.DeleteCloneSetPodRequest_DELETE_ONLY) {
		currentReplicas := *cloneSet.Spec.Replicas
		newReplicas := currentReplicas - int32(len(req.PodNames))
		cloneSet.Spec.Replicas = &newReplicas
	}

	_, err = clientSet.AppsV1alpha1().CloneSets(req.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新 CloneSet 失败: %v", err)
		return fmt.Errorf("更新 CloneSet 失败: %w", err)
	}
	return nil
}

// UpdateScaleStrategy 流式扩容
func (x *CloneSetUseCase) UpdateScaleStrategy(ctx context.Context, req *UpdateScaleStrategyRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询cloneSet失败: %v", err)
		return fmt.Errorf("查询cloneSet失败: %w", err)
	}
	// 更新 ScaleStrategy
	cloneSet.Spec.ScaleStrategy = kruiseappsv1alpha1.CloneSetScaleStrategy{}

	// 指定 ScaleStrategy.MaxUnavailable
	if req.MaxUnavailable != "" {
		intOrStr := intstr.FromString(req.MaxUnavailable)
		cloneSet.Spec.ScaleStrategy.MaxUnavailable = &intOrStr
	}
	// 指定 ScaleStrategy.MinReadySeconds
	cloneSet.Spec.MinReadySeconds = int32(req.MinReadySeconds)

	// 更新 cloneSet
	_, err = clientSet.AppsV1alpha1().CloneSets(req.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新cloneSet扩容策略失败: %v", err)
		return fmt.Errorf("更新cloneSet扩容策略失败: %w", err)
	}
	return nil
}

// UpdateUpgradeStrategy 修改更新策略
func (x *CloneSetUseCase) UpdateUpgradeStrategy(ctx context.Context, req *UpdateUpgradeStrategyRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	cloneSet, err := clientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询cloneSet失败: %v", err)
		return fmt.Errorf("查询cloneSet失败: %w", err)
	}
	var (
		maxSurge, maxUnavailable *intstr.IntOrString
	)
	if req.MaxSurge != "" {
		maxSurge, err = parseIntOrPercent(req.MaxSurge)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 maxSurge 失败: %v", err)
			return fmt.Errorf("解析 maxSurge 失败: %v", err)
		}
	}
	if req.MaxUnavailable != "" {
		maxUnavailable, err = parseIntOrPercent(req.MaxUnavailable)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 maxUnavailable 失败: %v", err)
			return fmt.Errorf("解析 maxUnavailable 失败: %v", err)
		}
	}

	switch req.UpdateStrategyType {
	case uint32(pb.UpdateUpgradeStrategyRequest_Recreate):
		// 重建策略
		cloneSet.Spec.UpdateStrategy = kruiseappsv1alpha1.CloneSetUpdateStrategy{
			Type: kruiseappsv1alpha1.RecreateCloneSetUpdateStrategyType,
		}
	case uint32(pb.UpdateUpgradeStrategyRequest_InPlaceOnly):
		// 就地升级策略
		cloneSet.Spec.UpdateStrategy = kruiseappsv1alpha1.CloneSetUpdateStrategy{
			Type: kruiseappsv1alpha1.InPlaceOnlyCloneSetUpdateStrategyType,
			InPlaceUpdateStrategy: &appspub.InPlaceUpdateStrategy{
				GracePeriodSeconds: int32(req.GracePeriodSeconds),
			},
		}
		cloneSet.Spec.UpdateStrategy.InPlaceUpdateStrategy.GracePeriodSeconds = int32(req.GracePeriodSeconds)
	case uint32(pb.UpdateUpgradeStrategyRequest_InPlaceIfPossible):
		// 就地升级策略（如果可能）
		cloneSet.Spec.UpdateStrategy = kruiseappsv1alpha1.CloneSetUpdateStrategy{
			Type: kruiseappsv1alpha1.InPlaceIfPossibleCloneSetUpdateStrategyType,
			InPlaceUpdateStrategy: &appspub.InPlaceUpdateStrategy{
				GracePeriodSeconds: int32(req.GracePeriodSeconds),
			},
			MaxUnavailable: maxUnavailable,
			MaxSurge:       maxSurge,
		}

	default:
		return fmt.Errorf("未知的更新策略类型: %d", req.UpdateStrategyType)
	}
	// 更新 cloneSet
	_, err = clientSet.AppsV1alpha1().CloneSets(req.Namespace).Update(ctx, cloneSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新cloneSet升级策略失败: %v", err)
		return fmt.Errorf("更新cloneSet升级策略失败: %w", err)
	}
	return nil
}

// ListCloneSetControllerRevisions 获取CloneSet的历史版本
func (x *CloneSetUseCase) ListCloneSetControllerRevisions(ctx context.Context, req *CloneSetCommonParams) ([]*appsv1.ControllerRevision, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	kruiseClientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	cloneSet, err := kruiseClientSet.AppsV1alpha1().CloneSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if errors.IsForbidden(err) {
		return nil, fmt.Errorf("没有权限查询ControllerRevision: %v", err)
	}
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("未找到cloneSet对应的历史版本: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("查询cloneSet历史版本失败: %w", err)
	}
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: cloneSet.Status.LabelSelector,
		Limit:         10,
	})
	if err != nil {
		return nil, fmt.Errorf("查询CloneSet历史版本失败: %w", err)
	}
	// 过滤出CloneSet对应的ControllerRevision
	revisions := make([]*appsv1.ControllerRevision, 0)
	for _, rs := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&rs); ownerRef != nil && ownerRef.UID == cloneSet.UID {
			revisions = append(revisions, &rs)
		}
	}
	// 按照创建时间排序
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].CreationTimestamp.Time.After(revisions[j].CreationTimestamp.Time)
	})
	return revisions, nil

}
