package entity

import "gorm.io/gorm"

type TrxDetail struct {
	gorm.Model
	TrxID        uint
	ProductLogID uint
	ShopID       uint
	Quantity     int
	TotalPrice   int
	ProductLog   ProductLog `gorm:"constraint:OnDelete:SET NULL;"`
	Shop         Shop       `gorm:"constraint:OnDelete:SET NULL;"`
}
