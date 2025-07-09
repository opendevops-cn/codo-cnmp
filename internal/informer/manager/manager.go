package manager

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewK8sInformerManager, NewMetricsInformerManager, NewManagers)
