package repository

import (
	"context"
	"github.com/TiktokCommence/productService/internal/biz"
	"github.com/TiktokCommence/productService/internal/conf"
	"github.com/TiktokCommence/productService/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type contextTxKey struct{}

var _ biz.Transaction = (*Gdb)(nil)

type Gdb struct {
	db *gorm.DB
}

func NewGdb(db *gorm.DB) *Gdb {
	return &Gdb{db: db}
}
func NewDB(cf *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cf.Database.Source), &gorm.Config{})
	if err != nil {
		panic("connect mysql failed")
	}
	if err := db.AutoMigrate(&model.Product{}, &model.ProductCategory{}); err != nil {
		panic(err)
	}
	return db
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
