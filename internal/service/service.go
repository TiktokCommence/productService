package service

import (
	"context"
	"errors"
	pb "github.com/TiktokCommence/productService/api/product/v1"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet()

const (
	UpdateInfoKey = "update_info"
)

type ProductHandler interface {
	CreateProduct(ctx context.Context, productInfo *model.ProductInfo) error
	UpdateProduct(ctx context.Context, productInfo *model.ProductInfo) error
	GetProductInfoByID(ctx context.Context, ID uint64) (*model.ProductInfo, error)
	DeleteProduct(ctx context.Context, ID uint64) error

	ListProducts(ctx context.Context, page uint32, pageSize uint32, category *string) ([]*model.ProductInfo, error)
}

var (
	ErrCreateProduct = errors.New("create product failed")
	ErrUpdateProduct = errors.New("update product failed")

	ErrCategoryIsEmpty = errors.New("category is empty")
	ErrDeleteProduct   = errors.New("delete product failed")
	ErrListProduct     = errors.New("list products failed")
	ErrGetProduct      = errors.New("get product failed")
)

func transformProducts(pdis []*model.ProductInfo) []*pb.ProductInfo {
	products := make([]*pb.ProductInfo, 0)
	for _, pdi := range pdis {
		products = append(products, transformProduct(pdi))
	}
	return products
}
func transformProduct(pdi *model.ProductInfo) *pb.ProductInfo {
	return &pb.ProductInfo{
		Id:          pdi.Pd.ID,
		Name:        pdi.Pd.Name,
		Description: pdi.Pd.Description,
		Price:       pdi.Pd.Price,
		Picture:     pdi.Pd.PictureUrl,
		Merchant:    pdi.Pd.MerchantID,
		Categories:  pdi.Categories,
	}
}
