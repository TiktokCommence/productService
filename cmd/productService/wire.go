//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/TiktokCommence/productService/internal/conf"
	"github.com/TiktokCommence/productService/internal/data"
	"github.com/TiktokCommence/productService/internal/data/cache"
	"github.com/TiktokCommence/productService/internal/data/repository"
	"github.com/TiktokCommence/productService/internal/registry"
	"github.com/TiktokCommence/productService/internal/server"
	"github.com/TiktokCommence/productService/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.RegistryConf, *conf.Expiration, *conf.ListOptions, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		registry.ProviderSet,
		wire.Bind(new(service.ProductHandler), new(*biz.ProductBiz)),
		wire.Bind(new(biz.Transaction), new(*repository.Gdb)),
		wire.Bind(new(biz.ProductInfoRepository), new(*repository.ProductInfoRepository)),
		wire.Bind(new(biz.ProductInfoCache), new(*cache.ProductCache)),
		wire.Bind(new(biz.GenerateIDer), new(*cache.GenerateIDImplement)),
		newApp,
	))
}
