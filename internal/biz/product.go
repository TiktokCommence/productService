package biz

import (
	"context"
	"fmt"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/TiktokCommence/productService/internal/service"
	"math"
)

var _ service.ProductHandler = (*ProductBiz)(nil)

type ProductBiz struct {
	tx Transaction
	pr ProductInfoRepository
	g  GenerateIDer
}

func (p *ProductBiz) CreateProduct(ctx context.Context, productInfo *model.ProductInfo) error {
	var err error
	productInfo.Pd.ID, err = p.g.GenerateID()
	if err != nil {
		return fmt.Errorf("%v,reason:%w", "generate ID failed", err)
	}
	err = p.tx.InTx(ctx, func(ctx context.Context) error {
		err = p.pr.CreateProductInfo(ctx, productInfo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("create product error:%w", err)
	}
	return nil
}

func (p *ProductBiz) UpdateProduct(ctx context.Context, productInfo *model.ProductInfo) error {
	var err error
	err = p.tx.InTx(ctx, func(ctx context.Context) error {
		// 注意更新时可根据ctx中的service.UpdateInfoKey来判断是否需要更新category
		err = p.pr.UpdateProductInfo(ctx, productInfo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("update product error:%w", err)
	}
	return nil
}

func (p *ProductBiz) GetProductInfoByID(ctx context.Context, ID uint64) (*model.ProductInfo, error) {
	pdi, err := p.pr.GetProductInfoByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get product info failed:%w", err)
	}
	return pdi, nil
}

func (p *ProductBiz) DeleteProduct(ctx context.Context, ID uint64) error {
	var err error
	err = p.tx.InTx(ctx, func(ctx context.Context) error {
		err = p.pr.DeleteProductInfo(ctx, ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete product error:%w", err)
	}
	return nil
}

func (p *ProductBiz) ListProducts(ctx context.Context, page uint32, pageSize uint32, category *string, totalPage *uint32) ([]*model.ProductInfo, error) {
	ListOpts := NewListOptions(WithCategory(category))
	res, err := p.pr.GetTotalNum(ctx, ListOpts)
	if err != nil {
		return nil, fmt.Errorf("get total num error: %w", err)
	}
	*totalPage = uint32(math.Ceil(float64(res) / float64(pageSize)))

	ids, err := p.pr.ListProductIDs(ctx, page, pageSize, ListOpts)
	if err != nil {
		return nil, fmt.Errorf("list product ids error: %w", err)
	}
	pdi, err := p.pr.GetProductInfosByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get product info by ids error:%w", err)
	}
	return pdi, nil
}
