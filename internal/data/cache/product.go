package cache

import (
	"context"
	"fmt"
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/TiktokCommence/productService/internal/conf"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

var _ biz.ProductInfoCache = (*ProductCache)(nil)

type ProductCache struct {
	cli *redis.Client
	log *log.Helper
}

func NewProductCache(cli *redis.Client, logger log.Logger) *ProductCache {
	return &ProductCache{
		cli: cli,
		log: log.NewHelper(logger),
	}
}

func NewRedisClient(cf *conf.Data) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:         cf.Redis.Addr,
		Password:     cf.Redis.Password,
		MaxRetries:   int(cf.Redis.MaxRetry),
		ReadTimeout:  time.Duration(cf.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cf.Redis.WriteTimeout) * time.Second,
		PoolSize:     int(cf.Redis.PoolSize),
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return cli

}

func (p *ProductCache) SetProductInfo(ctx context.Context, id uint64, pi *model.ProductInfo, expire int) error {
	key := p.generateKey(id)
	expTime := p.generateExpiration(expire)
	val, err := pi.Write()
	if err != nil {
		return err
	}
	_, err = p.cli.Set(ctx, key, val, expTime).Result()
	if err != nil {
		p.log.Errorf("set product info cache failed,key:%v,value:%v,error:%v", key, *pi, err)
		return err
	}
	return nil
}

func (p *ProductCache) GetProductInfo(ctx context.Context, id uint64) (*model.ProductInfo, error) {
	key := p.generateKey(id)
	val, err := p.cli.Get(ctx, key).Result()
	if err != nil {
		p.log.Errorf("get product info cache failed,key:%s,error:%v", key, err)
		return nil, err
	}
	pi := &model.ProductInfo{}
	err = pi.Read(val)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (p *ProductCache) DeleteProductInfo(ctx context.Context, id uint64) error {
	key := p.generateKey(id)
	_, err := p.cli.Del(ctx, key).Result()
	if err != nil {
		p.log.Errorf("delete product info cache failed,key:%s,error:%v", key, err)
		return err
	}
	return nil
}

func (p *ProductCache) MgetProductInfo(ctx context.Context, ids []uint64) ([]uint64, []*model.ProductInfo, error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = p.generateKey(id)
	}
	values, err := p.cli.MGet(ctx, keys...).Result()
	if err != nil {
		p.log.Errorf("mget product info cache failed,keys:%v,error:%v", keys, err)
		return ids, nil, err
	}
	unfoundIDs := make([]uint64, 0, len(ids))
	res := make([]*model.ProductInfo, 0)
	for k, v := range values {
		if v == nil {
			unfoundIDs = append(unfoundIDs, ids[k])
			continue
		} else {
			pi := &model.ProductInfo{}
			err = pi.Read(v.(string))
			if err != nil {
				return ids, nil, err
			}
			res = append(res, pi)
		}
	}
	return unfoundIDs, res, nil
}

func (p *ProductCache) MsetProductInfo(ctx context.Context, mp map[uint64]*model.ProductInfo) error {
	var val []string
	for id, pdi := range mp {
		k := p.generateKey(id)
		v, err := pdi.Write()
		if err != nil {
			return err
		}
		val = append(val, k, v)
	}
	_, err := p.cli.MSet(ctx, val).Result()
	if err != nil {
		p.log.Errorf("mset product info cache failed,map:%v,error:%v", mp, err)
		return err
	}
	return nil
}
func (p *ProductCache) generateKey(id uint64) string {
	return fmt.Sprintf("product:%d", id)
}
func (p *ProductCache) generateExpiration(expire int) time.Duration {
	// 生成一个随机倍数，范围在 [1, 2) 之间
	randomFactor := 1 + rand.Float64() // 生成一个随机浮动值，范围为 [1, 2)
	// 计算新的过期时间
	newExpireTime := time.Duration(float64(expire)*randomFactor) * time.Second
	return newExpireTime
}
