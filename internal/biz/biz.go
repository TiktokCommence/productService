package biz

import (
	"context"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/TiktokCommence/productService/internal/service"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewProductBiz)

type Transaction interface {
	// DB(ctx context.Context) *gorm.DB
	InTx(ctx context.Context, fc func(ctx context.Context) error) error
}
type ProductInfoRepository interface {
	CreateProductInfo(ctx context.Context, pi *model.ProductInfo) error
	UpdateProductInfo(ctx context.Context, pi *model.ProductInfo) error
	GetProductInfoByID(ctx context.Context, ID uint64) (*model.ProductInfo, error)
	DeleteProductInfo(ctx context.Context, ID uint64) error
	ListProductIDs(ctx context.Context, currentPage uint32, options service.ListOptions) ([]uint64, error)

	GetTotalNum(ctx context.Context, options service.ListOptions) (uint32, error)
	GetProductInfosByIDs(ctx context.Context, ids []uint64) ([]*model.ProductInfo, error)
}
type ProductInfoCache interface {
	SetProductInfo(ctx context.Context, id uint64, pi *model.ProductInfo, expire int) error
	GetProductInfo(ctx context.Context, id uint64) (*model.ProductInfo, error)
	DeleteProductInfo(ctx context.Context, id uint64) error
	MgetProductInfo(ctx context.Context, ids []uint64) ([]uint64, []*model.ProductInfo, error)
	MsetProductInfo(ctx context.Context, mp map[uint64]*model.ProductInfo) error
}

type GenerateIDer interface {
	GenerateID() (uint64, error)
}
