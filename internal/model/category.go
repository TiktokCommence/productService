package model

import (
	"encoding/json"
	"fmt"
	"time"
)

const ProductCategoryTableName = "categories"

type ProductCategory struct {
	//使用唯一联合索引，（最左前缀法则）可以通过pid来查询
	Pid       uint64    `gorm:"column:p_id;uniqueIndex:uidx_category,priority:1"`
	Category  string    `gorm:"column:category;type:varchar(50);uniqueIndex:uidx_category,priority:2"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (p *ProductCategory) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *ProductCategory) TableName() string {
	return ProductCategoryTableName
}
func (p *ProductCategory) Read(val string) error {
	err := json.Unmarshal([]byte(val), &p)
	return err
}
func (p *ProductCategory) Write() (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
