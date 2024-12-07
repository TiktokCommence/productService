package repository

import (
	"context"
	"fmt"
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/TiktokCommence/productService/internal/service"
	"strings"
)

type ProductInfoRepository struct {
	DB *Gdb
}

func NewProductInfoRepository(db *Gdb) *ProductInfoRepository {
	return &ProductInfoRepository{
		DB: db,
	}
}

func (p *ProductInfoRepository) CreateProductInfo(ctx context.Context, pi *model.ProductInfo) error {
	db := p.DB.DB(ctx)

	err := db.Create(pi.Pd).Error
	if err != nil {
		return err
	}
	categories := transformStringsToCategories(pi.Pd.ID, pi.Categories)
	err = db.Create(&categories).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductInfoRepository) UpdateProductInfo(ctx context.Context, pi *model.ProductInfo) error {
	var err error
	db := p.DB.DB(ctx)
	updateCategories, ok := ctx.Value(service.UpdateInfoKey).(bool)
	if !ok {
		return fmt.Errorf("ctx value not found")
	}
	if updateCategories {
		err = db.Table(model.ProductCategoryTableName).Where("p_id = ?", pi.Pd.ID).Delete(&model.ProductCategory{}).Error
		if err != nil {
			return err
		}
		categories := transformStringsToCategories(pi.Pd.ID, pi.Categories)
		err = db.Create(&categories).Error
		if err != nil {
			return err
		}
	}
	err = db.Updates(pi.Pd).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductInfoRepository) GetProductInfoByID(ctx context.Context, ID uint64) (*model.ProductInfo, error) {
	db := p.DB.db.WithContext(ctx)
	var pd model.Product
	err := db.Table(model.ProductTableName).Where("id = ?", ID).First(&pd).Error
	if err != nil {
		return nil, err
	}
	var categories []string
	err = db.Table(model.ProductCategoryTableName).Where("p_id = ?", ID).Pluck("category", &categories).Error
	if err != nil {
		return nil, err
	}
	return &model.ProductInfo{Pd: &pd, Categories: categories}, nil
}

func (p *ProductInfoRepository) DeleteProductInfo(ctx context.Context, ID uint64) error {
	db := p.DB.DB(ctx)
	err := db.Delete(&model.Product{}, ID).Error
	if err != nil {
		return err
	}
	err = db.Delete(&model.ProductCategory{}, "p_id =?", ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductInfoRepository) ListProductIDs(ctx context.Context, currentPage uint32, pageSize uint32, options biz.ListOptions) ([]uint64, error) {
	query := p.DB.db.WithContext(ctx).Table(model.ProductTableName)
	if options.Category != nil {
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.p_id", model.ProductCategoryTableName, model.ProductTableName, model.ProductCategoryTableName)).
			Where(fmt.Sprintf("%s.category =?", model.ProductCategoryTableName), *options.Category)
	}
	var ids []uint64
	err := query.Offset(int((currentPage-1)*pageSize)).Limit(int(pageSize)).Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (p *ProductInfoRepository) GetTotalNum(ctx context.Context, options biz.ListOptions) (uint32, error) {
	query := p.DB.db.WithContext(ctx).Table(model.ProductTableName)
	if options.Category != nil {
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.p_id", model.ProductCategoryTableName, model.ProductTableName, model.ProductCategoryTableName)).
			Where(fmt.Sprintf("%s.category = ?", model.ProductCategoryTableName), *options.Category)
	}
	var cnt int64
	err := query.Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return uint32(cnt), nil
}

func (p *ProductInfoRepository) GetProductInfosByIDs(ctx context.Context, ids []uint64) ([]*model.ProductInfo, error) {
	type Res struct {
		ID          uint64  `gorm:"column:id;primaryKey"`
		Name        string  `gorm:"column:name"`
		Description string  `gorm:"column:description"`
		PictureUrl  string  `gorm:"column:picture_url"`
		Price       float64 `gorm:"column:price"`
		MerchantID  uint64  `gorm:"column:merchant_id"`
		Categories  string  `gorm:"column:categories"`
	}
	results := make([]Res, len(ids))
	db := p.DB.db.WithContext(ctx)
	sql := fmt.Sprintf(`
		SELECT 
		    p.id,
		    p.name,
		    p.description,
		    p.picture_url,
		    p.price,
		    p.merchant_id,
		    GROUP_CONCAT(pc.category) AS categories
		FROM 
		    %s p
		LEFT JOIN 
		    %s pc ON p.id = pc.p_id
		WHERE 
		    p.id IN ?
		GROUP BY 
		    p.id
	`, model.ProductTableName, model.ProductCategoryTableName)
	err := db.Raw(sql, ids).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	pdis := make([]*model.ProductInfo, len(ids))
	for k, v := range results {
		pdis[k] = &model.ProductInfo{
			Pd: &model.Product{
				ID:          v.ID,
				Name:        v.Name,
				Description: v.Description,
				PictureUrl:  v.PictureUrl,
				Price:       v.Price,
				MerchantID:  v.MerchantID,
			},
			Categories: strings.Split(v.Categories, ","),
		}
	}
	return pdis, nil
}

func transformStringsToCategories(pid uint64, categoriesStr []string) []model.ProductCategory {
	categories := make([]model.ProductCategory, len(categoriesStr))
	for k, v := range categoriesStr {
		categories[k].Pid = pid
		categories[k].Category = v
	}
	return categories
}
