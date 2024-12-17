//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/bilibili/HCP/app/interface/v1/internal/data"
	"github.com/bilibili/HCP/app/interface/v1/internal/server"
	"github.com/bilibili/HCP/app/interface/v1/internal/service"
	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(data.ProviderSet, service.ProviderSet, server.NewHttpServer, NewApp))
}
