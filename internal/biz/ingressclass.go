package biz

import (
	"context"
	"fmt"

	"codo-cnmp/common/utils"
	"github.com/go-kratos/kratos/v2/log"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressClassCommonParams struct {
	ClusterName string
}
type ListIngressClassRequest struct {
	IngressClassCommonParams
	Keyword  string
	Page     uint32
	PageSize uint32
	ListAll  bool
}

type IIngressClassUseCase interface {
	ListIngressClass(ctx context.Context, req *ListIngressClassRequest) ([]*networkingv1.IngressClass, uint32, error)
}

type IngressClassUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *IngressClassUseCase) ListIngressClass(ctx context.Context, req *ListIngressClassRequest) ([]*networkingv1.IngressClass, uint32, error) {
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
		allFilteredIngressClassList = make([]*networkingv1.IngressClass, 0)
		continueToken               = ""
		limit                       = int64(req.PageSize)
	)
	for {
		IngressClassListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		IngressClassList, err := clientSet.NetworkingV1().IngressClasses().List(ctx, IngressClassListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取ingress列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取ingress列表失败: %w", err)
		}
		for _, ingressClass := range IngressClassList.Items {
			allFilteredIngressClassList = append(allFilteredIngressClassList, &ingressClass)
		}
		if IngressClassList.Continue == "" {
			break
		}
		continueToken = IngressClassList.Continue
	}
	if req.ListAll {
		return allFilteredIngressClassList, uint32(len(allFilteredIngressClassList)), nil
	}
	if len(allFilteredIngressClassList) == 0 {
		return allFilteredIngressClassList, 0, nil
	}
	// 否则分页返回结果
	paginatedIngresses, total := utils.K8sPaginate(allFilteredIngressClassList, req.Page, req.PageSize)
	return paginatedIngresses, total, nil

}

func NewIngressClassUseCase(cluster IClusterUseCase, logger log.Logger) *IngressClassUseCase {
	return &IngressClassUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/ingressclass")),
	}
}

func NewIIngressClassUseCase(x *IngressClassUseCase) IIngressClassUseCase {
	return x
}
