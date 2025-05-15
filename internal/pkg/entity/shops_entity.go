package entity

import "gorm.io/gorm"

type Shop struct {
	gorm.Model
	UserID   uint
	ShopName string
	PhotoURL string
}

type FilterShops struct {
	Limit, Offset int
	ShopName      string
}
