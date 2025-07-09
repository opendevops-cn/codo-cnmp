package biz

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"regexp"
	"time"

	"codo-cnmp/common/utils"
	"github.com/go-kratos/kratos/v2/log"
	v1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type CRRCommonParams struct {
	ClusterName string
	Namespace   string
}

type CreateCRRRequest struct {
	CRRCommonParams
	PodName    string
	Containers []string
}

type ListCRRRequest struct {
	CRRCommonParams
	Keyword  string
	Page     uint32
	PageSize uint32
	ListAll  bool
}

type DeleteCRRRequest struct {
	CRRCommonParams
	CRRName string
}

type GetCRRRequest struct {
	CRRCommonParams
	CRRName string
}

type ContainerItems struct {
	PodName    string
	Containers []string
}

type BatchCreateCRRRequest struct {
	CRRCommonParams
	ContainerItems []ContainerItems
}

type BatchQueryCRRDetailRequest struct {
	CRRCommonParams
	CRRNames []string
}

type ICRRUseCase interface {
	ListCRR(ctx context.Context, req *ListCRRRequest) ([]*v1alpha1.ContainerRecreateRequest, uint32, error)
	CreateCRR(ctx context.Context, req *CreateCRRRequest) (string, error)
	DeleteCRR(ctx context.Context, req *DeleteCRRRequest) error
	GetCRR(ctx context.Context, req *GetCRRRequest) (*v1alpha1.ContainerRecreateRequest, error)
	BatchCreateCRR(ctx context.Context, req *BatchCreateCRRRequest) ([]string, error)
	BatchQueryCRR(ctx context.Context, req *BatchQueryCRRDetailRequest) ([]*v1alpha1.ContainerRecreateRequest, error)
}

type CRRUseCase struct {
	cluster IClusterUseCase
	redis   *redis.Client
	log     *log.Helper
}

func (x *CRRUseCase) BatchQueryCRR(ctx context.Context, req *BatchQueryCRRDetailRequest) ([]*v1alpha1.ContainerRecreateRequest, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	var CRRList []*v1alpha1.ContainerRecreateRequest
	for _, CRRName := range req.CRRNames {
		CRR, err := clientSet.AppsV1alpha1().ContainerRecreateRequests(req.Namespace).Get(ctx, CRRName, metav1.GetOptions{})
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询CRR失败: %v", err)
			return nil, err
		}
		CRRList = append(CRRList, CRR)
	}
	return CRRList, err
}

func (x *CRRUseCase) BatchCreateCRR(ctx context.Context, req *BatchCreateCRRRequest) ([]string, error) {
	var CRRNameList []string
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return CRRNameList, err
	}
	for _, containerItem := range req.ContainerItems {
		containers := make([]v1alpha1.ContainerRecreateRequestContainer, 0)
		for _, container := range containerItem.Containers {
			containers = append(containers, v1alpha1.ContainerRecreateRequestContainer{
				Name: container,
			})
		}
		CRRName := fmt.Sprintf("%s.%s.%s", req.Namespace, containerItem.PodName, time.Now().Format("20060102150405"))
		CRRRequest := &v1alpha1.ContainerRecreateRequest{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: req.Namespace,
				Name:      CRRName,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ContainerRecreateRequest",
				APIVersion: "apps.kruise.io/v1alpha1",
			},
			Spec: v1alpha1.ContainerRecreateRequestSpec{
				PodName:    containerItem.PodName,
				Containers: containers,
				Strategy: &v1alpha1.ContainerRecreateRequestStrategy{
					FailurePolicy:             v1alpha1.ContainerRecreateRequestFailurePolicyFail,
					OrderedRecreate:           false,
					UnreadyGracePeriodSeconds: pointer.Int64(3),
					MinStartedSeconds:         10,
				},
				ActiveDeadlineSeconds:   pointer.Int64(300),
				TTLSecondsAfterFinished: pointer.Int32(1800),
			},
		}
		_, err = clientSet.AppsV1alpha1().ContainerRecreateRequests(req.Namespace).Create(ctx, CRRRequest, metav1.CreateOptions{})
		if err != nil {
			x.log.WithContext(ctx).Errorf("创建CRR失败: %v", err)
			return CRRNameList, err
		}
		// 设置半个小时过期时间
		cachedCRR := fmt.Sprintf("codo:cnmp:%s:crr:%s", req.Namespace, CRRName)
		_, err := x.redis.Set(ctx, cachedCRR, len(containers), time.Second*1800).Result()
		if err != nil {
			x.log.WithContext(ctx).Errorf("设置CRR缓存失败: %v", err)
		}
		CRRNameList = append(CRRNameList, CRRName)
	}
	return CRRNameList, err
}

func (x *CRRUseCase) GetCRR(ctx context.Context, req *GetCRRRequest) (*v1alpha1.ContainerRecreateRequest, error) {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	CRR, err := clientSet.AppsV1alpha1().ContainerRecreateRequests(req.Namespace).Get(ctx, req.CRRName, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询CRR失败: %v", err)
		return nil, err
	}
	return CRR, err
}

func (x *CRRUseCase) ListCRR(ctx context.Context, req *ListCRRRequest) ([]*v1alpha1.ContainerRecreateRequest, uint32, error) {
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

	continueToken := ""
	limit := int64(req.PageSize)
	var allFilteredCRRList []*v1alpha1.ContainerRecreateRequest

	for {
		ListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		CRRList, err := clientSet.AppsV1alpha1().ContainerRecreateRequests(req.Namespace).List(ctx, ListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询CRR失败: %v", err)
			return nil, 0, fmt.Errorf("查询CRR失败: %w", err)
		}
		filteredCRRList := filterCRRByKeyword(CRRList, req.Keyword)
		for _, CRR := range filteredCRRList.Items {
			allFilteredCRRList = append(allFilteredCRRList, &CRR)
		}

		if CRRList.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = CRRList.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredCRRList, uint32(len(allFilteredCRRList)), nil
	} else {
		if len(allFilteredCRRList) == 0 {
			return nil, 0, nil
		}
		// 否则分页返回结果
		paginatedCRRList, total := utils.K8sPaginate(allFilteredCRRList, req.Page, req.PageSize)
		return paginatedCRRList, total, nil
	}

}

func filterCRRByKeyword(list *v1alpha1.ContainerRecreateRequestList, keyword string) *v1alpha1.ContainerRecreateRequestList {
	if keyword == "" {
		return list
	}
	filteredCRRList := &v1alpha1.ContainerRecreateRequestList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, CRR := range list.Items {
		if utils.MatchString(pattern, CRR.Name) {
			filteredCRRList.Items = append(filteredCRRList.Items, CRR)
		}
	}
	return filteredCRRList

}

func (x *CRRUseCase) CreateCRR(ctx context.Context, req *CreateCRRRequest) (string, error) {
	var CRRName string
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return CRRName, err
	}
	containers := make([]v1alpha1.ContainerRecreateRequestContainer, 0)
	for _, container := range req.Containers {
		containers = append(containers, v1alpha1.ContainerRecreateRequestContainer{
			Name: container,
		})
	}
	CRRName = fmt.Sprintf("%s.%s.%s", req.Namespace, req.PodName, time.Now().Format("20060102150405"))
	CRRRequest := &v1alpha1.ContainerRecreateRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      CRRName,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ContainerRecreateRequest",
			APIVersion: "apps.kruise.io/v1alpha1",
		},
		Spec: v1alpha1.ContainerRecreateRequestSpec{
			Containers: containers,
			PodName:    req.PodName,
			Strategy: &v1alpha1.ContainerRecreateRequestStrategy{
				FailurePolicy:             v1alpha1.ContainerRecreateRequestFailurePolicyFail,
				OrderedRecreate:           false,
				UnreadyGracePeriodSeconds: pointer.Int64(3),
				MinStartedSeconds:         10,
			},
			ActiveDeadlineSeconds:   pointer.Int64(300),
			TTLSecondsAfterFinished: pointer.Int32(1800),
		},
	}
	_, err = clientSet.AppsV1alpha1().ContainerRecreateRequests(req.Namespace).Create(ctx, CRRRequest, metav1.CreateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建CRR失败: %v", err)
		return CRRName, err
	}
	return CRRName, err
}

func (x *CRRUseCase) DeleteCRR(ctx context.Context, req *DeleteCRRRequest) error {
	clientSet, err := x.cluster.GetKruiseClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.AppsV1alpha1().ContainerRecreateRequests(req.Namespace).Delete(ctx, req.CRRName, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除CRR失败: %v", err)
		return err
	}
	return err
}

func NewCRRUseCase(cluster IClusterUseCase, logger log.Logger, redis *redis.Client) *CRRUseCase {
	return &CRRUseCase{
		cluster: cluster,
		redis:   redis,
		log:     log.NewHelper(log.With(logger, "module", "biz/crr")),
	}
}

func NewICRRUseCase(x *CRRUseCase) ICRRUseCase {
	return x
}
