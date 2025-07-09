package dep

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewRedis, NewMysql, NewMeterProvider, NewTracerProvider, NewTextMapPropagator,
	NewLogger, NewMsgBus, NewTopicManager, NewSharedStorage, NewCAS, NewCODOAPIGateway, NewKafka)
