package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

type ListPersistentVolumeRequest struct {
	ClusterName string
	Keyword     string
	PageSize    uint32
	Page        uint32
	ListAll     bool
}

type IPersistentVolumeUseCase interface {
	// ListPersistentVolume list PersistentVolume
	ListPersistentVolume(ctx context.Context, req *ListPersistentVolumeRequest) ([]*corev1.PersistentVolume, uint32, error)
}

type PersistentVolumeUseCase struct {
	Cluster IClusterUseCase
	log     *log.Helper
}

func (x *PersistentVolumeUseCase) ListPersistentVolume(ctx context.Context, req *ListPersistentVolumeRequest) ([]*corev1.PersistentVolume, uint32, error) {
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
		allFilteredPersistentVolumes []*corev1.PersistentVolume
		continueToken                = ""
		limit                        = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		PersistentVolumes, err := clientSet.CoreV1().PersistentVolumes().List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询PersistentVolumes失败: %v", err)
			return nil, 0, fmt.Errorf("查询PersistentVolumes失败: %w", err)
		}
		filteredPersistentVolumes := filterPersistentVolumesByKeyword(PersistentVolumes, req.Keyword)
		for _, pv := range filteredPersistentVolumes.Items {
			allFilteredPersistentVolumes = append(allFilteredPersistentVolumes, &pv)
		}

		if PersistentVolumes.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = PersistentVolumes.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredPersistentVolumes, uint32(len(allFilteredPersistentVolumes)), nil
	}
	if len(allFilteredPersistentVolumes) == 0 {
		return nil, 0, nil
	}
	// 否则分页返回结果
	paginatedPersistentVolumes, total := utils.K8sPaginate(allFilteredPersistentVolumes, req.Page, req.PageSize)
	return paginatedPersistentVolumes, total, nil

}

func filterPersistentVolumesByKeyword(list *corev1.PersistentVolumeList, keyword string) *corev1.PersistentVolumeList {
	results := &corev1.PersistentVolumeList{}
	if keyword == "" {
		return list
	}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, pv := range list.Items {
		if utils.MatchString(pattern, pv.Name) {
			results.Items = append(results.Items, pv)
		}
	}
	return results

}

func NewIPersistentVolumeUseCase(x *PersistentVolumeUseCase) IPersistentVolumeUseCase {
	return x
}

func NewPersistentVolumeUseCase(cluster IClusterUseCase, logger log.Logger) *PersistentVolumeUseCase {
	return &PersistentVolumeUseCase{
		Cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/persistent_volume")),
	}
}
