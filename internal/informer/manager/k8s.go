package manager

import (
	"codo-cnmp/internal/informer"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

// InformerManager 实现 ClusterEventHandler
type InformerManager struct {
	informerList informer.K8sInformerList
	log          *log.Helper
}

func NewK8sInformerManager(list informer.K8sInformerList, logger log.Logger) *InformerManager {
	return &InformerManager{
		informerList: list,
		log:          log.NewHelper(log.With(logger, "module", "manager/k8s")),
	}
}

func (m *InformerManager) Name() string {
	return "informer-manager"
}

func (m *InformerManager) OnClusterAdd(ctx context.Context, clusterName string) error {
	for _, inf := range m.informerList {
		if err := inf.Add(ctx, clusterName); err != nil {
			m.log.WithContext(ctx).Errorf("添加集群 [%s] K8s informer 失败: %v", clusterName, err)
			return fmt.Errorf("添加集群 K8s informer 失败: %v", err)
		}
	}
	m.log.WithContext(ctx).Infof("集群 [%s] informer 添加成功", clusterName)
	return nil
}

func (m *InformerManager) OnClusterDelete(ctx context.Context, clusterName string) error {
	for _, inf := range m.informerList {
		inf.Remove(clusterName)
	}
	m.log.WithContext(ctx).Infof("集群 [%s] informer 移除成功", clusterName)
	return nil
}

// MetricsInformerManager 实现 ClusterEventHandler
type MetricsInformerManager struct {
	informerList informer.MetricsInformerList
	log          *log.Helper
}

func NewMetricsInformerManager(list informer.MetricsInformerList, logger log.Logger) *MetricsInformerManager {
	return &MetricsInformerManager{
		informerList: list,
		log:          log.NewHelper(log.With(logger, "module", "manager/metrics")),
	}
}

func (m *MetricsInformerManager) Name() string {
	return "metrics-informer-manager"
}

func (m *MetricsInformerManager) OnClusterAdd(ctx context.Context, clusterName string) error {
	for _, inf := range m.informerList {
		if err := inf.Add(ctx, clusterName); err != nil {
			m.log.WithContext(ctx).Errorf("添加集群 [%s] metrics informer 失败: %v", clusterName, err)
			return fmt.Errorf("添加集群 metrics informer 失败: %v", err)
		}
	}
	m.log.WithContext(ctx).Infof("集群 [%s] metrics informer 添加成功", clusterName)
	return nil
}

func (m *MetricsInformerManager) OnClusterDelete(ctx context.Context, clusterName string) error {
	for _, inf := range m.informerList {
		inf.Remove(clusterName)
	}
	m.log.WithContext(ctx).Infof("集群 [%s] metrics informer 移除成功", clusterName)
	return nil
}
