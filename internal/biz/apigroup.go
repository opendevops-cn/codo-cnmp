package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

type ListApiGroupRequest struct {
	ClusterName string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

type IApiGroupUseCase interface {
	ListApiGroup(ctx context.Context, req *ListApiGroupRequest) ([]*metav1.APIGroup, uint32, error)
}

type ApiGroupUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *ApiGroupUseCase) ListApiGroup(ctx context.Context, req *ListApiGroupRequest) ([]*metav1.APIGroup, uint32, error) {
	var result []*metav1.APIGroup
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return result, 0, err
	}
	discovery := clientSet.Discovery()
	apiGroups, err := discovery.ServerGroups()
	if err != nil {
		return result, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	filteredApiGroup := x.filterApiGroupByKeyword(apiGroups.Groups, req.Keyword)
	result = append(result, filteredApiGroup...)
	if req.ListAll {
		return result, uint32(len(result)), nil
	}
	if len(result) == 0 {
		return result, 0, nil
	}
	// 否则分页返回结果
	paginatedInstances, total := utils.K8sPaginate(result, req.Page, req.PageSize)
	return paginatedInstances, total, nil

}

func (x *ApiGroupUseCase) filterApiGroupByKeyword(list []metav1.APIGroup, keyword string) []*metav1.APIGroup {
	var result []*metav1.APIGroup
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, apiGroup := range list {
		if utils.MatchString(pattern, apiGroup.Name) || utils.MatchString(pattern, apiGroup.PreferredVersion.GroupVersion) {
			result = append(result, &apiGroup)
		}
	}
	return result
}

func NewApiGroupUseCase(cluster IClusterUseCase, logger log.Logger) *ApiGroupUseCase {
	return &ApiGroupUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/apiGroup")),
	}
}

func NewIApiGroupUseCase(x *ApiGroupUseCase) IApiGroupUseCase {
	return x
}
