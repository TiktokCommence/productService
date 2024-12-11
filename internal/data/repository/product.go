package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/TiktokCommence/productService/internal/errcode"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/TiktokCommence/productService/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"strings"
)

var _ biz.ProductInfoRepository = (*ProductInfoRepository)(nil)

type ProductInfoRepository struct {
	DB  *Gdb
	log *log.Helper
}

func NewProductInfoRepository(db *Gdb, logger log.Logger) *ProductInfoRepository {
	return &ProductInfoRepository{
		DB:  db,
		log: log.NewHelper(logger),
	}
}

func (p *ProductInfoRepository) CreateProductInfo(ctx context.Context, pi *model.ProductInfo) error {
	db := p.DB.DB(ctx)

	err := db.Debug().Create(pi.Pd).Error
	if err != nil {
		p.log.Errorf("mysql create %+v error", pi.Pd)
		return err
	}
	categories := transformStringsToCategories(pi.Pd.ID, pi.Categories)
	err = db.Debug().Create(&categories).Error
	if err != nil {
		p.log.Errorf("mysql create %+v error", categories)
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
		err = db.Table(model.ProductCategoryTableName).Where("p_id = ?", pi.Pd.ID).Debug().Delete(&model.ProductCategory{}).Error
		if err != nil {
			p.log.Errorf("mysql delete category p_id = %d", pi.Pd.ID)
			return err
		}
		categories := transformStringsToCategories(pi.Pd.ID, pi.Categories)
		err = db.Debug().Create(&categories).Error
		if err != nil {
			p.log.Errorf("mysql create category %+v error", categories)
			return err
		}
	}
	err = db.Debug().Updates(pi.Pd).Error
	if err != nil {
		p.log.Errorf("mysql update %+v", pi.Pd)
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.ErrProductNotExist
	}
	if err != nil {
		return nil, err
	}
	return &model.ProductInfo{Pd: &pd, Categories: categories}, nil
}

func (p *ProductInfoRepository) DeleteProductInfo(ctx context.Context, ID uint64) error {
	var err error
	db := p.DB.DB(ctx)
	res1 := db.Debug().Delete(&model.Product{}, "id = ?", ID)
	if res1.Error != nil || res1.RowsAffected == 0 {
		err = res1.Error
		if err == nil {
			err = errors.New("mysql don't exist the data")
		}
		p.log.Errorf("mysql delete product id=%v,reason:%v", ID, err)
		return err
	}
	res2 := db.Debug().Delete(&model.ProductCategory{}, "p_id = ?", ID)
	if res2.Error != nil || res2.RowsAffected == 0 {
		err = res2.Error
		if err == nil {
			err = errors.New("mysql don't exist the data")
		}
		p.log.Errorf("mysql delete product category p_id=%v,reason:%v", ID, err)
		return err
	}
	return nil
}

func (p *ProductInfoRepository) ListProductIDs(ctx context.Context, currentPage uint32, options service.ListOptions) ([]uint64, error) {
	query := p.DB.db.WithContext(ctx).Table(model.ProductTableName)
	if options.Category != nil {
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.p_id", model.ProductCategoryTableName, model.ProductTableName, model.ProductCategoryTableName)).
			Where(fmt.Sprintf("%s.category =?", model.ProductCategoryTableName), *options.Category)
	}
	var ids []uint64
	err := query.Offset(int((currentPage-1)*options.PageSize)).Limit(int(options.PageSize)).Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (p *ProductInfoRepository) GetTotalNum(ctx context.Context, options service.ListOptions) (uint32, error) {
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
