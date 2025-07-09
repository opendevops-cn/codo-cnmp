package event

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewClusterEventManager)
