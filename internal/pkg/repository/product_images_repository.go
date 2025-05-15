package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type ProductImagesRepository interface {
	Transactor

	CreateProductImage(ctx context.Context, data entity.ProductImage) (res uint, err error)
	UpdateProductImage(ctx context.Context, imageID string, data entity.ProductImage) (err error)
	GetImagesByProductID(ctx context.Context, productID uint) (res []entity.ProductImage, err error)
}

type ProductImagesRepositoryImpl struct {
	transactor
}

func NewProductImagesRepository(db *gorm.DB) ProductImagesRepository {
	return &ProductImagesRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *ProductImagesRepositoryImpl) CreateProductImage(ctx context.Context, data entity.ProductImage) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *ProductImagesRepositoryImpl) UpdateProductImage(ctx context.Context, imageID string, data entity.ProductImage) (err error) {
	if err := r.tx(ctx).Model(&entity.ProductImage{}).Where("id = ?", imageID).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductImagesRepositoryImpl) GetImagesByProductID(ctx context.Context, productID uint) (res []entity.ProductImage, err error) {
	if err := r.tx(ctx).Where("product_id = ?", productID).Find(&res).Error; err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return res, nil
}
