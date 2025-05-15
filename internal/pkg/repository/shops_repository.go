package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type ShopsRepository interface {
	Transactor

	CreateShop(ctx context.Context, data entity.Shop) (res uint, err error)
	GetShopByUserID(ctx context.Context, userID string) (res entity.Shop, err error)
	UpdateShopByID(ctx context.Context, shopID string, data entity.Shop) (err error)
	GetShopByID(ctx context.Context, shopID string) (res entity.Shop, err error)
	GetAllShops(ctx context.Context, params entity.FilterShops) (res []entity.Shop, err error)

	VerifyShopAvailability(ctx context.Context, shopID string) error
	VerifyShopOwner(ctx context.Context, shopID string, userID string) error
}

type ShopsRepositoryImpl struct {
	transactor
}

func NewShopsRepository(db *gorm.DB) ShopsRepository {
	return &ShopsRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *ShopsRepositoryImpl) CreateShop(ctx context.Context, data entity.Shop) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *ShopsRepositoryImpl) GetShopByUserID(ctx context.Context, userID string) (res entity.Shop, err error) {
	if err := r.tx(ctx).First(&res, userID).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *ShopsRepositoryImpl) UpdateShopByID(ctx context.Context, shopID string, data entity.Shop) (err error) {
	if err := r.tx(ctx).Model(&entity.Shop{}).Where("id = ?", shopID).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShopsRepositoryImpl) GetShopByID(ctx context.Context, shopID string) (res entity.Shop, err error) {
	if err := r.tx(ctx).First(&res, shopID).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *ShopsRepositoryImpl) GetAllShops(ctx context.Context, params entity.FilterShops) (res []entity.Shop, err error) {
	db := r.tx(ctx)

	keyword := "%" + params.ShopName + "%"
	db = db.Where("shop_name LIKE ?", keyword)

	if err := db.Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return res, nil
}

func (r *ShopsRepositoryImpl) VerifyShopAvailability(ctx context.Context, shopID string) error {
	var shop entity.Shop
	if err := r.tx(ctx).Where("id = ? ", shopID).First(&shop).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShopsRepositoryImpl) VerifyShopOwner(ctx context.Context, shopID string, userID string) error {
	var shop entity.Shop
	if err := r.tx(ctx).Where("id = ? AND user_id = ?", shopID, userID).First(&shop).Error; err != nil {
		return err
	}

	return nil
}
