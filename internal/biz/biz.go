package biz

import (
	"context"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet()

type Transaction interface {
	// DB(ctx context.Context) *gorm.DB
	InTx(ctx context.Context, fc func(ctx context.Context) error) error
}
type ProductInfoRepository interface {
	CreateProductInfo(ctx context.Context, pi *model.ProductInfo) error
	UpdateProductInfo(ctx context.Context, pi *model.ProductInfo) error
	GetProductInfoByID(ctx context.Context, ID uint64) (*model.ProductInfo, error)
	DeleteProductInfo(ctx context.Context, ID uint64) error
	ListProductIDs(ctx context.Context, currentPage uint32, pageSize uint32, options ListOptions) ([]uint64, error)

	GetTotalNum(ctx context.Context, options ListOptions) (uint32, error)
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

type ListOptions struct {
	Category *string
}
type ListOption func(options *ListOptions)

func NewListOptions(opt ...ListOption) ListOptions {
	var defaultListOptions = ListOptions{}
	for _, o := range opt {
		o(&defaultListOptions)
	}
	return defaultListOptions
}

func WithCategory(category *string) ListOption {
	return func(options *ListOptions) {
		options.Category = category
	}
}
