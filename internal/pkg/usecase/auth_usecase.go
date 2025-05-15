package usecase

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/entity"
	"backend-evermos/internal/pkg/model"
	userModel "backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/repository"
	"backend-evermos/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(ctx context.Context, params userModel.Login) (res userModel.LoginRes, err *helper.ErrorStruct)
	CreateUser(ctx context.Context, data userModel.UserReqCreate) (res uint, err *helper.ErrorStruct)
}

type AuthUseCaseImpl struct {
	usersRepository    repository.UsersRepository
	shopsRepository    repository.ShopsRepository
	provcityRepository repository.ProvcityRepository
}

func NewAuthUseCase(
	usersRepository repository.UsersRepository,
	shopsRepository repository.ShopsRepository,
	provcityRepository repository.ProvcityRepository,
) AuthUseCase {
	return &AuthUseCaseImpl{
		usersRepository:    usersRepository,
		shopsRepository:    shopsRepository,
		provcityRepository: provcityRepository,
	}
}

func (alc *AuthUseCaseImpl) Login(ctx context.Context, params userModel.Login) (res userModel.LoginRes, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.usersRepository.GetUserByPhoneNumber(ctx, params.PhoneNumber)
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusUnauthorized,
			Err:  errors.New("no telp atau kata sandi salah"),
		}
	}

	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAllUsers : %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	isValid := utils.CheckPasswordHash(params.Password, resRepo.Password)
	if !isValid {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusUnauthorized,
			Err:  errors.New("no telp atau kata sandi salah"),
		}
	}

	tokenInit := utils.NewToken(utils.DataClaims{
		ID:      fmt.Sprint(resRepo.ID),
		Email:   resRepo.Email,
		IsAdmin: resRepo.IsAdmin,
	})

	token, errToken := tokenInit.Create()
	if errToken != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusUnauthorized,
			Err:  errToken,
		}
	}

	provinceID := resRepo.ProvinceID
	provRepo, errRepo := alc.provcityRepository.GetProvinceByID(provinceID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("data provinsi tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetProvinceByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	province := model.ProvinceResp{
		ID:   provRepo.ID,
		Name: provRepo.Name,
	}

	cityID := resRepo.CityID
	cityRepo, errRepo := alc.provcityRepository.GetCityByID(cityID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("data kota tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetCityByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	city := model.CityResp{
		ID:         cityRepo.ID,
		ProvinceID: cityRepo.ProvinceID,
		Name:       cityRepo.Name,
	}

	res = userModel.LoginRes{
		Name:        resRepo.Name,
		PhoneNumber: resRepo.PhoneNumber,
		BirthDate:   resRepo.BirthDate.Format("02/01/2006"),
		About:       resRepo.About,
		JobTitle:    resRepo.JobTitle,
		Email:       resRepo.Email,
		ProvinceID:  province,
		CityID:      city,
		Token:       token,
	}

	return res, nil
}

func (alc *AuthUseCaseImpl) CreateUser(ctx context.Context, params userModel.UserReqCreate) (res uint, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(params); errValidate != nil {
		log.Println(errValidate)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	if errPhoneValidator := utils.ValidatePhoneNumber(params.PhoneNumber); errPhoneValidator != nil {
		log.Println(errPhoneValidator)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errPhoneValidator,
		}
	}

	if errRepo := alc.usersRepository.VerifyPhoneNumber(ctx, params.PhoneNumber); errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	if errRepo := alc.usersRepository.VerifyEmail(ctx, params.Email); errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	hashPass, errHash := utils.HashPassword(params.Password)
	if errHash != nil {
		log.Println(errHash)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusInternalServerError,
			Err:  errHash,
		}
	}

	birthDate, errParse := utils.ParseDate(params.BirthDate)
	if errParse != nil {
		log.Println(errParse)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errParse,
		}
	}

	shopName := utils.GenerateShopName(params.Email)

	errTransaction := alc.usersRepository.WithinTransaction(ctx, func(txCtx context.Context) (err error) {
		userID, err := alc.usersRepository.CreateUser(txCtx, entity.User{
			Email:       params.Email,
			Name:        params.Name,
			Password:    hashPass,
			PhoneNumber: params.PhoneNumber,
			BirthDate:   birthDate,
			About:       params.About,
			JobTitle:    params.JobTitle,
			ProvinceID:  params.ProvinceID,
			CityID:      params.CityID,
			IsAdmin:     false,
		})
		if err != nil {
			return err
		}

		_, err = alc.shopsRepository.CreateShop(txCtx, entity.Shop{
			UserID:   userID,
			ShopName: shopName,
			PhotoURL: "",
		})
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at Transactions: %s", errTransaction.Error()), errTransaction)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusInternalServerError,
			Err:  errors.New("gagal menambahkan user"),
		}
	}

	return res, nil
}
