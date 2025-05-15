package usecase

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/entity"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/repository"
	"backend-evermos/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ShopsUseCase interface {
	GetMyShop(ctx context.Context, userID string) (res model.MyShopResp, err *helper.ErrorStruct)
	UpdateShopByID(ctx context.Context, shopID string, userID string, data model.ShopReqUpdate, fileHeader *multipart.FileHeader) (res string, err *helper.ErrorStruct)
	GetShopByID(ctx context.Context, shopID string) (res model.ShopResp, err *helper.ErrorStruct)
	GetAllShops(ctx context.Context, params model.ShopsFilter) (res model.FilteredData, err *helper.ErrorStruct)
}

type ShopsUseCaseImpl struct {
	shopsRepository repository.ShopsRepository
}

func NewShopsUseCase(shopsRepository repository.ShopsRepository) ShopsUseCase {
	return &ShopsUseCaseImpl{
		shopsRepository: shopsRepository,
	}
}

func (alc *ShopsUseCaseImpl) GetMyShop(ctx context.Context, userID string) (res model.MyShopResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.shopsRepository.GetShopByUserID(ctx, userID)
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusNotFound,
			Err:  errors.New("record not found"),
		}
	}

	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetShopByUserID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	res = model.MyShopResp{
		ID:       resRepo.ID,
		ShopName: resRepo.ShopName,
		PhotoURL: resRepo.PhotoURL,
		UserID:   resRepo.UserID,
	}

	return res, nil
}

func (alc *ShopsUseCaseImpl) UpdateShopByID(ctx context.Context, shopID string, userID string, data model.ShopReqUpdate, fileHeader *multipart.FileHeader) (res string, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		log.Println(errValidate)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	if errRepo := alc.shopsRepository.VerifyShopAvailability(ctx, shopID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("toko tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyShopAvailability: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	if errRepo := alc.shopsRepository.VerifyShopOwner(ctx, shopID, userID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusForbidden,
				Err:  errors.New("anda tidak berhak mengakses resource ini"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyShopOwner: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	uploadDir := "files/shops"
	photoURL, errRepo := utils.SaveFileToDisk(fileHeader, uploadDir)
	if err != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	errRepo = alc.shopsRepository.UpdateShopByID(ctx, shopID, entity.Shop{
		ShopName: data.ShopName,
		PhotoURL: photoURL,
	})
	if errRepo != nil {
		_ = os.Remove(photoURL)
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at UpdateShopByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errors.New("gagal melakukan pembaruan"),
		}
	}

	return "updated", nil
}

func (alc *ShopsUseCaseImpl) GetShopByID(ctx context.Context, shopID string) (res model.ShopResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.shopsRepository.GetShopByID(ctx, shopID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("toko tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetShopByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	res = model.ShopResp{
		ID:       resRepo.ID,
		ShopName: resRepo.ShopName,
		PhotoURL: resRepo.PhotoURL,
	}

	return res, nil
}

func (alc *ShopsUseCaseImpl) GetAllShops(ctx context.Context, params model.ShopsFilter) (res model.FilteredData, err *helper.ErrorStruct) {
	var shops []model.ShopResp

	limit, offset := func(limit, page int) (int, int) {
		if limit < 1 {
			limit = 10
		}

		var offset int
		if page < 1 {
			offset = 0
		} else {
			offset = (page - 1) * limit
		}
		return limit, offset
	}(params.Limit, params.Page)

	resRepo, errRepo := alc.shopsRepository.GetAllShops(ctx, entity.FilterShops{
		Limit:    limit,
		Offset:   offset,
		ShopName: params.ShopName,
	})

	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusNotFound,
			Err:  errors.New("toko tidak ditemukan"),
		}
	}

	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAllShops: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		shops = append(shops, model.ShopResp{
			ID:       v.ID,
			ShopName: v.ShopName,
			PhotoURL: v.PhotoURL,
		})
	}

	res = model.FilteredData{
		Data:  shops,
		Page:  params.Page,
		Limit: params.Limit,
	}

	return res, nil
}
