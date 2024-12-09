package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/TiktokCommence/productService/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math"
	"time"
)

var _ service.ProductHandler = (*ProductBiz)(nil)

type ProductBiz struct {
	tx    Transaction
	pr    ProductInfoRepository
	cache ProductInfoCache
	g     GenerateIDer
	log   *log.Helper
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
	//采用延迟双删来解决数据一致性问题
	var err error
	err = p.cache.DeleteProductInfo(ctx, productInfo.Pd.ID)
	if err != nil {
		return fmt.Errorf("%v,reason:%w", "delete old product info failed", err)
	}

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
	go func() {
		//睡眠1秒钟，再删除缓存
		time.Sleep(time.Second)

		err1 := p.cache.DeleteProductInfo(context.Background(), productInfo.Pd.ID)
		if err1 != nil {
			p.log.Warn("delete old product info cache failed ", err)
		}
	}()

	return nil
}

func (p *ProductBiz) GetProductInfoByID(ctx context.Context, ID uint64) (*model.ProductInfo, error) {
	pdi, err := p.cache.GetProductInfo(ctx, ID)
	//如果返回的pdi是nil，说明其存入的时候就是个空值，也就是不存在的数据，直接返回
	//（为了应对缓存穿透，即大量请求去访问一个不存在的数据，增加数据库的压力）
	//如果获取到了数据，当然就直接返回
	if pdi == nil || err == nil {
		return pdi, nil
	}
	if !errors.Is(err, redis.Nil) {
		p.log.Warn("get product info cache failed ", err)
	}
	//如果redis不存在该key，就从数据库中查
	pdi, err = p.pr.GetProductInfoByID(ctx, ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get product info failed:%w", err)
	}
	go func() {
		var err1 error
		if err == nil {
			//err==nil说明数据库中有数据
			//TODO:将expire的配置传入，完成expire的离散化，防止缓存雪崩
			err1 = p.cache.SetProductInfo(context.Background(), ID, pdi, 60*60*24*7)
		} else {
			//err!=nil说明数据库中没有数据
			//TODO: 设定一个 expire，让 cache 5分钟后过期，防止大量的 cache 占用
			// 并且在这个 expire 期内，如果有新的数据更新，会重新 set 这个 cache
			// 保证了 cache 的高可用性和快速的过期
			// 但是在大量的并发场景下，需要保证 cache 过期后的数据是最新的
			err1 = p.cache.SetProductInfo(context.Background(), ID, nil, 60*5)
		}
		if err != nil {
			p.log.Warn("set product info cache failed ", err1)
		}
	}()

	return pdi, nil
}

func (p *ProductBiz) DeleteProduct(ctx context.Context, ID uint64) error {
	var err error
	err = p.cache.DeleteProductInfo(ctx, ID)
	if err != nil {
		p.log.Warn("delete product info cache failed ", err)
	}
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
	go func() {
		time.Sleep(time.Second)

		err1 := p.cache.DeleteProductInfo(context.Background(), ID)
		if err1 != nil {
			p.log.Warn("delay delete old product info cache failed ", err)
		}
	}()
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
	// 从缓存中获取产品信息，加速查询，缓存中不一定有所有id的信息，所以要返回缓存中没有的id的信息
	unfoundIDs, pdis, err := p.cache.MgetProductInfo(ctx, ids)
	if err == nil && len(unfoundIDs) == 0 {
		return pdis, nil
	}

	dbpdis, err := p.pr.GetProductInfosByIDs(ctx, unfoundIDs)
	if err != nil {
		return nil, fmt.Errorf("get product info by ids error:%w", err)
	}
	pdis = append(pdis, dbpdis...)
	go func() {
		mp := make(map[uint64]*model.ProductInfo, len(dbpdis))
		for _, pdi := range dbpdis {
			mp[pdi.Pd.ID] = pdi
		}
		err = p.cache.MsetProductInfo(context.Background(), mp)
		if err != nil {
			p.log.Warn("set product info cache failed ", err)
		}
	}()
	return pdis, nil
}
