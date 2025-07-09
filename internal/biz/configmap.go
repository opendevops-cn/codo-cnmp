package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"fmt"
	_ "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"regexp"
	"strings"
)

// ConfigMapCommonParams ConfigMap通用参数.
type ConfigMapCommonParams struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

// ListConfigMapRequest 列出ConfigMap请求.
type ListConfigMapRequest struct {
	ClusterName string
	Namespace   string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

// GetConfigMapRequest 获取ConfigMap请求.
type GetConfigMapRequest struct {
	ConfigMapCommonParams
}

// CreateConfigMapRequest 创建ConfigMap请求.
type CreateConfigMapRequest struct {
	ConfigMapCommonParams
	Data        map[string]string `json:"data"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// UpdateConfigMapRequest 更新ConfigMap请求.
type UpdateConfigMapRequest struct {
	ConfigMapCommonParams
	Data        map[string]string `json:"data"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// DeleteConfigMapRequest 删除ConfigMap请求.
type DeleteConfigMapRequest struct {
	ConfigMapCommonParams
}

// CreateOrUpdateConfigMapByYamlRequest 通过Yaml创建ConfigMap请求.
type CreateOrUpdateConfigMapByYamlRequest struct {
	ClusterName string `json:"cluster_name"`
	Yaml        string `json:"yaml"`
}

type IConfigMapUseCase interface {
	ListConfigMap(ctx context.Context, req *ListConfigMapRequest) ([]*corev1.ConfigMap, uint32, error)
	GetConfigMap(ctx context.Context, req *ConfigMapCommonParams) (*corev1.ConfigMap, error)
	CreateConfigMap(ctx context.Context, req *CreateConfigMapRequest) error
	UpdateConfigMap(ctx context.Context, req *UpdateConfigMapRequest) error
	DeleteConfigMap(ctx context.Context, req *DeleteConfigMapRequest) error
	CreateOrUpdateConfigMapByYaml(ctx context.Context, req *CreateOrUpdateConfigMapByYamlRequest) error
}

type ConfigMapUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *ConfigMapUseCase) ListConfigMap(ctx context.Context, req *ListConfigMapRequest) ([]*corev1.ConfigMap, uint32, error) {
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
		allFilteredConfigMap = make([]*corev1.ConfigMap, 0)
		continueToken        = ""
		limit                = int64(req.PageSize)
	)
	for {
		configMapListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		configMapList, err := clientSet.CoreV1().ConfigMaps(req.Namespace).List(ctx, configMapListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取ConfigMap列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取ConfigMap列表失败: %w", err)
		}
		filteredConfigMapList := x.filterConfigMapByKeyword(configMapList, req.Keyword)
		for _, configMap := range filteredConfigMapList.Items {
			allFilteredConfigMap = append(allFilteredConfigMap, &configMap)
		}
		if configMapList.Continue == "" {
			break
		}
		continueToken = configMapList.Continue
	}
	if len(allFilteredConfigMap) > 0 {
		gvks, _, err := scheme.Scheme.ObjectKinds(allFilteredConfigMap[0])
		if err != nil {
			return nil, 0, fmt.Errorf("查询conigmap失败: %w", err)
		}
		if len(gvks) > 0 {
			for _, deployment := range allFilteredConfigMap {
				deployment.Kind = gvks[0].Kind
				deployment.APIVersion = gvks[0].GroupVersion().String()
			}
		}
	}
	if req.ListAll {
		return allFilteredConfigMap, uint32(len(allFilteredConfigMap)), nil
	}
	if len(allFilteredConfigMap) == 0 {
		return allFilteredConfigMap, 0, nil
	}
	// 否则分页返回结果
	paginatedConfigMaps, total := utils.K8sPaginate(allFilteredConfigMap, req.Page, req.PageSize)
	return paginatedConfigMaps, total, nil
}

func (x *ConfigMapUseCase) filterConfigMapByKeyword(configMapList *corev1.ConfigMapList, keyword string) *corev1.ConfigMapList {
	result := &corev1.ConfigMapList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, configMap := range configMapList.Items {
		if utils.MatchString(pattern, configMap.Name) {
			result.Items = append(result.Items, configMap)
		}
	}
	return result
}

func (x *ConfigMapUseCase) GetConfigMap(ctx context.Context, req *ConfigMapCommonParams) (*corev1.ConfigMap, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	configMap, err := clientSet.CoreV1().ConfigMaps(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取ConfigMap失败: %v", err)
		return nil, fmt.Errorf("获取ConfigMap失败: %w", err)
	}
	return configMap, nil
}

func (x *ConfigMapUseCase) CreateConfigMap(ctx context.Context, req *CreateConfigMapRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	// 根据n是否已经存在
	if _, err := clientSet.CoreV1().ConfigMaps(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{}); err == nil {
		x.log.WithContext(ctx).Errorf("同名configMap已存在: %v", err)
		return fmt.Errorf("创建configMap失败: %w", err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: req.Data,
	}
	_, err = clientSet.CoreV1().ConfigMaps(req.Namespace).Create(ctx, configMap, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("DryRun参数校验失败: %w", err)
	}
	_, err = clientSet.CoreV1().ConfigMaps(req.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建ConfigMap失败: %v", err)
		return fmt.Errorf("创建ConfigMap失败: %w", err)
	}
	return nil
}

func (x *ConfigMapUseCase) UpdateConfigMap(ctx context.Context, req *UpdateConfigMapRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: req.Data,
	}
	_, err = clientSet.CoreV1().ConfigMaps(req.Namespace).Update(ctx, configMap, metav1.UpdateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("DryRun参数校验失败: %w", err)
	}
	_, err = clientSet.CoreV1().ConfigMaps(req.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新ConfigMap失败: %v", err)
		return fmt.Errorf("更新ConfigMap失败: %w", err)
	}
	return nil
}

func (x *ConfigMapUseCase) DeleteConfigMap(ctx context.Context, req *DeleteConfigMapRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().ConfigMaps(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除ConfigMap失败: %v", err)
		return fmt.Errorf("删除ConfigMap失败: %w", err)
	}
	return nil
}

func (x *ConfigMapUseCase) CreateOrUpdateConfigMapByYaml(ctx context.Context, req *CreateOrUpdateConfigMapByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	configMap := &corev1.ConfigMap{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&configMap); err != nil {
		x.log.WithContext(ctx).Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		return fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
	}
	_, err = clientSet.CoreV1().ConfigMaps(configMap.Namespace).Get(ctx, configMap.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 如果不存在则创建
			_, err = clientSet.CoreV1().ConfigMaps(configMap.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
			if err != nil {
				x.log.WithContext(ctx).Errorf("创建ConfigMap失败: %v", err)
				return fmt.Errorf("创建ConfigMap失败: %w", err)
			}
			return nil
		}
		return fmt.Errorf("获取ConfigMap失败: %w", err)
	}
	// 如果存在则更新
	_, err = clientSet.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新ConfigMap失败: %v", err)
		return fmt.Errorf("更新ConfigMap失败: %w", err)
	}
	return nil
}

func (x *ConfigMapUseCase) GetConfigMapReferences(ctx context.Context, req *ConfigMapCommonParams) (int32, []map[string]string, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return 0, nil, err
	}

	refMap := make(map[string]struct{})
	continueToken := ""

	for {
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{
			Limit:    100,
			Continue: continueToken,
		})
		if err != nil {
			return 0, nil, err
		}

		// 处理当前页的 pods
		for _, pod := range pods.Items {
			if !isPodReferencingConfigMap(&pod, req.Name) {
				continue
			}
			topOwner, err := x.findTopOwner(ctx, clientSet, pod.Namespace, pod.OwnerReferences)
			if err != nil {
				return 0, nil, err
			}
			if topOwner != nil {
				refMap[fmt.Sprintf("%s/%s", topOwner.Kind, topOwner.Name)] = struct{}{}
			} else {
				refMap[fmt.Sprintf("Pod/%s", pod.Name)] = struct{}{}
			}
		}
		continueToken = pods.Continue
		if continueToken == "" {
			break
		}
	}

	references := make([]string, 0, len(refMap))
	for ref := range refMap {
		references = append(references, ref)
	}

	// 返回结果
	result := make([]map[string]string, len(references))
	for i, ref := range references {
		kind, name, _ := strings.Cut(ref, "/")
		result[i] = map[string]string{
			"kind": kind,
			"name": name,
		}
	}
	return int32(len(refMap)), result, nil
}

// isPodReferencingConfigMap 检查 Pod 是否引用了指定的 ConfigMap
func isPodReferencingConfigMap(pod *corev1.Pod, configMapName string) bool {
	// 检查卷挂载
	for _, volume := range pod.Spec.Volumes {
		if volume.ConfigMap != nil && volume.ConfigMap.Name == configMapName {
			return true
		}
	}

	// 检查所有容器
	for _, container := range pod.Spec.Containers {
		// 检查 envFrom
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == configMapName {
				return true
			}
		}

		// 检查 env
		for _, env := range container.Env {
			if env.ValueFrom != nil &&
				env.ValueFrom.ConfigMapKeyRef != nil &&
				env.ValueFrom.ConfigMapKeyRef.Name == configMapName {
				return true
			}
		}
	}

	return false
}

// findTopOwner 递归查找最顶层的所有者
func (x *ConfigMapUseCase) findTopOwner(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, owners []metav1.OwnerReference) (*metav1.OwnerReference, error) {
	if len(owners) == 0 {
		return nil, nil
	}
	owner := owners[0]
	switch owner.Kind {
	case "ReplicaSet":
		rs, err := clientSet.AppsV1().ReplicaSets(namespace).Get(ctx, owner.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if len(rs.OwnerReferences) > 0 {
			return x.findTopOwner(ctx, clientSet, namespace, rs.OwnerReferences)
		}
	case "StatefulSet", "DaemonSet", "Deployment", "Job", "CronJob", "CloneSet":
		return &owner, nil
	}

	return &owner, nil
}

func NewConfigMapUseCase(cluster IClusterUseCase, logger log.Logger) *ConfigMapUseCase {
	return &ConfigMapUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/configmap")),
	}
}

func NewIConfigMapUseCase(x *ConfigMapUseCase) IConfigMapUseCase {
	return x
}
