package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type ProductsRepository interface {
	Transactor

	CreateProduct(ctx context.Context, data entity.Product) (res uint, err error)
	GetAllProducts(ctx context.Context, params entity.FilterProducts) (res []entity.Product, err error)
	GetProductByID(ctx context.Context, productID string) (res entity.Product, err error)
	UpdateProductByID(ctx context.Context, productID string, data entity.Product) (err error)
	DeleteProductByID(ctx context.Context, productID string) (err error)

	VerifyProductAvailability(ctx context.Context, productID string) (err error)
	VerifyProductOwner(ctx context.Context, productID string, shopID string) (err error)
}

type ProductsRepositoryImpl struct {
	transactor
}

func NewProductsRepository(db *gorm.DB) ProductsRepository {
	return &ProductsRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *ProductsRepositoryImpl) CreateProduct(ctx context.Context, data entity.Product) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *ProductsRepositoryImpl) GetAllProducts(ctx context.Context, params entity.FilterProducts) (res []entity.Product, err error) {
	db := r.tx(ctx).
		Preload("Shop").
		Preload("Category").
		Preload("Images")

	if params.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+params.ProductName+"%")
	}
	if params.CategoryID > 0 {
		db = db.Where("category_id LIKE ?", params.CategoryID)
	}
	if params.ShopID > 0 {
		db = db.Where("shop_id LIKE ?", params.ShopID)
	}
	if params.MinPrice > 0 {
		db = db.Where("consumer_price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		db = db.Where("consumer_price <= ?", params.MaxPrice)
	}

	if err := db.Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *ProductsRepositoryImpl) GetProductByID(ctx context.Context, productID string) (res entity.Product, err error) {
	db := r.tx(ctx).Preload("Shop").Preload("Category").Preload("Images")
	if err := db.First(&res, productID).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *ProductsRepositoryImpl) UpdateProductByID(ctx context.Context, productID string, data entity.Product) (err error) {
	if err := r.tx(ctx).Model(&entity.Product{}).Where("id = ?", productID).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepositoryImpl) DeleteProductByID(ctx context.Context, productID string) (err error) {
	if err := r.tx(ctx).Delete(&entity.Product{}, productID).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepositoryImpl) VerifyProductAvailability(ctx context.Context, productID string) (err error) {
	var product entity.Product
	if err := r.tx(ctx).Where("id = ? ", productID).First(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepositoryImpl) VerifyProductOwner(ctx context.Context, productID string, shopID string) (err error) {
	var product entity.Product
	if err := r.tx(ctx).Where("id = ? AND shop_id = ?", productID, shopID).First(&product).Error; err != nil {
		return err
	}

	return nil
}
