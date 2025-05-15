package entity

import "gorm.io/gorm"

type Trx struct {
	gorm.Model
	UserID        uint
	AddressID     uint
	TotalPrice    int
	InvoiceCode   string
	PaymentMethod string
	TrxDetails    []TrxDetail `gorm:"constraint:OnDelete:CASCADE;"`
}

type FilterTrx struct {
	Limit, Offset int
	Search        string
}

type ProductTrx struct {
	Quantity      int
	TotalPrice    int
	ProductID     uint
	ProductName   string
	Slug          string
	ResellerPrice string
	ConsumerPrice string
	Stock         int
	Description   string
	ShopID        uint
	CategoryID    *uint
}
