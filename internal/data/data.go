package data

import (
	"github.com/TiktokCommence/productService/internal/data/cache"
	"github.com/TiktokCommence/productService/internal/data/repository"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	cache.NewRedisClient,
	cache.NewProductCache,
	cache.NewGenerateIDImplement,
	repository.NewDB,
	repository.NewGdb,
	repository.NewProductInfoRepository,
)
