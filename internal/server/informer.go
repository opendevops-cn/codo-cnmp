package server

import (
	"codo-cnmp/internal/informer"
	"context"
	"fmt"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

//var _ event.ClusterEventHandler = (*InformerServerWrapper)(nil)

type InformerServerWrapper struct {
	K8sInformerList informer.K8sInformerList
	cancel          func()
	eg              *errgroup.Group
	log             *log.Helper
	mu              sync.RWMutex
	started         bool
	processCtx      context.Context
}

func NewInformerServer(K8sInformerList informer.K8sInformerList, logger log.Logger) *InformerServerWrapper {
	return &InformerServerWrapper{
		K8sInformerList: K8sInformerList,
		log:             log.NewHelper(log.With(logger, "module", "server/informer")),
		started:         false,
	}
}

func (x *InformerServerWrapper) Name() string {
	return "informer-server"
}

func (x *InformerServerWrapper) Start(ctx context.Context) error {
	// 只锁 started 检查和修改
	if !x.setStarted(true) {
		return fmt.Errorf("informer server 已启动")
	}
	x.processCtx, x.cancel = context.WithCancel(ctx)
	x.eg = errgroup.WithCancel(ctx)

	// 启动 informer
	for _, myInformer := range x.K8sInformerList {
		inf := myInformer // 创建副本用于闭包
		x.eg.Go(func(ctx context.Context) error {
			return inf.Start(x.processCtx)
		})
	}
	if err := x.eg.Wait(); err != nil {
		x.log.WithContext(ctx).Errorf("informer启动失败: %v", err)
		return err
	}
	x.log.WithContext(ctx).Info("informer启动成功")
	return nil
}

func (x *InformerServerWrapper) Stop(ctx context.Context) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	if !x.setStarted(false) {
		return nil
	}
	if x.cancel != nil {
		x.cancel()
	}

	for _, myInformer := range x.K8sInformerList {
		inf := myInformer // 创建副本用于闭包
		x.log.WithContext(ctx).Info("stopping informer(%s) server", inf.Name())
		_ = inf.Stop(ctx)
		defer func() {
			x.log.Infof("informer(%s) server stopped", inf.Name())
		}()
	}
	x.log.WithContext(ctx).Info("informer停止成功")
	return x.eg.Wait()
}

func (x *InformerServerWrapper) setStarted(start bool) bool {
	x.mu.Lock()
	defer x.mu.Unlock()

	if start && x.started {
		return false
	}
	if !start && !x.started {
		return false
	}
	x.started = start
	return true
}
