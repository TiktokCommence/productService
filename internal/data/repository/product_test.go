package repository

import (
	"context"
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/TiktokCommence/productService/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func initDB() *ProductInfoRepository {
	db, err := gorm.Open(mysql.Open("root:12345678@tcp(127.0.0.1:13306)/testDB?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("connect mysql failed")
	}
	if err := db.AutoMigrate(&model.Product{}, &model.ProductCategory{}); err != nil {
		panic(err)
	}
	gdb := NewGdb(db)
	pdr := NewProductInfoRepository(gdb)
	return pdr
}

func TestCreateProductInfo(t *testing.T) {
	p := initDB()
	var err error
	err = p.CreateProductInfo(context.Background(), &model.ProductInfo{
		Pd: &model.Product{
			ID:          1,
			Name:        "甜甜圈",
			Description: "很好吃",
			PictureUrl:  "123.com",
			Price:       10,
			MerchantID:  1,
		},
		Categories: []string{"food", "dessert"},
	})
	if err != nil {
		t.Error(err)
	}
	err = p.CreateProductInfo(context.Background(), &model.ProductInfo{
		Pd: &model.Product{
			ID:          2,
			Name:        "冰淇淋",
			Description: "很好吃",
			PictureUrl:  "123.com",
			Price:       10,
			MerchantID:  1,
		},
		Categories: []string{"food", "ice food"},
	})
	if err != nil {
		t.Error(err)
	}
}
func TestProductInfoRepository_GetProductInfosByIDs(t *testing.T) {
	p := initDB()
	ids := []uint64{1, 2}
	pdis, err := p.GetProductInfosByIDs(context.Background(), ids)
	if err != nil {
		t.Error(err)
	}
	for _, pdi := range pdis {
		t.Log(*pdi.Pd)
		t.Log(pdi.Categories)
	}
}

func TestProductInfoRepository_GetTotalNum(t *testing.T) {
	p := initDB()
	res, err := p.GetTotalNum(context.Background(), biz.ListOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	str := "ice food"
	res, err = p.GetTotalNum(context.Background(), biz.ListOptions{Category: &str})
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
