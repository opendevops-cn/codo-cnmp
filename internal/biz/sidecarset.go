package biz

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	kruiseappsv1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"regexp"
)

type ListSideCarSetRequest struct {
	ClusterName string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

type GetSideCarSetRequest struct {
	ClusterName string
	Name        string
}

type UpdateSideCarSetStrategyRequest struct {
	ClusterName        string
	Name               string
	Paused             bool
	Partition          string
	MaxUnavailable     string
	UpdateStrategyType uint32
}

type DeleteSideCarSetRequest struct {
	ClusterName string
	Name        string
}

type ISideCarSetUseCase interface {
	ListSideCarSet(ctx context.Context, req *ListSideCarSetRequest) ([]*kruiseappsv1alpha1.SidecarSet, uint32, error)
	GetSideCarSet(ctx context.Context, req *GetSideCarSetRequest) (*kruiseappsv1alpha1.SidecarSet, error)
	UpdateSideCarSetStrategy(ctx context.Context, req *UpdateSideCarSetStrategyRequest) (bool, error)
	DeleteSideCarSet(ctx context.Context, req *DeleteSideCarSetRequest) (bool, error)
}

type SideCarSetUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *SideCarSetUseCase) DeleteSideCarSet(ctx context.Context, req *DeleteSideCarSetRequest) (bool, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, err
	}
	sidecarSet, err := clientSet.AppsV1alpha1().SidecarSets().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return false, fmt.Errorf("没有权限访问SideCarSet: %w", err)
		}
		if k8serrors.IsNotFound(err) {
			return false, fmt.Errorf("SideCarSet不存在: %w", err)
		}
		return false, fmt.Errorf("查询SideCarSet失败: %w", err)
	}
	err = clientSet.AppsV1alpha1().SidecarSets().Delete(ctx, sidecarSet.Name, metav1.DeleteOptions{})
	if err != nil {
		return false, fmt.Errorf("删除SideCarSet失败: %w", err)
	}
	return true, nil
}

func (x *SideCarSetUseCase) UpdateSideCarSetStrategy(ctx context.Context, req *UpdateSideCarSetStrategyRequest) (bool, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, err
	}
	sidecarSet, err := clientSet.AppsV1alpha1().SidecarSets().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return false, fmt.Errorf("没有权限访问SideCarSet: %w", err)
		}
		if k8serrors.IsNotFound(err) {
			return false, fmt.Errorf("SideCarSet不存在: %w", err)
		}
		return false, fmt.Errorf("查询SideCarSet失败: %w", err)
	}
	var (
		maxUnavailable *intstr.IntOrString
		partition      *intstr.IntOrString
	)
	if req.MaxUnavailable != "" {
		maxUnavailable, err = parseIntOrPercent(req.MaxUnavailable)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 maxUnavailable 失败: %v", err)
			return false, fmt.Errorf("解析 maxUnavailable 失败: %v", err)
		}
	}
	if req.Partition != "" {
		partition, err = parseIntOrPercent(req.Partition)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析 partition 失败: %v", err)
			return false, fmt.Errorf("解析 partition 失败: %v", err)
		}
	}
	// 更新策略
	switch req.UpdateStrategyType {
	case uint32(pb.UpdateSideCarSetStrategyRequest_RollingUpdate):
		sidecarSet.Spec.UpdateStrategy = kruiseappsv1alpha1.SidecarSetUpdateStrategy{
			Type:           kruiseappsv1alpha1.RollingUpdateSidecarSetStrategyType,
			MaxUnavailable: maxUnavailable,
			Partition:      partition,
			Paused:         req.Paused,
		}
	case uint32(pb.UpdateSideCarSetStrategyRequest_NotUpdate):
		sidecarSet.Spec.UpdateStrategy = kruiseappsv1alpha1.SidecarSetUpdateStrategy{
			Type: kruiseappsv1alpha1.NotUpdateSidecarSetStrategyType,
		}
	default:
		return false, fmt.Errorf("未知的策略类型: %d", req.UpdateStrategyType)
	}
	_, err = clientSet.AppsV1alpha1().SidecarSets().Update(ctx, sidecarSet, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新 SideCarSet 策略失败: %v", err)
		return false, fmt.Errorf("更新 SideCarSet 策略失败")
	}
	return true, nil
}

func (x *SideCarSetUseCase) GetSideCarSet(ctx context.Context, req *GetSideCarSetRequest) (*kruiseappsv1alpha1.SidecarSet, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	sidecarSet, err := clientSet.AppsV1alpha1().SidecarSets().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return nil, fmt.Errorf("没有权限访问SideCarSet: %w", err)
		}
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("SideCarSet不存在: %w", err)
		}
		return nil, fmt.Errorf("查询SideCarSet失败: %w", err)
	}
	return sidecarSet, nil
}

func (x *SideCarSetUseCase) ListSideCarSet(ctx context.Context, req *ListSideCarSetRequest) ([]*kruiseappsv1alpha1.SidecarSet, uint32, error) {
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

	var (
		allFilteredSetCarSets []*kruiseappsv1alpha1.SidecarSet
		continueToken         = ""
		limit                 = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		sideCarSets, err := clientSet.AppsV1alpha1().SidecarSets().List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询SideCarSets失败: %v", err)
			return nil, 0, fmt.Errorf("查询SideCarSets失败: %w", err)
		}
		filteredCloneSets := filterSideCarSetsByKeyword(sideCarSets, req.Keyword)
		for _, sideCarSet := range filteredCloneSets.Items {
			allFilteredSetCarSets = append(allFilteredSetCarSets, &sideCarSet)
		}

		if sideCarSets.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = sideCarSets.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredSetCarSets, uint32(len(allFilteredSetCarSets)), nil
	}
	if len(allFilteredSetCarSets) == 0 {
		return nil, 0, nil
	}
	// 否则分页返回结果
	paginatedSideCarSets, total := utils.K8sPaginate(allFilteredSetCarSets, req.Page, req.PageSize)
	return paginatedSideCarSets, total, nil

}

func filterSideCarSetsByKeyword(list *kruiseappsv1alpha1.SidecarSetList, keyword string) *kruiseappsv1alpha1.SidecarSetList {
	if keyword == "" {
		return list
	}
	result := &kruiseappsv1alpha1.SidecarSetList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, sideCarSet := range list.Items {
		var containers []corev1.Container
		for _, container := range sideCarSet.Spec.Containers {
			containers = append(containers, container.Container)
		}
		if utils.MatchString(pattern, sideCarSet.Name) ||
			utils.MatchLabels(pattern, sideCarSet.Labels) ||
			utils.MatchContainerImages(pattern, containers) {
			result.Items = append(result.Items, sideCarSet)
		}
	}
	return result

}

func NewSideCarSetUseCase(cluster IClusterUseCase, logger log.Logger) *SideCarSetUseCase {
	return &SideCarSetUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/sidecarset")),
	}
}

func NewISideCarSetUseCase(x *SideCarSetUseCase) ISideCarSetUseCase {
	return x
}
