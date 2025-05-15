package entity

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	ProductID uint
	PhotoURL  string
}
