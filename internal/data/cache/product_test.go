package cache

import (
	"context"
	"errors"
	"github.com/TiktokCommence/productService/internal/model"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"os"
	"testing"
)

func initCache() *ProductCache {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:16379",
		DB:   0,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	logger := log.NewStdLogger(os.Stdout)
	return &ProductCache{
		cli: cli,
		log: log.NewHelper(logger),
	}
}

func TestProductCache_SetProductInfo(t *testing.T) {
	c := initCache()
	pi := &model.ProductInfo{
		Pd: &model.Product{
			ID:   1,
			Name: "甜甜圈",
		},
		Categories: []string{"food"},
		//...
	}
	err := c.SetProductInfo(context.Background(), pi.Pd.ID, pi, 60)
	if err != nil {
		t.Error(err)
	}
	t.Log("set product info cache success")
}
func TestProductCache_GetProductInfo(t *testing.T) {
	c := initCache()
	pi := &model.ProductInfo{
		Pd: &model.Product{
			ID:   1,
			Name: "甜甜圈",
		},
		Categories: []string{"food"},
		//...
	}
	err := c.SetProductInfo(context.Background(), pi.Pd.ID, pi, 60)
	if err != nil {
		t.Error(err)
	}
	got, err := c.GetProductInfo(context.Background(), pi.Pd.ID)
	if err != nil {
		t.Error(err)
	}
	t.Log(*got.Pd)
	t.Log(got.Categories)
}
func TestProductCache_MsetProductInfo(t *testing.T) {
	c := initCache()
	pis := []*model.ProductInfo{
		{
			Pd: &model.Product{
				ID:   1,
				Name: "小明",
			},
			Categories: []string{"food"},
			//...
		},
		{
			Pd: &model.Product{
				ID:   2,
				Name: "小爱",
			},
			Categories: []string{"food"},
			//...
		},
	}
	mp := map[uint64]*model.ProductInfo{
		pis[0].Pd.ID: pis[0],
		pis[1].Pd.ID: pis[1],
	}
	err := c.MsetProductInfo(context.Background(), mp)
	if err != nil {
		t.Error(err)
	}
	t.Log("mset product info cache success")
}
func TestProductCache_MgetProductInfo(t *testing.T) {
	c := initCache()
	pis := []*model.ProductInfo{
		{
			Pd: &model.Product{
				ID:   1,
				Name: "小明",
			},
			Categories: []string{"food"},
			//...
		},
		{
			Pd: &model.Product{
				ID:   2,
				Name: "小爱",
			},
			Categories: []string{"food"},
			//...
		},
	}
	mp := map[uint64]*model.ProductInfo{
		pis[0].Pd.ID: pis[0],
		pis[1].Pd.ID: pis[1],
	}
	err := c.MsetProductInfo(context.Background(), mp)
	if err != nil {
		t.Error(err)
	}

	ids, infos, err := c.MgetProductInfo(context.Background(), []uint64{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
	for i, info := range infos {
		t.Log(i, *info.Pd)
		t.Log(info.Categories)
	}
}
func TestProductCache_DeleteProductInfo(t *testing.T) {
	c := initCache()
	pi := &model.ProductInfo{
		Pd: &model.Product{
			ID:   1,
			Name: "甜甜圈",
		},
		Categories: []string{"food"},
		//...
	}
	err := c.SetProductInfo(context.Background(), pi.Pd.ID, pi, 60)
	if err != nil {
		t.Error(err)
	}
	t.Log("set product info cache success")

	err = c.DeleteProductInfo(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
	t.Log("delete product info cache success")
	_, err = c.GetProductInfo(context.Background(), 1)
	if errors.Is(err, redis.Nil) {
		t.Log("pass")
	} else {
		t.Error(err)
	}
}
