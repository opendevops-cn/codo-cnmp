package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"codo-cnmp/common/utils"
	"github.com/go-kratos/kratos/v2/log"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type IngressCommonParams struct {
	ClusterName string
	Namespace   string
	Name        string
}

type ListIngressRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	ListAll     bool
	Host        string
	Page        uint32
	PageSize    uint32
}

type GetIngressRequest struct {
	IngressCommonParams
}

type DeleteIngressRequest struct {
	IngressCommonParams
}

type ListIngressHostRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	ListAll     bool
	Page        uint32
	PageSize    uint32
}

type CreateIngressRequest struct {
	IngressCommonParams
	IngressSpec networkingv1.IngressSpec
	Labels      map[string]string
	Annotations map[string]string
}

type IIngressUseCase interface {
	ListIngress(ctx context.Context, req *ListIngressRequest) ([]*networkingv1.Ingress, uint32, error)
	GetIngress(ctx context.Context, req *GetIngressRequest) (*networkingv1.Ingress, error)
	DeleteIngress(ctx context.Context, req *DeleteIngressRequest) error
	ListIngressHost(ctx context.Context, req *ListIngressHostRequest) ([]string, error)
	// CreateIngress 创建ingress
	CreateIngress(ctx context.Context, req *CreateIngressRequest) error
	// UpdateIngress 更新ingress
	UpdateIngress(ctx context.Context, req *CreateIngressRequest) error
}

type IngressUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *IngressUseCase) ListIngressHost(ctx context.Context, req *ListIngressHostRequest) ([]string, error) {
	ingressReq := &ListIngressRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Keyword:     req.Keyword,
		ListAll:     req.ListAll,
		Page:        req.Page,
	}
	ingresses, _, err := x.ListIngress(ctx, ingressReq)
	if err != nil {
		return nil, err
	}
	hosts := make([]string, 0)
	for _, ingress := range ingresses {
		for _, rule := range ingress.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}
	}
	return hosts, nil
}

func (x *IngressUseCase) CreateIngress(ctx context.Context, req *CreateIngressRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.IngressSpec,
	}
	_, err = clientSet.NetworkingV1().Ingresses(req.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// UpdateIngress 更新ingress
func (x *IngressUseCase) UpdateIngress(ctx context.Context, req *CreateIngressRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}

	// 构建 JSON Patch 操作数组
	var patchOperations []map[string]interface{}

	// 更新 spec
	patchOperations = append(patchOperations, map[string]interface{}{
		"op":    "replace",
		"path":  "/spec/ingressClassName",
		"value": req.IngressSpec.IngressClassName,
	})
	patchOperations = append(patchOperations, map[string]interface{}{
		"op":    "replace",
		"path":  "/spec/rules",
		"value": req.IngressSpec.Rules,
	})
	patchOperations = append(patchOperations, map[string]interface{}{
		"op":    "replace",
		"path":  "/spec/tls",
		"value": req.IngressSpec.TLS,
	})

	// 替换 labels - 无论是否为 nil 都执行替换
	patchOperations = append(patchOperations, map[string]interface{}{
		"op":    "replace",
		"path":  "/metadata/labels",
		"value": req.Labels,
	})

	// 替换 annotations - 无论是否为 nil 都执行替换
	patchOperations = append(patchOperations, map[string]interface{}{
		"op":    "replace",
		"path":  "/metadata/annotations",
		"value": req.Annotations,
	})

	// 序列化 JSON Patch 操作
	patchBytes, err := json.Marshal(patchOperations)
	if err != nil {
		return err
	}

	// 应用 JSON Patch
	_, err = clientSet.NetworkingV1().Ingresses(req.Namespace).Patch(
		ctx,
		req.Name,
		types.JSONPatchType,
		patchBytes,
		metav1.PatchOptions{},
	)
	return err
}

func NewIngressUseCase(cluster IClusterUseCase, logger log.Logger) *IngressUseCase {
	return &IngressUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/ingress")),
	}
}

func (x *IngressUseCase) ListIngress(ctx context.Context, req *ListIngressRequest) ([]*networkingv1.Ingress, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, fmt.Errorf("创建k8s client失败: %w", err)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	var (
		allFilteredIngressList = make([]*networkingv1.Ingress, 0)
		continueToken          = ""
		limit                  = int64(req.PageSize)
	)
	for {
		IngressListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		IngressList, err := clientSet.NetworkingV1().Ingresses(req.Namespace).List(ctx, IngressListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取ingress列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取ingress列表失败: %w", err)
		}
		filteredIngressList := x.filterIngressByKeyword(IngressList, req.Keyword, req.Host)
		for _, ingress := range filteredIngressList.Items {
			allFilteredIngressList = append(allFilteredIngressList, &ingress)
		}
		if IngressList.Continue == "" {
			break
		}
		continueToken = IngressList.Continue
	}
	if req.ListAll {
		return allFilteredIngressList, uint32(len(allFilteredIngressList)), nil
	}
	if len(allFilteredIngressList) == 0 {
		return allFilteredIngressList, 0, nil
	}
	// 否则分页返回结果
	paginatedIngresses, total := utils.K8sPaginate(allFilteredIngressList, req.Page, req.PageSize)
	return paginatedIngresses, total, nil
}

func (x *IngressUseCase) GetIngress(ctx context.Context, req *GetIngressRequest) (*networkingv1.Ingress, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	ingress, err := clientSet.NetworkingV1().Ingresses(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取ingress失败: %v", err)
		return nil, fmt.Errorf("获取ingress失败: %w", err)
	}
	return ingress, nil
}

func (x *IngressUseCase) DeleteIngress(ctx context.Context, req *DeleteIngressRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.NetworkingV1().Ingresses(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除ingress失败: %v", err)
		return fmt.Errorf("删除ingress失败: %w", err)
	}
	return nil
}

func matchIngressRules(ingress networkingv1.Ingress, keyword string) bool {
	if keyword == "" {
		return true
	}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, rule := range ingress.Spec.Rules {
		if utils.MatchString(pattern, rule.Host) {
			return true
		}
	}
	return false
}

func (x *IngressUseCase) filterIngressByKeyword(list *networkingv1.IngressList, keyword string, host string) *networkingv1.IngressList {
	if keyword == "" && host == "" {
		return list
	}
	result := &networkingv1.IngressList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, Ingress := range list.Items {
		if host != "" {
			for _, rule := range Ingress.Spec.Rules {
				if rule.Host == host {
					result.Items = append(result.Items, Ingress)
					break
				}
			}
		} else {
			if matchIngressRules(Ingress, keyword) {
				result.Items = append(result.Items, Ingress)
				continue
			}
			if utils.MatchString(pattern, Ingress.Name) {
				result.Items = append(result.Items, Ingress)
			}
		}

	}
	return result
}

func NewIIngressUseCase(x *IngressUseCase) IIngressUseCase {
	return x
}
