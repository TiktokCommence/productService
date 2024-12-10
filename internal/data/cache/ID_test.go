package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestGenerateIDImplement_GenerateID(t *testing.T) {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:16379",
		DB:   0,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	imp := NewGenerateIDImplement(cli)
	id, err := imp.GenerateID()
	if err != nil {
		t.Error(err)
	}
	t.Log("Generated ID: ", id)
}
