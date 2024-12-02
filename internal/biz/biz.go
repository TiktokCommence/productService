package biz

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet()

type Transaction interface {
	DB(ctx context.Context) *gorm.DB
	InTx(ctx context.Context, fc func(ctx context.Context) error) error
}
