package cache

import (
	"context"
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	IdKey = "productID"
)

var _ biz.GenerateIDer = (*GenerateIDImplement)(nil)

type GenerateIDImplement struct {
	client *redis.Client
}

func (g *GenerateIDImplement) GenerateID() (uint64, error) {
	timeStamp := uint64(time.Now().Unix())
	res, err := g.client.Incr(context.Background(), IdKey).Result()
	if err != nil {
		return 0, err
	}
	return timeStamp<<32 | uint64(res), nil
}

func NewGenerateIDImplement(cli *redis.Client) *GenerateIDImplement {
	return &GenerateIDImplement{
		client: cli,
	}
}
