package model

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ProductTableName = "products"
)

type Product struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(100);index"`
	Description string    `gorm:"column:description"`
	PictureUrl  string    `gorm:"column:picture_url"`
	Price       float64   `gorm:"column:price"`
	MerchantID  uint64    `gorm:"column:merchant_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (p *Product) String() string {
	return fmt.Sprintf("%+v", *p)
}
func (p *Product) TableName() string {
	return ProductTableName
}

func (p *Product) Read(val string) error {
	err := json.Unmarshal([]byte(val), &p)
	return err
}

func (p *Product) Write() (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type ProductInfo struct {
	Pd         *Product
	Categories []string
}

func NewProductInfo(pd *Product, cats []string) *ProductInfo {
	return &ProductInfo{
		Pd:         pd,
		Categories: cats,
	}
}
func (p *ProductInfo) String() string {
	return fmt.Sprintf("%+v", *p)
}
func (p *ProductInfo) Write() (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (p *ProductInfo) Read(val string) error {
	err := json.Unmarshal([]byte(val), &p)
	return err
}
