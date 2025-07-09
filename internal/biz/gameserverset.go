package biz

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	appspub "github.com/openkruise/kruise-api/apps/pub"
	kruiseV1beta1 "github.com/openkruise/kruise-api/apps/v1beta1"
	gamekruisev1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	_ "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
	"regexp"
	"strings"
	"time"
)

type GameServerSetCommonParams struct {
	ClusterName string // 集群名称
	Namespace   string // 命名空间
	Name        string // 控制器名称
}

type LisGameServerSetRequest struct {
	GameServerSetCommonParams
	Page     uint32 // 页码
	PageSize uint32 // 每页数量
	Keyword  string // 关键字
	ListAll  bool   // 是否查询所有
}

type CreateOrUpdateGameServerSetByYamlRequest struct {
	GameServerSetCommonParams
	Yaml string // YAML文件内容
}

type DeleteGameServerSetRequest struct {
	GameServerSetCommonParams
}

type RestartGameServerSetRequest struct {
	GameServerSetCommonParams
}

type RollbackGameServerSetRequest struct {
	GameServerSetCommonParams
	Revision uint32 // 回滚版本
}

type ScaleGameServerSetRequest struct {
	GameServerSetCommonParams
	Replicas int32 // 副本数量
}

type GetGameServerSetDetailRequest struct {
	GameServerSetCommonParams
}

type DeleteGameServerSetPodsRequest struct {
	GameServerSetCommonParams
	PodNames     []string // Pod名称列表
	DeletePolicy uint32   // 删除策略
}

type UpdateGameServerSetScaleStrategyRequest struct {
	GameServerSetCommonParams
	Strategy gamekruisev1alpha1.ScaleStrategy // 扩容策略
}

type GameServerSetControllerRevisionParams struct {
	GameServerSetCommonParams
	Revision int64 // 版本号
}

type UpdateGameServerSetUpgradeStrategyRequest struct {
	GameServerSetCommonParams
	UpdateStrategyType    uint32 // 更新策略类型
	PodUpdateStrategyType uint32 // Pod更新策略
	GracePeriodSeconds    uint32 // 优雅关闭时间
	MaxUnavailable        string // 最大不可用副本数
	MaxSurge              string // 最大超载副本数
	Partition             uint32 // 分区
}

type IGameServerSetUseCase interface {
	// ListGameServerSet 查看GameServerSet列表
	ListGameServerSet(ctx context.Context, req *LisGameServerSetRequest) ([]*gamekruisev1alpha1.GameServerSet, uint32, error)
	// CreateOrUpdateGameServerSetByYaml 通过YAML创建或者更新GameServerSet
	CreateOrUpdateGameServerSetByYaml(ctx context.Context, req *CreateOrUpdateGameServerSetByYamlRequest) error
	// DeleteGameServerSet 删除GameServerSet
	DeleteGameServerSet(ctx context.Context, req *DeleteGameServerSetRequest) error
	// RestartGameServerSet 重启GameServerSet
	RestartGameServerSet(ctx context.Context, req *RestartGameServerSetRequest) error
	// RollbackGameServerSet 回滚GameServerSet
	RollbackGameServerSet(ctx context.Context, req *RollbackGameServerSetRequest) error
	// ScaleGameServerSet 伸缩GameServerSet
	ScaleGameServerSet(ctx context.Context, req *ScaleGameServerSetRequest) error
	// GetGameServerSetDetail 获取GameServerSet详情
	GetGameServerSetDetail(ctx context.Context, req *GetGameServerSetDetailRequest) (*gamekruisev1alpha1.GameServerSet, error)
	// DeleteGameServerSetPods 删除GameServerSet的Pod
	// DeleteGameServerSetPods(ctx context.Context, req *DeleteGameServerSetPodsRequest) error
	// UpdateScaleStrategy 流式扩容策略
	UpdateScaleStrategy(ctx context.Context, req *UpdateGameServerSetScaleStrategyRequest) error
	// UpdateUpgradeStrategy 修改升级策略
	UpdateUpgradeStrategy(ctx context.Context, req *UpdateGameServerSetUpgradeStrategyRequest) error
	// ListGameServerSetControllerRevisions 获取GameServerSet的历史版本
	ListGameServerSetControllerRevisions(ctx context.Context, req *GameServerSetCommonParams) ([]*appsv1.ControllerRevision, error)
}

func NewIGameServerSetUseCase(x *GameServerSetUseCase) IGameServerSetUseCase {
	return x
}

type GameServerSetUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *GameServerSetUseCase) ListGameServerSet(ctx context.Context, req *LisGameServerSetRequest) ([]*gamekruisev1alpha1.GameServerSet, uint32, error) {
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
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
		allFilteredObjects = make([]*gamekruisev1alpha1.GameServerSet, 0)
		continueToken      = ""
		limit              = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		gameServerSets, err := kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).List(ctx, ListOptions)
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
				x.log.WithContext(ctx).Errorf("命名空间: %s GameServerSet未找到", req.Namespace)
				return allFilteredObjects, 0, err
			}
			return allFilteredObjects, 0, fmt.Errorf("查询gameserverset失败: %w", err)
		}
		filteredGameServerSets := x.filterObjectsByKeyword(gameServerSets, req.Keyword)
		for _, object := range filteredGameServerSets.Items {
			allFilteredObjects = append(allFilteredObjects, &object)
		}

		if gameServerSets.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = gameServerSets.Continue
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

func (x *GameServerSetUseCase) filterObjectsByKeyword(objects *gamekruisev1alpha1.GameServerSetList, keyword string) *gamekruisev1alpha1.GameServerSetList {
	result := &gamekruisev1alpha1.GameServerSetList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, object := range objects.Items {
		if utils.MatchString(pattern, object.Name) ||
			utils.MatchLabels(pattern, object.Labels) ||
			utils.MatchContainerImages(pattern, object.Spec.GameServerTemplate.Spec.Containers) {
			result.Items = append(result.Items, object)
		}
	}
	return result
}

func (x *GameServerSetUseCase) CreateOrUpdateGameServerSetByYaml(ctx context.Context, req *CreateOrUpdateGameServerSetByYamlRequest) error {
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
			return xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请确认是否已安装kruise-game插件")
		}
		if k8serrors.IsForbidden(err) {
			x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
			return err
		}
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return fmt.Errorf("获取kruise-game client失败: %w", err)
	}
	gameServerSet := &gamekruisev1alpha1.GameServerSet{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&gameServerSet); err != nil {
		//return fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		x.log.WithContext(ctx).Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "解析YAML失败, 请检查格式是否正确")
	}

	// 获取当前的gss对象
	_, err = kruiseGameClientSet.GameV1alpha1().GameServerSets(gameServerSet.Namespace).Get(ctx, gameServerSet.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		// 创建gss
		_, err = kruiseGameClientSet.GameV1alpha1().GameServerSets(gameServerSet.Namespace).Create(ctx, gameServerSet, metav1.CreateOptions{})
		if err != nil {
			x.log.WithContext(ctx).Errorf("创建gameserverset失败: %v", err)
			return fmt.Errorf("创建gameserverset失败: %w", err)
		}
		return nil
	}
	// 更新gss
	_, err = kruiseGameClientSet.GameV1alpha1().GameServerSets(gameServerSet.Namespace).Update(ctx, gameServerSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新gameserverset失败: %v", err)
		return fmt.Errorf("更新gameserverset失败: %w", err)
	}
	return nil
}

func (x *GameServerSetUseCase) DeleteGameServerSet(ctx context.Context, req *DeleteGameServerSetRequest) error {
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
			return xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请确认是否已安装kruise-game插件")
		}
		if k8serrors.IsForbidden(err) {
			x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
			return err
		}
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return fmt.Errorf("获取kruise-game client失败: %w", err)
	}
	err = kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除gameserverset失败: %v", err)
		return fmt.Errorf("删除gameserverset失败: %w", err)
	}
	return nil
}

func (x *GameServerSetUseCase) RestartGameServerSet(ctx context.Context, req *RestartGameServerSetRequest) error {
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
			return xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请确认是否已安装kruise-game插件")
		}
		if k8serrors.IsForbidden(err) {
			x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
			return err
		}
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return fmt.Errorf("获取kruise-game client失败: %w", err)
	}

	gameServerSet, err := kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset失败: %v", err)
		return fmt.Errorf("查询gameserverset失败: %w", err)
	}
	// 更新GameServerSet的注解，标记为重启
	if gameServerSet.Spec.GameServerTemplate.Annotations == nil {
		gameServerSet.Spec.GameServerTemplate.Annotations = make(map[string]string)
	}
	gameServerSet.Spec.GameServerTemplate.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// 更新GameServerSet
	_, err = kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Update(ctx, gameServerSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("重启gameserverset失败: %v", err)
		return fmt.Errorf("重启gameserverset失败: %w", err)
	}
	return nil
}

func (x *GameServerSetUseCase) RollbackGameServerSet(ctx context.Context, req *RollbackGameServerSetRequest) error {
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
			return xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请确认是否已安装kruise-game插件")
		}
		if k8serrors.IsForbidden(err) {
			x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
			return err
		}
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return fmt.Errorf("获取kruise-game client失败: %w", err)
	}

	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取k8s client失败: %v", err)
		return err
	}

	gameServerSet, err := kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset失败: %v", err)
		return fmt.Errorf("查询gameserverset失败: %w", err)
	}
	// 获取历史版本
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: gameServerSet.Status.LabelSelector,
		Limit:         10,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset对应的controller revision失败: %v", err)
		return fmt.Errorf("查询gameserverset对应的controller revision失败: %w", err)
	}
	// 获取对应版本的 ControllerRevision
	var revision *appsv1.ControllerRevision
	for _, rs := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&rs); ownerRef != nil && ownerRef.Name == gameServerSet.Name && rs.Revision == int64(req.Revision) {
			revision = &rs
			break
		}
	}
	if revision == nil {
		return fmt.Errorf("未找到指定版本的 ControllerRevision")
	}

	// 将 CloneSet 的 spec 更新为目标版本的配置
	newGameServerSet := &gamekruisev1alpha1.GameServerSet{}
	err = json.Unmarshal(revision.Data.Raw, &newGameServerSet)
	if err != nil {
		return fmt.Errorf("无法反序列化目标版本的 CloneSet spec: %v", err)
	}

	// 回滚到指定版本
	gameServerSet.Spec.GameServerTemplate = newGameServerSet.Spec.GameServerTemplate
	_, err = kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Update(ctx, gameServerSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("回滚gameserverset失败: %v", err)
		return fmt.Errorf("回滚gameserverset失败: %w", err)
	}
	return nil
}

func (x *GameServerSetUseCase) ScaleGameServerSet(ctx context.Context, req *ScaleGameServerSetRequest) error {
	kruiseGameCientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
			return xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请确认是否已安装kruise-game插件")
		}
		if k8serrors.IsForbidden(err) {
			x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
			return err
		}
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return fmt.Errorf("获取kruise-game client失败: %w", err)
	}
	gameServerSet, err := kruiseGameCientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset失败: %v", err)
		return fmt.Errorf("查询gameserverset失败: %w", err)
	}

	gameServerSet.Spec.Replicas = &req.Replicas
	_, err = kruiseGameCientSet.GameV1alpha1().GameServerSets(req.Namespace).Update(ctx, gameServerSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("伸缩gameserverset失败: %v", err)
		return fmt.Errorf("伸缩gameserverset失败: %w", err)
	}
	return nil
}

func (x *GameServerSetUseCase) GetGameServerSetDetail(ctx context.Context, req *GetGameServerSetDetailRequest) (*gamekruisev1alpha1.GameServerSet, error) {
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if k8serrors.IsNotFound(err) {
		x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
		return nil, err
	}
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return nil, err
	}
	gameServerSet, err := kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset失败: %v", err)
		return nil, fmt.Errorf("查询gameserverset失败: %w", err)
	}
	return gameServerSet, nil
}

func (x *GameServerSetUseCase) UpdateScaleStrategy(ctx context.Context, req *UpdateGameServerSetScaleStrategyRequest) error {
	//TODO implement me
	panic("implement me")
}

func (x *GameServerSetUseCase) UpdateUpgradeStrategy(ctx context.Context, req *UpdateGameServerSetUpgradeStrategyRequest) error {
	clientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	gameServerSet, err := clientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var (
		maxUnavailable *intstr.IntOrString
	)
	if req.MaxUnavailable != "" {
		maxUnavailable, err = parseIntOrPercent(req.MaxUnavailable)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 maxUnavailable 失败: %v", err)
			return fmt.Errorf("解析 maxUnavailable 失败: %v", err)
		}
	}
	var podUpdateStrategy kruiseV1beta1.PodUpdateStrategyType
	switch req.PodUpdateStrategyType {
	case uint32(pb.UpdateGameServerSetUpgradeStrategyRequest_Recreate):
		podUpdateStrategy = kruiseV1beta1.RecreatePodUpdateStrategyType
	case uint32(pb.UpdateGameServerSetUpgradeStrategyRequest_InPlaceIfPossible):
		podUpdateStrategy = kruiseV1beta1.InPlaceIfPossiblePodUpdateStrategyType
	default:
		return fmt.Errorf("不支持的pod更新策略类型: %d", req.PodUpdateStrategyType)
	}
	partition := int32(req.Partition)
	switch req.UpdateStrategyType {
	case uint32(pb.UpdateGameServerSetUpgradeStrategyRequest_RollingUpdate):
		gameServerSet.Spec.UpdateStrategy = gamekruisev1alpha1.UpdateStrategy{
			RollingUpdate: &gamekruisev1alpha1.RollingUpdateStatefulSetStrategy{
				MaxUnavailable: maxUnavailable,
				Partition:      &partition,
				InPlaceUpdateStrategy: &appspub.InPlaceUpdateStrategy{
					GracePeriodSeconds: int32(req.GracePeriodSeconds),
				},
				PodUpdatePolicy: podUpdateStrategy,
			},
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		}
	case uint32(pb.UpdateGameServerSetUpgradeStrategyRequest_OnDelete):
		gameServerSet.Spec.UpdateStrategy = gamekruisev1alpha1.UpdateStrategy{
			RollingUpdate: &gamekruisev1alpha1.RollingUpdateStatefulSetStrategy{
				Partition: &partition,
				InPlaceUpdateStrategy: &appspub.InPlaceUpdateStrategy{
					GracePeriodSeconds: int32(req.GracePeriodSeconds),
				},
				PodUpdatePolicy: podUpdateStrategy,
			},
			Type: appsv1.OnDeleteStatefulSetStrategyType,
		}
	default:
		return fmt.Errorf("不支持的升级策略类型: %d", req.UpdateStrategyType)
	}

	// 更新 gameserverSet
	_, err = clientSet.GameV1alpha1().GameServerSets(req.Namespace).Update(ctx, gameServerSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新gameserverSet升级策略失败: %v", err)
		return fmt.Errorf("更新gameserverSet升级策略失败: %w", err)
	}
	return nil
}

func (x *GameServerSetUseCase) ListGameServerSetControllerRevisions(ctx context.Context, req *GameServerSetCommonParams) ([]*appsv1.ControllerRevision, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	revisions := make([]*appsv1.ControllerRevision, 0)
	// 获取kruise-game client
	kruiseGameClientSet, err := x.cluster.GetKruiseGameClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			x.log.WithContext(ctx).Errorf("请确认是否已安装kruise-game插件: %v", err)
			return revisions, xerr.NewErrCodeMsg(xerr.ErrResourceNotFound, "请确认是否已安装kruise-game插件")
		}
		if k8serrors.IsForbidden(err) {
			x.log.WithContext(ctx).Errorf("KubeConig权限不足: %v", err)
			return revisions, err
		}
		x.log.WithContext(ctx).Errorf("获取kruise-game client失败: %v", err)
		return revisions, fmt.Errorf("获取kruise-game client失败: %w", err)
	}
	gameServerSet, err := kruiseGameClientSet.GameV1alpha1().GameServerSets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset失败: %v", err)
		return nil, fmt.Errorf("查询gameserverset失败: %w", err)
	}
	// 获取controller revision
	controllerRevisions, err := clientSet.AppsV1().ControllerRevisions(req.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: gameServerSet.Status.LabelSelector,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询gameserverset对应的controller revision失败: %v", err)
		return nil, fmt.Errorf("查询gameserverset对应的controller revision失败: %w", err)
	}
	// 过滤controller revision
	for _, rs := range controllerRevisions.Items {
		if ownerRef := metav1.GetControllerOf(&rs); ownerRef != nil && ownerRef.Name == gameServerSet.Name {
			revisions = append(revisions, &rs)
		}
	}
	return revisions, nil
}

func NewGameServerSetUseCase(clusterUseCase IClusterUseCase, logger log.Logger) *GameServerSetUseCase {
	return &GameServerSetUseCase{
		cluster: clusterUseCase,
		log:     log.NewHelper(log.With(logger, "module", "biz/gameserverset")),
	}
}
