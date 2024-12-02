package service

import (
	"context"

	pb "github.com/TiktokCommence/productService/api/product/v1"
)

type ProductService struct {
	pb.UnimplementedProductServer
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductReq) (*pb.CreateProductResp, error) {
	return &pb.CreateProductResp{}, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductReq) (*pb.UpdateProductResp, error) {
	return &pb.UpdateProductResp{}, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.DeleteProductResp, error) {
	return &pb.DeleteProductResp{}, nil
}
func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	return &pb.ListProductsResp{}, nil
}
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.GetProductResp, error) {
	return &pb.GetProductResp{}, nil
}
func (s *ProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	return &pb.SearchProductsResp{}, nil
}
