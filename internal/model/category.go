package model

import (
	"encoding/json"
	"time"
)

const ProductCategoryTableName = "categories"

type ProductCategory struct {
	Pid       uint64    `gorm:"column:p_id;index"`
	Category  string    `gorm:"column:category;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
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
