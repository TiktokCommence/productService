package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrorNoWhere = errors.New("no where for query")
)

type contextTxKey struct{}

type Gdb struct {
	db *gorm.DB
}

func (g *Gdb) DB(ctx context.Context) *gorm.DB {
	// 从ctx中获取tx
	txKey := ctx.Value(contextTxKey{})
	tx, ok := txKey.(*gorm.DB)
	if ok {
		return tx
	}
	// Notice 如果 !ok 返回错误还是使用默认DB～这个根据实际情况来定！
	return g.db
}

func (g *Gdb) InTx(ctx context.Context, fc func(ctx context.Context) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将tx放入到ctx中
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fc(ctx)
	})
}
