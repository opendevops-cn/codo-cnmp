package biz

import (
	"context"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
)

type CreateOrUpdateResourceRequest struct {
	ClusterName string
	Yaml        string
}

type IResourceUseCase interface {
	// CreateOrUpdateResourceByYaml 创建或更新资源
	CreateOrUpdateResourceByYaml(ctx context.Context, req *CreateOrUpdateResourceRequest) error
	// DryRunResourceByYaml 验证资源是否可以创建或更新
	DryRunResourceByYaml(ctx context.Context, req *CreateOrUpdateResourceRequest) error
}

type ResourceUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *ResourceUseCase) DryRunResourceByYaml(ctx context.Context, req *CreateOrUpdateResourceRequest) error {
	return x.applyResource(ctx, req, true)
}

func (x *ResourceUseCase) CreateOrUpdateResourceByYaml(ctx context.Context, req *CreateOrUpdateResourceRequest) error {
	return x.applyResource(ctx, req, false)
}

func NewResourceUseCase(cluster IClusterUseCase, logger log.Logger) *ResourceUseCase {
	return &ResourceUseCase{cluster: cluster, log: log.NewHelper(log.With(logger, "module", "usecase/resource"))}
}

func NewIResourceUseCase(x *ResourceUseCase) IResourceUseCase {
	return x
}

func (x *ResourceUseCase) applyResource(ctx context.Context, req *CreateOrUpdateResourceRequest, dryRun bool) error {
	dynamicClient, err := x.cluster.GetDynamicClientByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("获取动态客户端失败: %v", err)
	}
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %v", err)
	}

	unstructuredObjs, err := parseYamlToUnstructured(req.Yaml)
	if err != nil {
		return fmt.Errorf("解析 YAML 失败: %v", err)
	}

	for _, unstructuredObj := range unstructuredObjs {
		gvr, err := getGVRMapping(ctx, clientSet, unstructuredObj)
		if err != nil {
			return err
		}

		client, err := getDynamicClientInterface(dynamicClient, gvr, unstructuredObj.GetNamespace())
		if err != nil {
			return err
		}

		if err := createOrUpdateResource(ctx, client, unstructuredObj, dryRun, x.log); err != nil {
			return err
		}
	}
	return nil
}

// 解析 YAML 并转换为 Unstructured 对象
func parseYamlToUnstructured(yamlContent string) ([]*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(yamlContent), 4096)
	var unstructuredObjs []*unstructured.Unstructured

	for {
		var rawObj map[string]interface{}
		if err := decoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("解析 YAML 失败: %v", err)
		}
		if len(rawObj) == 0 {
			continue
		}

		unstructuredObjs = append(unstructuredObjs, &unstructured.Unstructured{Object: rawObj})
	}
	return unstructuredObjs, nil
}

// 获取 GVR 资源映射
func getGVRMapping(ctx context.Context, clientSet *kubernetes.Clientset, obj *unstructured.Unstructured) (schema.GroupVersionResource, error) {
	gvk := obj.GroupVersionKind()
	mapping, err := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(clientSet.Discovery())).RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return mapping.Resource, nil
}

// 获取动态客户端接口
func getDynamicClientInterface(dynamicClient dynamic.Interface, gvr schema.GroupVersionResource, namespace string) (dynamic.ResourceInterface, error) {
	resourceInterface := dynamicClient.Resource(gvr)
	if namespace == "" {
		return resourceInterface, nil
	}
	return resourceInterface.Namespace(namespace), nil
}

// 创建或更新资源
func createOrUpdateResource(ctx context.Context, client dynamic.ResourceInterface, obj *unstructured.Unstructured, dryRun bool, logger *log.Helper) error {
	name := obj.GetName()
	jsonData, err := obj.MarshalJSON()
	if err != nil {
		return fmt.Errorf("序列化 JSON 失败: %w", err)
	}

	_, err = client.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, createErr := client.Create(ctx, obj, metav1.CreateOptions{
				DryRun: getDryRunOption(dryRun),
			})
			if createErr != nil {
				logger.WithContext(ctx).Errorf("资源创建失败: %v", createErr)
				return fmt.Errorf("资源创建失败: %v", createErr)
			}
			return nil
		}
		logger.WithContext(ctx).Errorf("查询失败: %v", err)
		return fmt.Errorf("查询失败: %v", err)
	}

	_, patchErr := client.Patch(ctx, name, types.MergePatchType, jsonData, metav1.PatchOptions{
		DryRun: getDryRunOption(dryRun),
	})
	if patchErr != nil {
		logger.WithContext(ctx).Errorf("资源更新失败: %v", patchErr)
		return fmt.Errorf("资源更新失败: %v", patchErr)
	}
	return nil
}

func getDryRunOption(dryRun bool) []string {
	if dryRun {
		return []string{metav1.DryRunAll}
	}
	return nil
}
