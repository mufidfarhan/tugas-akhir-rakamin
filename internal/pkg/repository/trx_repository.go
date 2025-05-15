package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type TrxRepository interface {
	Transactor

	CreateTrx(ctx context.Context, data entity.Trx) (res uint, err error)
	GetTrxByID(ctx context.Context, userID string, trxID string) (res entity.Trx, err error)
	GetAllTrxByUserID(ctx context.Context, userID string, params entity.FilterTrx) (res []entity.Trx, err error)
}

type TrxRepositoryImpl struct {
	transactor
}

func NewTrxRepository(db *gorm.DB) TrxRepository {
	return &TrxRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *TrxRepositoryImpl) CreateTrx(ctx context.Context, data entity.Trx) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *TrxRepositoryImpl) GetTrxByID(ctx context.Context, userID string, trxID string) (res entity.Trx, err error) {
	db := r.tx(ctx).
		Preload("TrxDetails").
		Preload("TrxDetails.ProductLog").
		Preload("TrxDetails.Shop").
		Preload("TrxDetails.ProductLog.Shop").
		Preload("TrxDetails.ProductLog.Category")

	if err := db.Where("user_id = ?", userID).First(&res, trxID).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *TrxRepositoryImpl) GetAllTrxByUserID(ctx context.Context, userID string, params entity.FilterTrx) (res []entity.Trx, err error) {
	db := r.tx(ctx).
		Joins("JOIN trx_details ON trx_details.trx_id = trxes.id").
		Joins("JOIN product_logs ON product_logs.id = trx_details.product_log_id").
		Preload("TrxDetails").
		Preload("TrxDetails.ProductLog").
		Preload("TrxDetails.Shop").
		Preload("TrxDetails.ProductLog.Shop").
		Preload("TrxDetails.ProductLog.Category")

	if params.Search != "" {
		db = db.Where("product_logs.product_name LIKE ?", "%"+params.Search+"%")
	}

	if err := db.Where("user_id = ?", userID).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}
