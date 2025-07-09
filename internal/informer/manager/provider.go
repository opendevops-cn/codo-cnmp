package manager

import (
	"codo-cnmp/internal/event"
	"codo-cnmp/internal/server"
)

func NewManagers(informerManager *InformerManager, metricsInformerManager *MetricsInformerManager,
	metricsServer *server.MetricsInformerServerWrapper, k8sServer *server.InformerServerWrapper) []event.ClusterEventHandler {
	return []event.ClusterEventHandler{
		informerManager,
		metricsInformerManager,
		//metricsServer,
		//k8sServer,
	}
}
