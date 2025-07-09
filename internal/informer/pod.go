package informer

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"sync"
	"time"

	"codo-cnmp/internal/biz"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type PodInformer struct {
	cluster       *biz.ClusterUseCase
	pod           *biz.PodUseCase
	gs            *biz.GameServerUseCase
	log           *log.Helper
	mu            sync.RWMutex
	informers     map[string]cache.SharedInformer
	stopChans     map[string]chan struct{}
	started       bool
	runningStatus map[string]bool
	rootCtx       context.Context
	cancel        context.CancelFunc
	eg            *errgroup.Group
}

func (x *PodInformer) Name() string {
	return "PodInformer"
}

func NewPodInformer(cluster *biz.ClusterUseCase, pod *biz.PodUseCase, gs *biz.GameServerUseCase, logger log.Logger) *PodInformer {
	// 创建 rootCtx
	rootCtx, cancel := context.WithCancel(context.Background())
	eg := errgroup.WithContext(rootCtx)
	return &PodInformer{
		cluster:       cluster,
		pod:           pod,
		gs:            gs,
		log:           log.NewHelper(log.With(logger, "module", "informer/pod")),
		informers:     make(map[string]cache.SharedInformer),
		stopChans:     make(map[string]chan struct{}),
		runningStatus: make(map[string]bool),
		rootCtx:       rootCtx,
		cancel:        cancel,
		eg:            eg,
	}
}

// createEventHandlers 创建统一的事件处理器
func (x *PodInformer) createEventHandlers(ctx context.Context, clusterName string) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			x.log.WithContext(ctx).Debugf("集群[%s]新增pod: %s", clusterName, pod.Name)
			// 随机sleep一下，避免创建过快
			time.Sleep(time.Duration(1000+rand.Intn(100)) * time.Millisecond)
		},
		//UpdateFunc: func(oldObj, newObj interface{}) {
		//	oldPod := oldObj.(*corev1.Pod)
		//	newPod := newObj.(*corev1.Pod)
		//	x.log.WithContext(ctx).Debugf("集群[%s]更新pod: %s -> %s",
		//		clusterName, oldPod.Name, newPod.Name)
		//	// TODO: 实现更新逻辑
		//},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
					pod, ok = tombstone.Obj.(*corev1.Pod)
					if !ok {
						x.log.WithContext(ctx).Errorf("无法解析已删除的对象为Pod: %v", obj)
						return
					}
				} else {
					x.log.WithContext(ctx).Errorf("无法解析对象为Pod: %v", obj)
					return
				}
			}
			x.log.WithContext(ctx).Debugf("集群[%s]删除pod: %s", clusterName, pod.Name)
			// 随机sleep一下，避免创建过快
			time.Sleep(time.Duration(1000+rand.Intn(100)) * time.Millisecond)
		},
	}
}

// Add 为指定集群添加 informer
func (x *PodInformer) Add(ctx context.Context, clusterName string) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	// 检查是否已存在
	if _, exists := x.informers[clusterName]; exists {
		return nil
	}

	// 获取 clientSet
	clientSet, err := x.cluster.GetClientSetByClusterName(x.rootCtx, clusterName)
	if err != nil {
		return fmt.Errorf("创建clientSet失败: %v", err)
	}

	// 创建 factory
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientSet,
		time.Minute,
		informers.WithNamespace(corev1.NamespaceAll),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.LabelSelector = ""
		}),
	)

	// 获取 pod informer
	informer := factory.Core().V1().Pods().Informer()

	// 添加事件处理器
	_, err = informer.AddEventHandler(x.createEventHandlers(x.rootCtx, clusterName))
	if err != nil {
		return fmt.Errorf("添加事件处理器失败: %v", err)
	}

	x.informers[clusterName] = informer
	x.log.WithContext(x.rootCtx).Infof("集群 [%s] pod informer 添加成功", clusterName)
	// 释放锁 并启动 informer
	// 如果服务已启动，则启动新添加的 informer
	if x.started && !x.runningStatus[clusterName] {
		if err := x.startInformerLocked(ctx, clusterName); err != nil {
			delete(x.informers, clusterName)
			return fmt.Errorf("启动 informer 失败: %v", err)
		}
	}

	return nil
}

// startInformerLocked 在已获得锁的情况下启动单个 informer
// 这不是线程安全的函数， 在调用之前， 需要确保已经 mu 加锁
func (x *PodInformer) startInformerLocked(ctx context.Context, clusterName string) error {
	// 检查 informer 是否已经在运行
	if x.runningStatus[clusterName] {
		x.log.WithContext(ctx).Infof("集群[%s] pod informer运行中", clusterName)
		return nil
	}

	informer, exists := x.informers[clusterName]
	if !exists {
		return fmt.Errorf("informer not found for cluster: %s", clusterName)
	}

	stopCh := make(chan struct{})
	x.stopChans[clusterName] = stopCh

	x.eg.Go(func(ctx context.Context) error {
		informer.Run(stopCh)
		return nil
	})

	// 等待缓存同步
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		x.log.WithContext(ctx).Errorf("集群[%s] pod informer 同步缓存超时", clusterName)
		close(stopCh)
		delete(x.stopChans, clusterName)
		x.runningStatus[clusterName] = false
		//return fmt.Errorf("集群[%s] pod informer 同步缓存超时", clusterName)
		return nil
	}

	x.runningStatus[clusterName] = true
	x.log.WithContext(ctx).Infof("集群[%s] pod informer 启动成功", clusterName)
	return nil
}

// startInformer 启动单个 informer（带锁保护）
func (x *PodInformer) startInformer(ctx context.Context, clusterName string) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.startInformerLocked(ctx, clusterName)
}

// Start 启动所有初始 informer
func (x *PodInformer) Start(ctx context.Context) error {
	x.mu.Lock()
	if x.started {
		x.mu.Unlock()
		return fmt.Errorf("pod informer already started")
	} else {
		x.started = true
		x.mu.Unlock()
	}

	// 获取集群列表
	clusters, err := x.cluster.FetchAllClusters(ctx)
	if err != nil {
		return fmt.Errorf("获取集群列表失败: %v", err)
	}

	// 为每个集群创建并启动 informer
	for _, cluster := range clusters {
		cluster := cluster // 创建副本用于闭包
		x.eg.Go(func(ctx context.Context) error {
			if err := x.Add(ctx, cluster.Name); err != nil {
				x.log.WithContext(ctx).Errorf("添加集群 [%s] informer 失败: %v",
					cluster.Name, err)
				return err
			}
			return nil
		})
	}

	return x.eg.Wait()
}

// Stop 停止所有 informer
func (x *PodInformer) Stop(ctx context.Context) error {
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
		x.log.WithContext(ctx).Infof("集群[%s] pod informer 停止成功", clusterName)
	}

	x.informers = make(map[string]cache.SharedInformer)
	x.stopChans = make(map[string]chan struct{})
	if x.cancel != nil {
		x.cancel()
	}
	return x.eg.Wait()
}

// Remove 移除指定集群的 informer
func (x *PodInformer) Remove(clusterName string) {
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
	x.log.Infof("集群 [%s] pod informer 移除成功", clusterName)
}
