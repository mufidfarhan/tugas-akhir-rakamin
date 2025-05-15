package entity

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID        uint
	AddressTitle  string
	RecipientName string
	PhoneNumber   string
	FullAddress   string
	Trx           Trx `gorm:"constraint:OnDelete:SET NULL;"`
}
