package informer

import (
	"codo-cnmp/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"sync"
	"time"
)

type NodeInformer struct {
	cluster           *biz.ClusterUseCase
	node              *biz.NodeUseCase
	log               *log.Helper
	informerFactories map[string]informers.SharedInformerFactory
	mu                sync.RWMutex
}

func NewNodeInformer(cluster *biz.ClusterUseCase, node *biz.NodeUseCase, logger log.Logger) *NodeInformer {
	return &NodeInformer{
		cluster:           cluster,
		node:              node,
		log:               log.NewHelper(log.With(logger, "module", "informer/node")),
		informerFactories: make(map[string]informers.SharedInformerFactory),
	}
}

func (x *NodeInformer) AddClusterInformer(ctx context.Context, clusterName string) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	// 检查 informer 是否已存在
	if _, exists := x.informerFactories[clusterName]; exists {
		return nil
	}

	// 创建 clientSet
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, clusterName)
	if err != nil {
		return fmt.Errorf("创建 clientSet 失败: 集群[%s]: %v", clusterName, err)
	}

	// 创建 informer factory
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Minute)
	nodeInformer := informerFactory.Core().V1().Nodes().Informer()

	// 添加事件处理
	_, err = nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*corev1.Node)
			x.log.WithContext(ctx).Debugf("集群[%s]添加节点[%s]", clusterName, node.Name)

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldNode := oldObj.(*corev1.Node)
			newNode := newObj.(*corev1.Node)
			x.log.WithContext(ctx).Debugf("集群[%s]更新节点[%s]，旧节点: %v，新节点: %v", clusterName, newNode.Name, oldNode, newNode)
		},
		DeleteFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if !ok {
				if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
					node, ok = tombstone.Obj.(*corev1.Node)
					if !ok {
						x.log.WithContext(ctx).Errorf("无法解析已删除的对象为Node: %v", obj)
						return
					}
				} else {
					x.log.WithContext(ctx).Errorf("无法解析对象为Node: %v", obj)
					return
				}
			}
			x.log.WithContext(ctx).Debugf("集群[%s]删除节点[%s]", clusterName, node.Name)
		},
	})

	if err != nil {
		return fmt.Errorf("添加事件处理器失败: 集群[%s]: %v", clusterName, err)
	}

	// 保存 informer factory
	x.informerFactories[clusterName] = informerFactory
	x.log.WithContext(ctx).Infof("集群 [%s] node informer 添加成功", clusterName)
	return nil
}

func (x *NodeInformer) RemoveClusterInformer(clusterName string) {
	x.mu.Lock()
	defer x.mu.Unlock()

	if factory, exists := x.informerFactories[clusterName]; exists {
		factory.Shutdown()
		delete(x.informerFactories, clusterName)
		x.log.Infof("集群 [%s] node informer 移除成功", clusterName)
	}
}

func (x *NodeInformer) SetInformerFactory(ctx context.Context) ([]informers.SharedInformerFactory, error) {
	informerFactories := make([]informers.SharedInformerFactory, 0)
	clusters, err := x.cluster.FetchAllClusters(ctx)
	if err != nil {
		return informerFactories, fmt.Errorf("查询集群列表失败: %v", err)
	}
	if clusters == nil {
		return informerFactories, nil
	}
	for _, cluster := range clusters {
		if err := x.AddClusterInformer(ctx, cluster.Name); err != nil {
			x.log.WithContext(ctx).Errorf("添加集群 [%s] informer 失败: %v", cluster.Name, err)
		}
	}

	x.mu.RLock()
	defer x.mu.RUnlock()

	factories := make([]informers.SharedInformerFactory, 0, len(x.informerFactories))
	for _, factory := range x.informerFactories {
		factories = append(factories, factory)
	}

	return factories, nil
}
