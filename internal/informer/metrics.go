package informer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type MetricsInformer struct {
	cluster       *biz.ClusterUseCase
	node          *biz.NodeUseCase
	log           *log.Helper
	mu            sync.RWMutex
	informers     map[string]cache.SharedInformer
	stopChans     map[string]chan struct{}
	runningStatus map[string]bool
	started       bool
	rootCtx       context.Context
	cancelFunc    context.CancelFunc
	eg            *errgroup.Group
}

func (x *MetricsInformer) Name() string {
	return "MetricsInformer"
}

func NewMetricsInformer(cluster *biz.ClusterUseCase, node *biz.NodeUseCase, logger log.Logger) *MetricsInformer {
	// 创建 root context
	rootCtx, cancel := context.WithCancel(context.Background())
	eg := errgroup.WithContext(rootCtx)
	return &MetricsInformer{
		cluster:       cluster,
		node:          node,
		log:           log.NewHelper(log.With(logger, "module", "informer/metrics")),
		informers:     make(map[string]cache.SharedInformer),
		runningStatus: make(map[string]bool),
		started:       false,
		rootCtx:       rootCtx,
		cancelFunc:    cancel,
		eg:            eg,
		stopChans:     make(map[string]chan struct{}),
	}
}

// createEventHandlers 创建统一的事件处理器
func (x *MetricsInformer) createEventHandlers(ctx context.Context, clusterName string) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nodeMetrics := obj.(*v1beta1.NodeMetrics)
			x.log.WithContext(ctx).Debugf("集群[%s]新增node: %s", clusterName, nodeMetrics.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			x.handleMetricsUpdate(ctx, clusterName, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			nodeMetrics, ok := obj.(*v1beta1.NodeMetrics)
			if !ok {
				if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
					nodeMetrics, ok = tombstone.Obj.(*v1beta1.NodeMetrics)
					if !ok {
						x.log.WithContext(ctx).Errorf("无法解析已删除的对象为NodeMetrics: %v", obj)
						return
					}
				} else {
					x.log.WithContext(ctx).Errorf("无法解析对象为NodeMetrics: %v", obj)
					return
				}
			}
			x.log.WithContext(ctx).Debugf("集群[%s]删除nodeMetrics: %s", clusterName, nodeMetrics.Name)
		},
	}
}

// handleMetricsUpdate 统一的更新处理逻辑
func (x *MetricsInformer) handleMetricsUpdate(ctx context.Context, clusterName string, newObj interface{}) {
	newNodeMetrics := newObj.(*v1beta1.NodeMetrics)
	cpuUsage := newNodeMetrics.Usage.Cpu().MilliValue()
	memoryUsage := utils.ConvertMemoryToGiB(*newNodeMetrics.Usage.Memory())
	CpuUsage := float32(utils.ConvertCPUToCores(*resource.NewQuantity(cpuUsage, resource.DecimalSI)))

	cluster, err := x.cluster.GetClusterByName(ctx, clusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取集群[%s]信息失败: %v", clusterName, err)
		return
	}

	err = x.node.UpdateNode(ctx, &biz.NodeItem{
		ClusterID:   cluster.ID,
		Name:        newNodeMetrics.Name,
		CpuUsage:    CpuUsage,
		MemoryUsage: float32(memoryUsage),
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("集群[%s]更新节点[%s]node metrics 指标失败: %v",
			clusterName, newNodeMetrics.Name, err)
	} else {
		x.log.WithContext(ctx).Infof("集群[%s]更新节点[%s]node metrics 指标成功",
			clusterName, newNodeMetrics.Name)
	}
}

// Add 为指定集群添加 informer
func (x *MetricsInformer) Add(ctx context.Context, clusterName string) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	if _, exists := x.informers[clusterName]; exists {
		return nil
	}

	// 获取 metrics client
	metricsClient, err := x.cluster.GetMetricsClientSetByClusterName(x.rootCtx, clusterName)
	if err != nil {
		return fmt.Errorf("创建 metrics clientset 失败: 集群[%s]: %v", clusterName, err)
	}

	// 创建 ListWatcher
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return metricsClient.MetricsV1beta1().NodeMetricses().List(x.rootCtx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return metricsClient.MetricsV1beta1().NodeMetricses().Watch(x.rootCtx, options)
		},
	}

	// 创建 informer
	informer := cache.NewSharedInformer(lw, &v1beta1.NodeMetrics{}, time.Second*30)

	// 添加事件处理器
	_, err = informer.AddEventHandler(x.createEventHandlers(x.rootCtx, clusterName))
	if err != nil {
		return fmt.Errorf("添加event handler 失败: 集群[%s]: %v", clusterName, err)
	}

	x.informers[clusterName] = informer
	x.log.WithContext(x.rootCtx).Infof("集群 [%s] metrics informer 添加成功", clusterName)
	// 释放锁 启动 informer
	// 如果 informer 已经启动，则启动新添加的 informer
	if x.started && !x.runningStatus[clusterName] {
		if err := x.startInformerLocked(ctx, clusterName); err != nil {
			delete(x.informers, clusterName)
			return fmt.Errorf("启动 informer 失败: %v", err)
		}
	}
	return nil
}

// startInformer 启动单个 informer
// 这不是线程安全的函数， 在调用之前， 需要确保已经 mu 加锁
func (x *MetricsInformer) startInformerLocked(ctx context.Context, clusterName string) error {
	informer, exists := x.informers[clusterName]
	if !exists {
		return fmt.Errorf("集群 [%s] 未设置metrics  informer", clusterName)
	}

	if x.runningStatus[clusterName] {

		return fmt.Errorf("集群 [%s] metrics informer 已经启动", clusterName)
	}

	stopCh := make(chan struct{})
	x.stopChans[clusterName] = stopCh

	x.eg.Go(func(ctx context.Context) error {
		informer.Run(stopCh)
		return nil
	})

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		x.log.WithContext(ctx).Errorf("集群[%s] metrics 同步缓存超时", clusterName)
		close(stopCh)
		delete(x.stopChans, clusterName)
		x.runningStatus[clusterName] = false
		//return fmt.Errorf("集群[%s] metrics 同步缓存超时", clusterName)
		return nil
	}
	x.runningStatus[clusterName] = true
	x.log.WithContext(ctx).Infof("集群[%s] metrics informer 启动成功", clusterName)
	return nil
}

// Start 启动所有 informer
func (x *MetricsInformer) Start(ctx context.Context) error {
	x.mu.Lock()
	if x.started {
		x.mu.Unlock()
		return fmt.Errorf("metrics informer 已经启动")
	} else {
		x.started = true
		x.mu.Unlock()
	}

	// 获取所有集群
	clusters, err := x.cluster.FetchAllClusters(ctx)
	if err != nil {
		return fmt.Errorf("获取集群列表失败: %v", err)
	}

	// 为每个集群创建 informer
	for _, cluster := range clusters {
		cluster := cluster
		x.eg.Go(func(ctx context.Context) error {
			if err := x.Add(ctx, cluster.Name); err != nil {
				x.log.WithContext(ctx).Errorf("添加集群 [%s] informer 失败: %v", cluster.Name, err)
				return nil
			}
			return nil
		})
	}

	return x.eg.Wait()
}

// Stop 停止所有 informer
func (x *MetricsInformer) Stop(ctx context.Context) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	if !x.started {
		return nil
	}
	x.started = false

	for clusterName, stopCh := range x.stopChans {
		close(stopCh)
		delete(x.stopChans, clusterName)
		delete(x.runningStatus, clusterName)
		x.log.WithContext(ctx).Infof("集群[%s] metrics informer 停止成功", clusterName)
	}

	x.informers = make(map[string]cache.SharedInformer)
	x.stopChans = make(map[string]chan struct{})
	if x.cancelFunc != nil {
		x.cancelFunc()
	}
	return x.eg.Wait()
}

// Remove 移除指定集群的 informer
func (x *MetricsInformer) Remove(clusterName string) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if _, exists := x.informers[clusterName]; !exists {
		return
	}
	if stopCh, exists := x.stopChans[clusterName]; exists {
		close(stopCh)
		delete(x.stopChans, clusterName)
		delete(x.runningStatus, clusterName)
	}

	delete(x.informers, clusterName)
	x.log.Infof("集群 [%s] metrics informer 移除成功", clusterName)
}
