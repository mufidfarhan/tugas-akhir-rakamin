package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName   string
	Slug          string
	ResellerPrice string
	ConsumerPrice string
	Stock         int
	Description   string
	ShopID        uint
	CategoryID    *uint
	Shop          Shop           `gorm:"constraint:OnDelete:CASCADE;"`
	Category      Category       `gorm:"constraint:OnDelete:SET NULL;"`
	Images        []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	ProductLog    ProductLog     `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL;"`
}

type FilterProducts struct {
	Limit, Offset int
	ProductName   string
	CategoryID    uint
	ShopID        uint
	MaxPrice      int
	MinPrice      int
}
