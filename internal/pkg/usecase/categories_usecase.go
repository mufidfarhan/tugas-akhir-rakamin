package usecase

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/entity"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/repository"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoriesUseCase interface {
	CreateCategory(ctx context.Context, data model.CategoryReqCreate) (res uint, err *helper.ErrorStruct)
	GetCategories(ctx context.Context) (res []model.CategoryResp, err *helper.ErrorStruct)
	GetCategoryByID(ctx context.Context, categoryID string) (res model.CategoryResp, err *helper.ErrorStruct)
	UpdateCategoryByID(ctx context.Context, categoryID string, data model.CategoryReqUpdate) (res string, err *helper.ErrorStruct)
	DeleteCategoryByID(ctx context.Context, categoryID string) (res string, err *helper.ErrorStruct)
}

type CategoriesUseCaseImpl struct {
	categoriesRepository repository.CategoriesRepository
	shopsRepository      repository.ShopsRepository
	productsRepository   repository.ProductsRepository
}

func NewCategoriesUseCase(
	categoriesRepository repository.CategoriesRepository,
	shopsRepository repository.ShopsRepository,
	productsRepository repository.ProductsRepository,
) CategoriesUseCase {
	return &CategoriesUseCaseImpl{
		categoriesRepository: categoriesRepository,
		shopsRepository:      shopsRepository,
		productsRepository:   productsRepository,
	}
}

func (alc *CategoriesUseCaseImpl) CreateCategory(ctx context.Context, data model.CategoryReqCreate) (res uint, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		log.Println(errValidate)
		return res, &helper.ErrorStruct{
			Err:  errValidate,
			Code: fiber.StatusBadRequest,
		}
	}

	resRepo, errRepo := alc.categoriesRepository.CreateCategory(ctx, entity.Category{
		CategoryName: data.CategoryName,
	})
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAllBooks : %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return resRepo, nil
}

func (alc *CategoriesUseCaseImpl) GetCategories(ctx context.Context) (res []model.CategoryResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.categoriesRepository.GetCategories(ctx)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("kategori tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetCategories: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		res = append(res, model.CategoryResp{
			ID:           v.ID,
			CategoryName: v.CategoryName,
		})
	}

	return res, nil
}

func (alc *CategoriesUseCaseImpl) GetCategoryByID(ctx context.Context, categoryID string) (res model.CategoryResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.categoriesRepository.GetCategoryByID(ctx, categoryID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("kategori tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetCategoryByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	res = model.CategoryResp{
		ID:           resRepo.ID,
		CategoryName: resRepo.CategoryName,
	}

	return res, nil
}

func (alc *CategoriesUseCaseImpl) UpdateCategoryByID(ctx context.Context, categoryID string, data model.CategoryReqUpdate) (res string, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		log.Println(errValidate)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	if errRepo := alc.categoriesRepository.VerifyCategoryAvailability(ctx, categoryID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("kategori tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyCategoryAvailability: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	errRepo := alc.categoriesRepository.UpdateCategoryByID(ctx, categoryID, entity.Category{
		CategoryName: data.CategoryName,
	})
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at UpdateCategoryByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errors.New("gagal melakukan pembaruan"),
		}
	}

	return "updated", nil
}

func (alc *CategoriesUseCaseImpl) DeleteCategoryByID(ctx context.Context, categoryID string) (res string, err *helper.ErrorStruct) {
	if errRepo := alc.categoriesRepository.VerifyCategoryAvailability(ctx, categoryID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("kategori tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyCategoryAvailability: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	errRepo := alc.categoriesRepository.DeleteCategoryByID(ctx, categoryID)
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at DeleteAddressByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return "deleted", nil
}
