package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"regexp"
	"strings"
)

type HpaCommonParams struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

// ListHpaRequest 查询 HPA 列表请求
type ListHpaRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

// GetHpaRequest 查询指定 HPA 请求
type GetHpaRequest struct {
	HpaCommonParams
}

// CreateOrUpdateHpaByYamlRequest 通过YAML创建或者更新HPA 请求
type CreateOrUpdateHpaByYamlRequest struct {
	ClusterName string
	Yaml        string
}

// DeleteHpaRequest 删除指定 HPA 请求
type DeleteHpaRequest struct {
	HpaCommonParams
}

type IHpaUseCase interface {
	// GetHpa 查询指定 HPA
	GetHpa(ctx context.Context, req *GetHpaRequest) (*autoscalingv2.HorizontalPodAutoscaler, error)
	// ListHpa 查询 HPA 列表
	ListHpa(ctx context.Context, req *ListHpaRequest) ([]*autoscalingv2.HorizontalPodAutoscaler, uint32, error)
	// CreateOrUpdateHpaByYaml 通过YAML创建或者更新HPA
	CreateOrUpdateHpaByYaml(ctx context.Context, req *CreateOrUpdateHpaByYamlRequest) error
	// DeleteHpa 删除指定 HPA
	DeleteHpa(ctx context.Context, req *DeleteHpaRequest) error
}

type HpaUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func NewHpaUseCase(cluster IClusterUseCase, logger log.Logger) *HpaUseCase {
	return &HpaUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/hpa")),
	}
}

func NewIHpaUseCase(x *HpaUseCase) IHpaUseCase {
	return x
}

func (x *HpaUseCase) GetHpa(ctx context.Context, req *GetHpaRequest) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	hpa, err := clientSet.AutoscalingV2().HorizontalPodAutoscalers(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取HPA失败: %v", err)
		return nil, fmt.Errorf("获取HPA失败: %w", err)
	}
	return hpa, nil

}

func (x *HpaUseCase) ListHpa(ctx context.Context, req *ListHpaRequest) ([]*autoscalingv2.HorizontalPodAutoscaler, uint32, error) {
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
		allFilteredHpaList = make([]*autoscalingv2.HorizontalPodAutoscaler, 0)
		continueToken      = ""
		limit              = int64(req.PageSize)
	)
	for {
		hpaListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		hpaList, err := clientSet.AutoscalingV2().HorizontalPodAutoscalers(req.Namespace).List(ctx, hpaListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取HPA列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取HPA列表失败: %w", err)
		}
		filteredHpaList := x.filterHpaByKeyword(hpaList, req.Keyword)
		for _, hpa := range filteredHpaList.Items {
			allFilteredHpaList = append(allFilteredHpaList, &hpa)
		}
		if hpaList.Continue == "" {
			break
		}
		continueToken = hpaList.Continue
	}
	if req.ListAll {
		return allFilteredHpaList, uint32(len(allFilteredHpaList)), nil
	}
	if len(allFilteredHpaList) == 0 {
		return allFilteredHpaList, 0, nil
	}
	// 否则分页返回结果
	paginatedHpas, total := utils.K8sPaginate(allFilteredHpaList, req.Page, req.PageSize)
	return paginatedHpas, total, nil
}

// filterHpaByKeyword 根据关键字过滤 HPA 列表, 模糊查询hpa名称、工作负载名称、 标签
func (x *HpaUseCase) filterHpaByKeyword(hpaList *autoscalingv2.HorizontalPodAutoscalerList, keyword string) *autoscalingv2.HorizontalPodAutoscalerList {
	result := &autoscalingv2.HorizontalPodAutoscalerList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, hpa := range hpaList.Items {
		if utils.MatchString(pattern, hpa.Name) ||
			utils.MatchString(pattern, hpa.Spec.ScaleTargetRef.Name) ||
			utils.MatchLabels(pattern, hpa.Labels) {
			result.Items = append(result.Items, hpa)
		}
	}
	return result
}

func (x *HpaUseCase) CreateOrUpdateHpaByYaml(ctx context.Context, req *CreateOrUpdateHpaByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	hpa := &autoscalingv2.HorizontalPodAutoscaler{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&hpa); err != nil {
		x.log.WithContext(ctx).Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		return fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
	}
	_, err = clientSet.AutoscalingV2().HorizontalPodAutoscalers(hpa.Namespace).Get(ctx, hpa.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// 新建
			_, err = clientSet.AutoscalingV2().HorizontalPodAutoscalers(hpa.Namespace).Create(ctx, hpa, metav1.CreateOptions{})
			if err != nil {
				x.log.WithContext(ctx).Errorf("创建HPA失败: %v", err)
				return fmt.Errorf("创建HPA失败: %w", err)
			}
			return nil
		}
		x.log.WithContext(ctx).Errorf("获取HPA失败: %v", err)
		return fmt.Errorf("查询HPA失败: %w", err)
	}
	// 更新
	_, err = clientSet.AutoscalingV2().HorizontalPodAutoscalers(hpa.Namespace).Update(ctx, hpa, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新HPA失败: %v", err)
		return fmt.Errorf("更新HPA失败: %w", err)
	}
	return nil
}

func (x *HpaUseCase) DeleteHpa(ctx context.Context, req *DeleteHpaRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.AutoscalingV2().HorizontalPodAutoscalers(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除HPA失败: %v", err)
		return fmt.Errorf("删除HPA失败: %w", err)
	}
	return nil
}
