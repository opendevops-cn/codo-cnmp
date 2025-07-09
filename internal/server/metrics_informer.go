package server

import (
	"codo-cnmp/internal/informer"
	"context"
	"fmt"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

type MetricsInformerServerWrapper struct {
	MetricsInformerList informer.MetricsInformerList
	cancel              func()
	eg                  *errgroup.Group
	log                 *log.Helper
	mu                  sync.Mutex
	started             bool
	processCtx          context.Context
}

func (x *MetricsInformerServerWrapper) Name() string {
	return "metrics-informer-server"
}

func NewMetricsInformerServer(MetricsInformerList informer.MetricsInformerList, logger log.Logger) *MetricsInformerServerWrapper {
	return &MetricsInformerServerWrapper{
		MetricsInformerList: MetricsInformerList,
		log:                 log.NewHelper(log.With(logger, "module", "server/metrics_informer")),
	}
}

func (x *MetricsInformerServerWrapper) Start(ctx context.Context) error {
	if !x.setStarted(true) {
		return fmt.Errorf("informer server 已启动")
	}
	x.processCtx, x.cancel = context.WithCancel(ctx)
	x.eg = errgroup.WithCancel(x.processCtx)

	for _, informerObj := range x.MetricsInformerList {
		if err := informerObj.Start(ctx); err != nil {
			x.log.WithContext(ctx).Errorf("启动 metrics informer 失败: %v", err)
			return err
		}
	}

	if err := x.eg.Wait(); err != nil {
		x.log.WithContext(ctx).Errorf("启动 metrics informer 失败: %v", err)
		return fmt.Errorf("启动 metrics informer 失败: %v", err)
	}
	return nil
}

func (x *MetricsInformerServerWrapper) Stop(ctx context.Context) error {
	if !x.setStarted(false) {
		return nil
	}
	if x.cancel != nil {
		x.cancel()
	}

	for _, informerObj := range x.MetricsInformerList {
		if err := informerObj.Stop(ctx); err != nil {
			x.log.WithContext(ctx).Errorf("停止 metrics informer 失败: %v", err)
		}
	}

	return x.eg.Wait()
}

func (x *MetricsInformerServerWrapper) setStarted(start bool) bool {
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
