package informer

import (
	"context"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewInformerList, NewNodeInformer, NewPodInformer, NewMetricsInformerList, NewMetricsInformer)

type K8sInformerList []IInformer
type MetricsInformerList []IMetricsInformer

// IInformer 资源informer接口
type IInformer interface {
	// Add 添加informer
	Add(ctx context.Context, clusterName string) error
	// Remove 删除informer
	Remove(clusterName string)
	// Start 启动informer
	Start(ctx context.Context) error
	// Stop 停止informer
	Stop(ctx context.Context) error

	Name() string
}

// IMetricsInformer metrics informer接口
type IMetricsInformer interface {
	// Add 添加metrics informer
	Add(ctx context.Context, clusterName string) error
	// Remove 删除metrics informer
	Remove(clusterName string)
	// Start 启动metrics informer
	Start(ctx context.Context) error
	// Stop 停止metrics informer
	Stop(ctx context.Context) error
}

// NewInformerList informer列表
func NewInformerList(podInformer *PodInformer) K8sInformerList {
	return K8sInformerList{podInformer}
}

// NewMetricsInformerList metrics informer列表
func NewMetricsInformerList(metricsInformer *MetricsInformer) MetricsInformerList {
	return MetricsInformerList{metricsInformer}
}
