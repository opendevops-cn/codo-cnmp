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

type ListPersistentVolumeClaimRequest struct {
	ClusterName string
	NameSpace   string
	Keyword     string
	PageSize    uint32
	Page        uint32
	ListAll     bool
}

type DeletePersistentVolumeClaimRequest struct {
	ClusterName string
	NameSpace   string
	Name        string
}

type IPersistentVolumeClaimUseCase interface {
	// ListPersistentVolumeClaim list PersistentVolumeClaim
	ListPersistentVolumeClaim(ctx context.Context, req *ListPersistentVolumeClaimRequest) ([]*corev1.PersistentVolumeClaim, uint32, error)
	// DeletePersistentVolumeClaim delete PersistentVolumeClaim
	DeletePersistentVolumeClaim(ctx context.Context, req *DeletePersistentVolumeClaimRequest) error
}

type PersistentVolumeClaimUseCase struct {
	Cluster IClusterUseCase
	log     *log.Helper
}

func (x *PersistentVolumeClaimUseCase) DeletePersistentVolumeClaim(ctx context.Context, req *DeletePersistentVolumeClaimRequest) error {
	clientSet, err := x.Cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().PersistentVolumeClaims(req.NameSpace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除PersistentVolumeClaim失败: %v", err)
		return fmt.Errorf("删除PersistentVolumeClaim失败: %w", err)
	}
	return nil
}

func (x *PersistentVolumeClaimUseCase) ListPersistentVolumeClaim(ctx context.Context, req *ListPersistentVolumeClaimRequest) ([]*corev1.PersistentVolumeClaim, uint32, error) {
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
		allFilteredPersistentVolumeClaims []*corev1.PersistentVolumeClaim
		continueToken                     = ""
		limit                             = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		PersistentVolumeClaims, err := clientSet.CoreV1().PersistentVolumeClaims(req.NameSpace).List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询PersistentVolumeClaims失败: %v", err)
			return nil, 0, fmt.Errorf("查询PersistentVolumeClaims失败: %w", err)
		}
		filteredPersistentVolumeClaims := filterPersistentVolumeClaimsByKeyword(PersistentVolumeClaims, req.Keyword)
		for _, pv := range filteredPersistentVolumeClaims.Items {
			allFilteredPersistentVolumeClaims = append(allFilteredPersistentVolumeClaims, &pv)
		}

		if PersistentVolumeClaims.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = PersistentVolumeClaims.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredPersistentVolumeClaims, uint32(len(allFilteredPersistentVolumeClaims)), nil
	}
	if len(allFilteredPersistentVolumeClaims) == 0 {
		return nil, 0, nil
	}
	// 否则分页返回结果
	paginatedPersistentVolumeClaims, total := utils.K8sPaginate(allFilteredPersistentVolumeClaims, req.Page, req.PageSize)
	return paginatedPersistentVolumeClaims, total, nil

}

func filterPersistentVolumeClaimsByKeyword(list *corev1.PersistentVolumeClaimList, keyword string) *corev1.PersistentVolumeClaimList {
	results := &corev1.PersistentVolumeClaimList{}
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

func NewIPersistentVolumeClaimUseCase(x *PersistentVolumeClaimUseCase) IPersistentVolumeClaimUseCase {
	return x
}

func NewPersistentVolumeClaimUseCase(cluster IClusterUseCase, logger log.Logger) *PersistentVolumeClaimUseCase {
	return &PersistentVolumeClaimUseCase{
		Cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/persistent_volume")),
	}
}
