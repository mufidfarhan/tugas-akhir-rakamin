package entity

import "gorm.io/gorm"

type ProductLog struct {
	gorm.Model
	ProductID     uint
	ProductName   string
	Slug          string
	ResellerPrice string
	ConsumerPrice string
	Description   string
	ShopID        uint
	CategoryID    *uint
	Shop          Shop     `gorm:"constraint:OnDelete:SET NULL;"`
	Category      Category `gorm:"constraint:OnDelete:SET NULL;"`
}
