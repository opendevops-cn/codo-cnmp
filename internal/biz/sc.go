package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

type ListStorageClassRequest struct {
	ClusterName string
	Keyword     string
	PageSize    uint32
	Page        uint32
	ListAll     bool
}

type IStorageClassUseCase interface {
	ListStorageClass(ctx context.Context, req *ListStorageClassRequest) ([]*storagev1.StorageClass, uint32, error)
}

type StorageClassUseCase struct {
	Cluster IClusterUseCase
	log     *log.Helper
}

func (x *StorageClassUseCase) ListStorageClass(ctx context.Context, req *ListStorageClassRequest) ([]*storagev1.StorageClass, uint32, error) {
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
		allFilteredStorageClasses []*storagev1.StorageClass
		continueToken             = ""
		limit                     = int64(req.PageSize)
	)

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		StorageClasses, err := clientSet.StorageV1().StorageClasses().List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询StoargeClasses失败: %v", err)
			return nil, 0, fmt.Errorf("查询StoargeClasses失败: %w", err)
		}
		filteredStorageClasses := filterStorageClassesByKeyword(StorageClasses, req.Keyword)
		for _, storageClass := range filteredStorageClasses.Items {
			allFilteredStorageClasses = append(allFilteredStorageClasses, &storageClass)
		}

		if StorageClasses.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = StorageClasses.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredStorageClasses, uint32(len(allFilteredStorageClasses)), nil
	}
	if len(allFilteredStorageClasses) == 0 {
		return nil, 0, nil
	}
	// 否则分页返回结果
	paginatedStorageClasses, total := utils.K8sPaginate(allFilteredStorageClasses, req.Page, req.PageSize)
	return paginatedStorageClasses, total, nil

}

func filterStorageClassesByKeyword(classes *storagev1.StorageClassList, keyword string) *storagev1.StorageClassList {
	results := &storagev1.StorageClassList{}
	if keyword == "" {
		return classes
	}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, class := range classes.Items {
		if utils.MatchString(pattern, class.Name) {
			results.Items = append(results.Items, class)
		}
	}
	return results

}

func NewIStorageClassUseCase(x *StorageClassUseCase) IStorageClassUseCase {
	return x
}

func NewStorageClassUseCase(cluster IClusterUseCase, logger log.Logger) *StorageClassUseCase {
	return &StorageClassUseCase{
		Cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/storage_class")),
	}
}
