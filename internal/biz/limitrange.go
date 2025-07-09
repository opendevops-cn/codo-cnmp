package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"regexp"
)

type LimitRangeCommonParams struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

type ListLimitRangeRequest struct {
	LimitRangeCommonParams
	Page     uint32
	PageSize uint32
	ListAll  bool
	Keyword  string
}

type GetLimitRangeRequest struct {
	LimitRangeCommonParams
}

// CreateLimitRangeRequest is the request struct for creating a limit range
type CreateLimitRangeRequest struct {
	LimitRangeCommonParams
	Limits []corev1.LimitRangeItem
}

type UpdateLimitRangeRequest struct {
	LimitRangeCommonParams
	Limits []corev1.LimitRangeItem
}

type DeleteLimitRangeRequest struct {
	LimitRangeCommonParams
}

type ILimitRangeUseCase interface {
	GetLimitRange(ctx context.Context, req *GetLimitRangeRequest) (*corev1.LimitRange, error)
	ListLimitRange(ctx context.Context, req *ListLimitRangeRequest) ([]*corev1.LimitRange, uint32, error)
	CreateLimitRange(ctx context.Context, req *CreateLimitRangeRequest) error
	UpdateLimitRange(ctx context.Context, req *UpdateLimitRangeRequest) error
	DeleteLimitRange(ctx context.Context, req *DeleteLimitRangeRequest) error
	CreateOrUpdateLimitRange(ctx context.Context, req *CreateLimitRangeRequest) error
}

type LimitRangeUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *LimitRangeUseCase) CreateOrUpdateLimitRange(ctx context.Context, req *CreateLimitRangeRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	_, err = clientSet.CoreV1().LimitRanges(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// limitRange 不存在，创建
			return x.CreateLimitRange(ctx, req)
		}
		return fmt.Errorf("查询limitRange失败 %s.%s: %w", req.Namespace, req.Name, err)

	}
	// limitRange已存在，更新
	return x.UpdateLimitRange(ctx, &UpdateLimitRangeRequest{
		LimitRangeCommonParams: req.LimitRangeCommonParams,
		Limits:                 req.Limits,
	})
}

func (x *LimitRangeUseCase) DeleteLimitRange(ctx context.Context, req *DeleteLimitRangeRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().LimitRanges(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除limitRange失败 %s.%s: %w", req.Namespace, req.Name, err)
	}
	return nil
}

func (x *LimitRangeUseCase) GetLimitRange(ctx context.Context, req *GetLimitRangeRequest) (*corev1.LimitRange, error) {
	cluster, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	limitRange, err := cluster.CoreV1().LimitRanges(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询limitRange失败 %s: %w", req.Namespace, err)
	}
	return limitRange, nil
}

func (x *LimitRangeUseCase) ListLimitRange(ctx context.Context, req *ListLimitRangeRequest) ([]*corev1.LimitRange, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
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
		allFilteredLimitRange = make([]*corev1.LimitRange, 0)
		continueToken         = ""
		limit                 = int64(req.PageSize)
	)
	for {
		configMapListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		configMapList, err := clientSet.CoreV1().LimitRanges(req.Namespace).List(ctx, configMapListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取LimitRange列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取LimitRange列表失败: %w", err)
		}
		filteredConfigMapList := x.filterLimitRangeByKeyword(configMapList, req.Keyword)
		for _, configMap := range filteredConfigMapList.Items {
			allFilteredLimitRange = append(allFilteredLimitRange, &configMap)
		}
		if configMapList.Continue == "" {
			break
		}
		continueToken = configMapList.Continue
	}
	if len(allFilteredLimitRange) > 0 {
		gvks, _, err := scheme.Scheme.ObjectKinds(allFilteredLimitRange[0])
		if err != nil {
			return nil, 0, fmt.Errorf("查询limitRange失败: %w", err)
		}
		if len(gvks) > 0 {
			for _, deployment := range allFilteredLimitRange {
				deployment.Kind = gvks[0].Kind
				deployment.APIVersion = gvks[0].GroupVersion().String()
			}
		}
	}
	if req.ListAll {
		return allFilteredLimitRange, uint32(len(allFilteredLimitRange)), nil
	}
	if len(allFilteredLimitRange) == 0 {
		return allFilteredLimitRange, 0, nil
	}
	// 否则分页返回结果
	paginatedConfigMaps, total := utils.K8sPaginate(allFilteredLimitRange, req.Page, req.PageSize)
	return paginatedConfigMaps, total, nil

}

func (x *LimitRangeUseCase) filterLimitRangeByKeyword(limitRangeList *corev1.LimitRangeList, keyword string) *corev1.LimitRangeList {
	if keyword == "" {
		return limitRangeList
	}
	result := &corev1.LimitRangeList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, limitRange := range limitRangeList.Items {
		if utils.MatchString(pattern, limitRange.Name) {
			result.Items = append(result.Items, limitRange)
		}
	}
	return result
}

func (x *LimitRangeUseCase) CreateLimitRange(ctx context.Context, req *CreateLimitRangeRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	limitRange := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: corev1.LimitRangeSpec{
			Limits: req.Limits,
		},
	}
	// 根据n是否已经存在
	if _, err := clientSet.CoreV1().LimitRanges(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{}); err == nil {
		x.log.WithContext(ctx).Errorf("同名limitRange已存在: %v", err)
		return fmt.Errorf("创建limitRange失败: %w", err)
	}
	// Dry-run 校验
	_, err = clientSet.CoreV1().LimitRanges(req.Namespace).Create(ctx, limitRange, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("dry-run校验失败: %w", err)
	}
	// 非 dry-run 模式下创建 limitRange
	_, err = clientSet.CoreV1().LimitRanges(req.Namespace).Create(ctx, limitRange, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建limitRange失败 %s: %w", req.Namespace, err)
	}
	return nil
}

func (x *LimitRangeUseCase) UpdateLimitRange(ctx context.Context, req *UpdateLimitRangeRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(context.Background(), req.ClusterName)
	if err != nil {
		return err
	}
	limitRange, err := clientSet.CoreV1().LimitRanges(req.Namespace).Get(context.Background(), req.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("查询limitRange失败 %s: %w", req.Namespace, err)
	}
	limitRange.Spec.Limits = req.Limits
	// dry-run 校验
	_, err = clientSet.CoreV1().LimitRanges(req.Namespace).Update(ctx, limitRange, metav1.UpdateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("dry-run校验失败: %w", err)
	}
	_, err = clientSet.CoreV1().LimitRanges(req.Namespace).Update(ctx, limitRange, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新limitRange失败 %s: %w", req.Namespace, err)
	}
	return nil
}

func NewLimitRangeUseCase(cluster IClusterUseCase, logger log.Logger) *LimitRangeUseCase {
	return &LimitRangeUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/limitrange")),
	}
}

func NewILimitRangeUseCase(x *LimitRangeUseCase) ILimitRangeUseCase {
	return x
}
