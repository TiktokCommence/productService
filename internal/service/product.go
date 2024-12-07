package service

import (
	"context"
	"fmt"
	pb "github.com/TiktokCommence/productService/api/product/v1"
	"github.com/TiktokCommence/productService/internal/conf"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/TiktokCommence/productService/internal/tool"
)

type ProductService struct {
	pb.UnimplementedProductServer
	phandler    ProductHandler
	pageOptions *conf.ListOptions
}

func NewProductService(phandler ProductHandler) *ProductService {
	return &ProductService{
		phandler: phandler,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductReq) (*pb.CreateProductResp, error) {
	if req.Categories == nil || len(req.Categories) == 0 {
		return &pb.CreateProductResp{}, fmt.Errorf("%w,reason:%w", ErrCreateProduct, ErrCategoryIsEmpty)
	}
	pd := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		PictureUrl:  req.Description,
		Price:       req.Price,
		MerchantID:  req.Merchant,
	}
	err := s.phandler.CreateProduct(ctx, &model.ProductInfo{
		Pd:         pd,
		Categories: req.Categories,
	})
	if err != nil {
		return &pb.CreateProductResp{}, fmt.Errorf("%w,reason:%w", ErrCreateProduct, err)
	}
	return &pb.CreateProductResp{
		Id: pd.ID,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductReq) (*pb.UpdateProductResp, error) {

	/*
		Notice:
			因为一般情况下，这个操作不会发生并发问题
			所以我们没有去针对并发冲突去做处理
	*/

	var updateCategory bool
	oldPdi, err := s.phandler.GetProductInfoByID(ctx, req.GetId())
	if err != nil || oldPdi == nil {
		return &pb.UpdateProductResp{
			Success: false,
		}, fmt.Errorf("%w,reason:%w", ErrUpdateProduct, err)
	}

	newPd := oldPdi.Pd
	if req.Name != nil {
		newPd.Name = req.GetName()
	}
	if req.Description != nil {
		newPd.Description = req.GetDescription()
	}
	if req.Picture != nil {
		newPd.PictureUrl = req.GetPicture()
	}
	if req.Price != nil {
		newPd.Price = req.GetPrice()
	}

	//由于更新类别字段涉及到另一个表，为了避免额外的IO，对类别字段是否需要更新做出判断
	//如果更新请求的类别不为空，则要更新类别
	if req.Categories != nil && !tool.CheckSliceEqual(req.Categories, oldPdi.Categories) {
		updateCategory = true
	}

	//将是否更新类别字段存入ctx中
	vctx := context.WithValue(ctx, UpdateInfoKey, updateCategory)
	//更新操作根据vctx中的值来判断是否需要更新类别
	err = s.phandler.UpdateProduct(vctx, &model.ProductInfo{
		Pd:         oldPdi.Pd,
		Categories: req.Categories,
	})
	if err != nil {
		return &pb.UpdateProductResp{
			Success: false,
		}, fmt.Errorf("%w,reason:%w", ErrUpdateProduct, err)
	}
	return &pb.UpdateProductResp{
		Success: true,
	}, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.DeleteProductResp, error) {
	err := s.phandler.DeleteProduct(ctx, req.GetId())
	if err != nil {
		return &pb.DeleteProductResp{
			Success: false,
		}, fmt.Errorf("%w,reason:%w", ErrDeleteProduct, err)
	}
	return &pb.DeleteProductResp{
		Success: true,
	}, nil
}
func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	var totalPage uint32
	pdis, err := s.phandler.ListProducts(ctx, req.GetPage(), s.pageOptions.Pagesize, req.CategoryName, &totalPage)
	if err != nil {
		return &pb.ListProductsResp{}, fmt.Errorf("%w,reason:%w", ErrListProduct, err)
	}
	return &pb.ListProductsResp{
		Products:    transformProducts(pdis),
		CurrentPage: req.GetPage(),
		TotalPages:  totalPage,
	}, nil
}
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.GetProductResp, error) {
	pdi, err := s.phandler.GetProductInfoByID(ctx, req.GetId())
	if err != nil {
		return &pb.GetProductResp{}, fmt.Errorf("%w,reason:%w", ErrGetProduct, err)
	}

	return &pb.GetProductResp{
		Product: transformProduct(pdi),
	}, nil
}
func (s *ProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	return &pb.SearchProductsResp{}, nil
}
