package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"

	"gorm.io/gorm"
)

type CategoriesRepository interface {
	Transactor

	CreateCategory(ctx context.Context, data entity.Category) (res uint, err error)
	GetCategories(ctx context.Context) (res []entity.Category, err error)
	GetCategoryByID(ctx context.Context, categoryID string) (res entity.Category, err error)
	UpdateCategoryByID(ctx context.Context, categoryID string, data entity.Category) (err error)
	DeleteCategoryByID(ctx context.Context, categoryID string) (err error)

	VerifyCategoryAvailability(ctx context.Context, categoryID string) error
}

type CategoriesRepositoryImpl struct {
	transactor
}

func NewCategoriesRepository(db *gorm.DB) CategoriesRepository {
	return &CategoriesRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *CategoriesRepositoryImpl) CreateCategory(ctx context.Context, data entity.Category) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *CategoriesRepositoryImpl) GetCategories(ctx context.Context) (res []entity.Category, err error) {
	if err := r.tx(ctx).Find(&res).Error; err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return res, nil
}

func (r *CategoriesRepositoryImpl) GetCategoryByID(ctx context.Context, categoryID string) (res entity.Category, err error) {
	if err := r.tx(ctx).First(&res, categoryID).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *CategoriesRepositoryImpl) UpdateCategoryByID(ctx context.Context, categoryID string, data entity.Category) (err error) {
	if err := r.tx(ctx).Model(&entity.Category{}).Where("id = ?", categoryID).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoriesRepositoryImpl) DeleteCategoryByID(ctx context.Context, categoryID string) (err error) {
	if err := r.tx(ctx).Delete(&entity.Category{}, categoryID).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoriesRepositoryImpl) VerifyCategoryAvailability(ctx context.Context, categoryID string) error {
	var category entity.Category
	if err := r.tx(ctx).Where("id = ? ", categoryID).First(&category).Error; err != nil {
		return err
	}

	return nil
}
