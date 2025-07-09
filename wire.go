//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"codo-cnmp/initialization"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/conf"
	"codo-cnmp/internal/data"
	"codo-cnmp/internal/dep"
	"codo-cnmp/internal/event"
	"codo-cnmp/internal/informer"
	"codo-cnmp/internal/informer/manager"
	"codo-cnmp/internal/job"
	"codo-cnmp/internal/middleware"
	"codo-cnmp/internal/server"
	"codo-cnmp/internal/service"
	"codo-cnmp/migrate"
	"context"
	"github.com/go-kratos/kratos/v2"

	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(context.Context, *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, dep.ProviderSet,
		job.ProviderSet, informer.ProviderSet, middleware.ProviderSet, migrate.ProviderSet, initialization.ProviderSet,
		event.ProviderSet, manager.ProviderSet, newApp))
}
