package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewPprofServer, NewPrometheusServer, NewCronServer, NewInformerServer,
	NewWebSocketServer, NewMetricsInformerServer, NewAPIServerProxy, NewAPIServerHandler)
