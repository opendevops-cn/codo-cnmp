package event

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// ClusterEventHandler 定义集群事件处理接口
type ClusterEventHandler interface {
	Name() string
	OnClusterAdd(ctx context.Context, clusterName string) error
	OnClusterDelete(ctx context.Context, clusterName string) error
}

// RegisteredHandler 已注册的处理器
type RegisteredHandler struct {
	Name    string
	Handler ClusterEventHandler
}

// ClusterEventManager 管理集群事件
type ClusterEventManager struct {
	handlers []RegisteredHandler
	log      *log.Helper
}

// NewClusterEventManager 创建事件管理器
func NewClusterEventManager(logger log.Logger, handlers []ClusterEventHandler) *ClusterEventManager {
	x := &ClusterEventManager{
		handlers: make([]RegisteredHandler, 0),
		log:      log.NewHelper(log.With(logger, "module", "event/cluster")),
	}
	for _, h := range handlers {
		x.RegisterHandler(h.Name(), h)
	}

	return x
}

// RegisterHandler 注册事件处理器
func (m *ClusterEventManager) RegisterHandler(name string, handler ClusterEventHandler) {
	m.handlers = append(m.handlers, RegisteredHandler{
		Name:    name,
		Handler: handler,
	})
}

// OnClusterAdd 处理集群添加事件
func (m *ClusterEventManager) OnClusterAdd(ctx context.Context, clusterName string) error {
	for _, h := range m.handlers {
		if err := h.Handler.OnClusterAdd(ctx, clusterName); err != nil {
			m.log.WithContext(ctx).Errorf("eventManager 处理器 [%s] 处理集群 [%s] 添加事件失败: %v", h.Name, clusterName, err)
			return err
		}
	}
	return nil
}

// OnClusterDelete 处理集群删除事件
func (m *ClusterEventManager) OnClusterDelete(ctx context.Context, clusterName string) error {
	for _, h := range m.handlers {
		if err := h.Handler.OnClusterDelete(ctx, clusterName); err != nil {
			m.log.WithContext(ctx).Errorf("enventManager 处理器 [%s] 处理集群 [%s] 删除事件失败: %v", h.Name, clusterName, err)
			return err
		}
	}
	return nil
}
