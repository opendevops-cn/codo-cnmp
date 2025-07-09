package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"regexp"
)

type CRDCommonParams struct {
	ClusterName string
	Namespace   string
	Name        string
}

type ListCRDRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	ApiGroup    string
	ListAll     bool
	Page        uint32
	PageSize    uint32
}

type GetCRDRequest struct {
	CRDCommonParams
}

type DeleteCRDRequest struct {
	CRDCommonParams
}

type ListCRDInstancesRequest struct {
	ClusterName string
	Namespace   string
	Name        string
	ApiGroup    string
	ApiVersion  string
	Keyword     string
	ListAll     bool
	Page        uint32
	PageSize    uint32
}

type ICRDUseCase interface {
	// ListCRD CRD列表
	ListCRD(ctx context.Context, req *ListCRDRequest) ([]*apiextensionsv1.CustomResourceDefinition, uint32, error)
	// GetCRD CRD详情
	GetCRD(ctx context.Context, req *GetCRDRequest) (*apiextensionsv1.CustomResourceDefinition, error)
	// DeleteCRD 删除CRD
	DeleteCRD(ctx context.Context, req *DeleteCRDRequest) error
	// ListCRDInstances CRD实例列表
	ListCRDInstances(ctx context.Context, req *ListCRDInstancesRequest) ([]*unstructured.Unstructured, uint32, error)
}

type CRDUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func NewCRDUseCase(cluster IClusterUseCase, logger log.Logger) *CRDUseCase {
	return &CRDUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/crd")),
	}
}

func (x *CRDUseCase) ListCRD(ctx context.Context, req *ListCRDRequest) ([]*apiextensionsv1.CustomResourceDefinition, uint32, error) {
	dynamicClient, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, fmt.Errorf("创建dynamic client失败: %w", err)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var (
		allFilteredCRDList = make([]*apiextensionsv1.CustomResourceDefinition, 0)
		continueToken      = ""
		limit              = int64(req.PageSize)
	)
	gvr := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	for {
		CRDListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		unstructuredList, err := dynamicClient.Resource(gvr).List(ctx, CRDListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取CRD列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取CRD列表失败: %w", err)
		}
		crdList := &apiextensionsv1.CustomResourceDefinitionList{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredList.UnstructuredContent(), crdList)
		if err != nil {
			x.log.WithContext(ctx).Errorf("解析CRD列表失败: %v", err)
			return nil, 0, fmt.Errorf("解析CRD列表失败: %w", err)
		}
		filteredCRDList := x.filterCRDByKeyword(crdList, req.Keyword, req.ApiGroup)
		for _, crd := range filteredCRDList.Items {
			allFilteredCRDList = append(allFilteredCRDList, &crd)
		}
		if unstructuredList.GetContinue() == "" {
			break
		}
		continueToken = unstructuredList.GetContinue()
	}
	if req.ListAll {
		return allFilteredCRDList, uint32(len(allFilteredCRDList)), nil
	}
	if len(allFilteredCRDList) == 0 {
		return allFilteredCRDList, 0, nil
	}
	// 否则分页返回结果
	paginatedCRDes, total := utils.K8sPaginate(allFilteredCRDList, req.Page, req.PageSize)
	return paginatedCRDes, total, nil
}

func (x *CRDUseCase) GetCRD(ctx context.Context, req *GetCRDRequest) (*apiextensionsv1.CustomResourceDefinition, error) {
	dynamicClient, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	gvr := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	crd, err := dynamicClient.Resource(gvr).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取crd失败: %v", err)
		return nil, fmt.Errorf("获取crd失败: %w", err)
	}
	crdObj := &apiextensionsv1.CustomResourceDefinition{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(crd.UnstructuredContent(), crdObj)
	if err != nil {
		x.log.WithContext(ctx).Errorf("转换CRD对象失败: %v", err)
		return nil, fmt.Errorf("转换CRD对象失败: %w", err)
	}
	return crdObj, nil
}

func (x *CRDUseCase) DeleteCRD(ctx context.Context, req *DeleteCRDRequest) error {
	dynamicClient, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	gvr := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	_, err = dynamicClient.Resource(gvr).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("获取CRD失败: %w", err)
	}
	err = dynamicClient.Resource(gvr).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除CRD失败: %w", err)
	}
	return nil
}

func (x *CRDUseCase) filterCRDByKeyword(list *apiextensionsv1.CustomResourceDefinitionList, keyword string, apiGroup string) *apiextensionsv1.CustomResourceDefinitionList {
	if keyword == "" && apiGroup == "" {
		return list
	}
	result := &apiextensionsv1.CustomResourceDefinitionList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, CRD := range list.Items {
		name := CRD.GetName()
		if utils.MatchString(pattern, name) {
			if apiGroup == "" || apiGroup == CRD.Spec.Group {
				result.Items = append(result.Items, CRD)
			}
		}
	}
	return result
}

func (x *CRDUseCase) ListCRDInstances(ctx context.Context, req *ListCRDInstancesRequest) ([]*unstructured.Unstructured, uint32, error) {
	dynamicClient, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, fmt.Errorf("创建dynamic client失败: %w", err)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var (
		allFilteredInstances = make([]*unstructured.Unstructured, 0)
		continueToken        = ""
		limit                = int64(req.PageSize)
	)
	gr := schema.ParseGroupResource(req.Name)
	gvr := schema.GroupVersionResource{
		Group:    req.ApiGroup,
		Version:  req.ApiVersion,
		Resource: gr.Resource,
	}
	for {
		instanceListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		unstructuredList, err := dynamicClient.Resource(gvr).Namespace(req.Namespace).List(ctx, instanceListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取CRD实例列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取CRD实例列表失败: %w", err)
		}
		filteredUnstructuredList := x.filterCRDInstanceByKeyword(unstructuredList, req.Keyword)
		for _, unstructured := range filteredUnstructuredList.Items {
			allFilteredInstances = append(allFilteredInstances, &unstructured)
		}
		if unstructuredList.GetContinue() == "" {
			break
		}
		continueToken = unstructuredList.GetContinue()
	}
	if req.ListAll {
		return allFilteredInstances, uint32(len(allFilteredInstances)), nil
	}
	if len(allFilteredInstances) == 0 {
		return allFilteredInstances, 0, nil
	}
	// 否则分页返回结果
	paginatedInstances, total := utils.K8sPaginate(allFilteredInstances, req.Page, req.PageSize)
	return paginatedInstances, total, nil
}

func (x *CRDUseCase) filterCRDInstanceByKeyword(list *unstructured.UnstructuredList, keyword string) *unstructured.UnstructuredList {
	if keyword == "" {
		return list
	}
	result := &unstructured.UnstructuredList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, instance := range list.Items {
		name := instance.GetName()
		if utils.MatchString(pattern, name) {
			result.Items = append(result.Items, instance)
		}

	}
	return result
}

func NewICRDUseCase(x *CRDUseCase) ICRDUseCase {
	return x
}
