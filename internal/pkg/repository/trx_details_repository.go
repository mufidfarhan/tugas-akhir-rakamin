package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type TrxDetailsRepository interface {
	Transactor
	CreateTrxDetails(ctx context.Context, data entity.TrxDetail) (res uint, err error)
}

type TrxDetailsRepositoryImpl struct {
	transactor
}

func NewTrxDetailsRepository(db *gorm.DB) TrxDetailsRepository {
	return &TrxDetailsRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *TrxDetailsRepositoryImpl) CreateTrxDetails(ctx context.Context, data entity.TrxDetail) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}
