package informer

import "context"

// ClusterInformer 定义集群 informer 的接口
type ClusterInformer interface {
	// AddCluster 添加集群的 informer
	AddCluster(ctx context.Context, clusterName string) error
	// RemoveCluster 移除集群的 informer
	RemoveCluster(ctx context.Context, clusterName string) error
}
