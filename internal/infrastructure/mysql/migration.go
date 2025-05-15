package mysql

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/entity"

	"gorm.io/gorm"
)

func RunMigration(mysqlDB *gorm.DB) {
	err := mysqlDB.AutoMigrate(
		&entity.User{},
		&entity.Address{},
		&entity.Shop{},
		&entity.Product{},
		&entity.ProductImage{},
		&entity.Category{},
		&entity.Trx{},
		&entity.TrxDetail{},
		&entity.ProductLog{},
	)
	if err != nil {
		helper.Logger(helper.LoggerLevelError, "Failed Database Migrated", err)
	}

	helper.Logger(helper.LoggerLevelInfo, "Database Migrated", nil)
}
