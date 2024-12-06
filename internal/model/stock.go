package model

const StockTable = "stocks"

type Stock struct {
	PID   uint64 `gorm:"column:p_id"`
	Stock uint   `gorm:"column:stock"`
}

func (s *Stock) TableName() string {
	return StockTable
}
