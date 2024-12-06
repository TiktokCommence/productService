package model

const SaleTable = "SaleTable"

type ProductSale struct {
	Pid  uint   `gorm:"column:p_id"`
	Sale uint64 `gorm:"column:sale"`
}

func (s *ProductSale) TableName() string {
	return SaleTable
}
