package biz

import (
	"bytes"
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"regexp"
	"time"
)

// StatefulSetCommonParams StatefulSet 通用参数
type StatefulSetCommonParams struct {
	ClusterName string // 集群名称
	Namespace   string // 命名空间
	Name        string // deployment 名称
}

// ListStatefulSetRequest 获取 StatefulSet 列表请求
type ListStatefulSetRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

// RollbackStatefulSetRequest 回滚 StatefulSet 请求
type RollbackStatefulSetRequest struct {
	StatefulSetCommonParams
	Revision uint32 // 回滚版本
}

// ScaleStatefulSetRequest 伸缩 StatefulSet 请求
type ScaleStatefulSetRequest struct {
	StatefulSetCommonParams
	Replicas uint32 // 副本数量
}

// CreateOrUpdateStatefulSetByYamlRequest 创建或更新 StatefulSet 请求
type CreateOrUpdateStatefulSetByYamlRequest struct {
	ClusterName string // 集群名称
	Yaml        string // YAML文件内容
}

// GetStatefulSetRequest 获取 StatefulSet 请求
type GetStatefulSetRequest struct {
	StatefulSetCommonParams
}

// RestartStatefulSetRequest 重启 StatefulSet 请求
type RestartStatefulSetRequest struct {
	StatefulSetCommonParams
}

// DeleteStatefulSetRequest 删除 StatefulSet 请求
type DeleteStatefulSetRequest struct {
	StatefulSetCommonParams
}

// UpdateStatefulSetStrategyRequest 更新 StatefulSet 策略请求
type UpdateStatefulSetStrategyRequest struct {
	StatefulSetCommonParams
	UpdateStrategyType uint32 // 更新策略类型
	Partition          uint32 // 分区数量
}

type DeleteStatefulSetPodRequest struct {
	StatefulSetCommonParams
	PodNames     []string // Pod 名称列表
	DeletePolicy uint32   // 删除策略
}

type IStatefulSetUseCase interface {
	// ListStatefulSet 获取 StatefulSet 列表
	ListStatefulSet(ctx context.Context, req *ListStatefulSetRequest) ([]*appsv1.StatefulSet, uint32, error)
	// DeleteStatefulSet 删除 StatefulSet
	DeleteStatefulSet(ctx context.Context, req *DeleteStatefulSetRequest) error
	// CreateOrUpdateStatefulSetByYaml 创建或更新 StatefulSet
	CreateOrUpdateStatefulSetByYaml(ctx context.Context, req *CreateOrUpdateStatefulSetByYamlRequest) error
	// RestartStatefulSet 重启 StatefulSet
	RestartStatefulSet(ctx context.Context, req *RestartStatefulSetRequest) error
	// RollbackStatefulSet 回滚 StatefulSet
	RollbackStatefulSet(ctx context.Context, req *RollbackStatefulSetRequest) error
	// ScaleStatefulSet 伸缩 StatefulSet
	ScaleStatefulSet(ctx context.Context, req *ScaleStatefulSetRequest) error
	// GetStatefulSetRevisions 获取 StatefulSet 历史版本
	GetStatefulSetRevisions(ctx context.Context, req *StatefulSetCommonParams) ([]*appsv1.ControllerRevision, error)
	// GetStatefulSetDetail 获取 StatefulSet 详情
	GetStatefulSetDetail(ctx context.Context, req *GetStatefulSetRequest) (*appsv1.StatefulSet, error)
	// GetStatefulSet 获取 StatefulSet
	GetStatefulSet(ctx context.Context, req *GetStatefulSetRequest) (*appsv1.StatefulSet, error)
	// UpdateStatefulSetStrategy 更新 StatefulSet 策略
	UpdateStatefulSetStrategy(ctx context.Context, req *UpdateStatefulSetStrategyRequest) (bool, error)
	// DeleteStatefulSetPod 删除 StatefulSet Pod
	//DeleteStatefulSetPod(ctx context.Context, req *DeleteStatefulSetPodRequest) error
}

type StatefulSetUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *StatefulSetUseCase) UpdateStatefulSetStrategy(ctx context.Context, req *UpdateStatefulSetStrategyRequest) (bool, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, err
	}
	statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return false, fmt.Errorf("查询 StatefulSet 失败: %w", err)
	}
	if statefulSet.Spec.UpdateStrategy.Type == appsv1.OnDeleteStatefulSetStrategyType {
		return true, nil
	}
	partition := int32(req.Partition)
	switch req.UpdateStrategyType {
	case uint32(pb.UpdateStatefulSetUpdateStrategyRequest_RollingUpdate):
		statefulSet.Spec.UpdateStrategy = appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
			RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{
				Partition: &partition,
			},
		}
	case uint32(pb.UpdateStatefulSetUpdateStrategyRequest_OnDelete):
		statefulSet.Spec.UpdateStrategy = appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.OnDeleteStatefulSetStrategyType,
		}
	default:
		return false, fmt.Errorf("无效的策略类型: %d", req.UpdateStrategyType)
	}
	_, err = clientSet.AppsV1().StatefulSets(req.Namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return false, fmt.Errorf("更新 StatefulSet 策略失败: %w", err)
	}
	return true, nil
}

func (x *StatefulSetUseCase) FilterStatefulSetsByKeyword(statefulSets *appsv1.StatefulSetList, keyword string) *appsv1.StatefulSetList {
	result := &appsv1.StatefulSetList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, statefulSet := range statefulSets.Items {
		if utils.MatchString(pattern, statefulSet.Name) ||
			utils.MatchLabels(pattern, statefulSet.Labels) ||
			utils.MatchContainerImages(pattern, statefulSet.Spec.Template.Spec.Containers) {
			result.Items = append(result.Items, statefulSet)
		}
	}
	return result
}

func (x *StatefulSetUseCase) ListStatefulSet(ctx context.Context, req *ListStatefulSetRequest) ([]*appsv1.StatefulSet, uint32, error) {
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
		allFilteredStatefulSets = make([]*appsv1.StatefulSet, 0)
		continueToken           = ""
		limit                   = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		statefulSets, err := clientSet.AppsV1().StatefulSets(req.Namespace).List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查statefulSets失败: %v", err)
			return nil, 0, fmt.Errorf("查statefulSets失败: %w", err)
		}
		filteredStatefulSets := x.FilterStatefulSetsByKeyword(statefulSets, req.Keyword)
		for _, statefulSet := range filteredStatefulSets.Items {
			allFilteredStatefulSets = append(allFilteredStatefulSets, &statefulSet)
		}

		if statefulSets.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = statefulSets.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredStatefulSets, uint32(len(allFilteredStatefulSets)), nil
	}
	if len(allFilteredStatefulSets) == 0 {
		return allFilteredStatefulSets, 0, nil
	}
	// 否则分页返回结果
	paginatedStatefulSets, total := utils.K8sPaginate(allFilteredStatefulSets, req.Page, req.PageSize)
	return paginatedStatefulSets, total, nil
}

func (x *StatefulSetUseCase) DeleteStatefulSet(ctx context.Context, req *DeleteStatefulSetRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.AppsV1().StatefulSets(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("删除 StatefulSet 失败: %w", err)
	}
	return nil
}

func (x *StatefulSetUseCase) CreateOrUpdateStatefulSetByYaml(ctx context.Context, req *CreateOrUpdateStatefulSetByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	statefulSet := &appsv1.StatefulSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(req.Yaml)), 4096)
	if err := decoder.Decode(statefulSet); err != nil {
		x.log.WithContext(ctx).Errorf("解析 YAML 失败: %v", err)
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "解析 YAML 失败")
	}
	// 查询 StatefulSet 是否存在
	_, err = clientSet.AppsV1().StatefulSets(statefulSet.Namespace).Get(ctx, statefulSet.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 创建
			_, err = clientSet.AppsV1().StatefulSets(statefulSet.Namespace).Create(ctx, statefulSet, metav1.CreateOptions{})
			if err != nil {
				x.log.WithContext(ctx).Errorf("创建statefulset失败: %s/%s: %v", statefulSet.Namespace, statefulSet.Name, err)
				return fmt.Errorf("创建statefulset失败: %w", err)
			}
			return nil
		}
		return fmt.Errorf("查询statefulset失败: %w", err)
	}

	// 更新
	_, err = clientSet.AppsV1().StatefulSets(statefulSet.Namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新statefulset失败: %s/%s: %v", statefulSet.Namespace, statefulSet.Name, err)
		return fmt.Errorf("更新statefulset失败: %w", err)
	}
	return nil
}

func (x *StatefulSetUseCase) RestartStatefulSet(ctx context.Context, req *RestartStatefulSetRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("查询 StatefulSet 失败: %w", err)
	}
	// 更新 Annotations 触发重启
	if statefulSet.Spec.Template.Annotations == nil {
		statefulSet.Spec.Template.Annotations = map[string]string{}
	}
	statefulSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = clientSet.AppsV1().StatefulSets(req.Namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("重启statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("重启 StatefulSet 失败: %w", err)
	}
	return nil
}

func (x *StatefulSetUseCase) DecodeControllerRevisionToStatefulSet(controllerRevision *appsv1.ControllerRevision) (*appsv1.StatefulSet, error) {
	statefulSet := &appsv1.StatefulSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(controllerRevision.Data.Raw), 4096)
	if err := decoder.Decode(statefulSet); err != nil {
		return nil, fmt.Errorf("failed to decode ControllerRevision: %v", err)
	}
	return statefulSet, nil
}

func (x *StatefulSetUseCase) RollbackStatefulSet(ctx context.Context, req *RollbackStatefulSetRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	currentStatefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("查询 StatefulSet 失败: %w", err)
	}
	// 获取 StatefulSet 的 ReplicaSet 历史版本
	selector, err := metav1.LabelSelectorAsSelector(currentStatefulSet.Spec.Selector)
	if err != nil {
		return fmt.Errorf("转换 Selector 失败: %v", err)
	}
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset历史版本失败: %v", err)
		return fmt.Errorf("查询历史版本失败: %w", err)
	}
	var revision *appsv1.ControllerRevision
	for _, rs := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&rs); ownerRef != nil && ownerRef.UID == currentStatefulSet.UID && rs.Revision == int64(req.Revision) {
			revision = &rs
			break
		}
	}
	if revision == nil {
		return fmt.Errorf("未找到对应版本的配置信息")
	}
	// 解析 ControllerRevision 数据
	statefulSet, err := x.DecodeControllerRevisionToStatefulSet(revision)
	if err != nil {
		return err
	}
	currentStatefulSet.Spec.Template = statefulSet.Spec.Template
	_, err = clientSet.AppsV1().StatefulSets(req.Namespace).Update(ctx, currentStatefulSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("回滚statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("回滚 StatefulSet 失败: %w", err)
	}
	return nil
}

func (x *StatefulSetUseCase) ScaleStatefulSet(ctx context.Context, req *ScaleStatefulSetRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("查询 StatefulSet 失败: %w", err)
	}
	replicas := int32(req.Replicas)
	statefulSet.Spec.Replicas = &replicas
	_, err = clientSet.AppsV1().StatefulSets(req.Namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("伸缩statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return fmt.Errorf("伸缩 StatefulSet 失败: %w", err)
	}
	return nil
}

func (x *StatefulSetUseCase) GetStatefulSetYaml(ctx context.Context, req *StatefulSetCommonParams) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (x *StatefulSetUseCase) GetStatefulSetRevisions(ctx context.Context, req *StatefulSetCommonParams) ([]*appsv1.ControllerRevision, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return nil, fmt.Errorf("查询 StatefulSet 失败: %w", err)
	}
	// 获取 StatefulSet 的 ReplicaSet 历史版本
	selector, err := metav1.LabelSelectorAsSelector(statefulSet.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("转换 Selector 失败: %v", err)
	}
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset历史版本失败: %v", err)
		return nil, fmt.Errorf("查询 statefulset 历史版本失败: %w", err)
	}
	revisions := make([]*appsv1.ControllerRevision, 0)
	for _, revision := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&revision); ownerRef != nil && ownerRef.UID == statefulSet.UID {
			revisions = append(revisions, &revision)
		}
	}
	return revisions, nil
}

func (x *StatefulSetUseCase) GetStatefulSetDetail(ctx context.Context, req *GetStatefulSetRequest) (*appsv1.StatefulSet, error) {
	//TODO implement me
	panic("implement me")
}

func (x *StatefulSetUseCase) GetStatefulSet(ctx context.Context, req *GetStatefulSetRequest) (*appsv1.StatefulSet, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	statefulSet, err := clientSet.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询statefulset失败: %s/%s: %v", req.Namespace, req.Name, err)
		return nil, fmt.Errorf("查询 StatefulSet 失败: %w", err)
	}
	return statefulSet, nil
}

func NewStatefulSetUseCase(cluster IClusterUseCase, logger log.Logger) *StatefulSetUseCase {
	return &StatefulSetUseCase{cluster: cluster, log: log.NewHelper(log.With(logger, "module", "biz/statefulset"))}
}

func NewIStatefulSetUseCase(x *StatefulSetUseCase) IStatefulSetUseCase {
	return x
}
