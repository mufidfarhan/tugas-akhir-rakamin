package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	Password    string
	PhoneNumber string `gorm:"unique"`
	BirthDate   time.Time
	Gender      string
	About       string
	JobTitle    string
	Email       string `gorm:"unique"`
	ProvinceID  string
	CityID      string
	IsAdmin     bool
	Address     []Address `gorm:"constraint:OnDelete:CASCADE;"`
	Shop        Shop      `gorm:"constraint:OnDelete:CASCADE;"`
	Trx         Trx       `gorm:"constraint:OnDelete:SET NULL;"`
}

type FilterUser struct {
	Limit, Offset int
	Title         string
}
