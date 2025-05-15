package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type ProductLogsRepository interface {
	Transactor
	CreateProductLogs(ctx context.Context, data entity.ProductLog) (res uint, err error)
}

type ProductLogsRepositoryImpl struct {
	transactor
}

func NewProductLogsRepository(db *gorm.DB) ProductLogsRepository {
	return &ProductLogsRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *ProductLogsRepositoryImpl) CreateProductLogs(ctx context.Context, data entity.ProductLog) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}
