package usecase

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/entity"
	"backend-evermos/internal/pkg/model"
	usersModel "backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/repository"
	"backend-evermos/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UsersUseCase interface {
	// Profile
	GetMyProfile(ctx context.Context, userID string) (res usersModel.UserResp, err *helper.ErrorStruct)
	UpdateMyProfile(ctx context.Context, userID string, data usersModel.UserReqUpdate) (res string, err *helper.ErrorStruct)

	// Address
	GetMyAddresses(ctx context.Context, userID string) (res []usersModel.AddressResp, err *helper.ErrorStruct)
	GetAddressByID(ctx context.Context, addressID string) (res usersModel.AddressResp, err *helper.ErrorStruct)
	AddUserAddress(ctx context.Context, userID string, data usersModel.AddressReqCreate) (res uint, err *helper.ErrorStruct)
	UpdateAddressByID(ctx context.Context, userID string, addressID string, data usersModel.AddressReqUpdate) (res string, err *helper.ErrorStruct)
	DeleteAddressByID(ctx context.Context, userID string, addressID string) (res string, err *helper.ErrorStruct)
}

type UsersUseCaseImpl struct {
	usersRepository     repository.UsersRepository
	addressesRepository repository.AddressesRepository
	provcityRepository  repository.ProvcityRepository
}

func NewUsersUseCase(
	usersRepository repository.UsersRepository,
	addressesRepository repository.AddressesRepository,
	provcityRepository repository.ProvcityRepository,
) UsersUseCase {
	return &UsersUseCaseImpl{
		usersRepository:     usersRepository,
		addressesRepository: addressesRepository,
		provcityRepository:  provcityRepository,
	}
}

func (alc *UsersUseCaseImpl) GetMyProfile(ctx context.Context, userID string) (res usersModel.UserResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.usersRepository.GetUserByID(ctx, userID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("profil tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetUserByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
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

	res = usersModel.UserResp{
		Name:        resRepo.Name,
		PhoneNumber: resRepo.PhoneNumber,
		BirthDate:   resRepo.BirthDate.Format("02/01/2006"),
		About:       resRepo.About,
		JobTitle:    resRepo.JobTitle,
		Email:       resRepo.Email,
		ProvinceID:  province,
		CityID:      city,
	}

	return res, nil
}

func (alc *UsersUseCaseImpl) UpdateMyProfile(ctx context.Context, userID string, data usersModel.UserReqUpdate) (res string, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	if data.PhoneNumber != "" {
		if errPhoneValidator := utils.ValidatePhoneNumber(data.PhoneNumber); errPhoneValidator != nil {
			log.Println(errPhoneValidator)
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  errPhoneValidator,
			}
		}
	}

	if errRepo := alc.usersRepository.VerifyPhoneNumber(ctx, data.PhoneNumber); errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	if errRepo := alc.usersRepository.VerifyEmail(ctx, data.Email); errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	hashPass, errHash := utils.HashPassword(data.Password)
	if errHash != nil {
		log.Println(errHash)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusInternalServerError,
			Err:  errHash,
		}
	}

	birthDate, errParse := utils.ParseDate(data.BirthDate)
	if errParse != nil {
		log.Println(errParse)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errParse,
		}
	}

	errRepo := alc.usersRepository.UpdateUserByID(ctx, userID, entity.User{
		Name:        data.Name,
		Password:    hashPass,
		PhoneNumber: data.PhoneNumber,
		BirthDate:   birthDate,
		JobTitle:    data.JobTitle,
		About:       data.About,
		Email:       data.Email,
		ProvinceID:  data.ProvinceID,
		CityID:      data.CityID,
	})
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("user tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at UpdateUserByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return "updated", nil
}

func (alc *UsersUseCaseImpl) GetMyAddresses(ctx context.Context, userID string) (res []usersModel.AddressResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.addressesRepository.GetAddressesByUserID(ctx, userID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("alamat tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAddressesByUserID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		res = append(res, usersModel.AddressResp{
			ID:            v.ID,
			AddressTitle:  v.AddressTitle,
			RecipientName: v.RecipientName,
			PhoneNumber:   v.PhoneNumber,
			FullAddress:   v.FullAddress,
		})
	}

	return res, nil
}

func (alc *UsersUseCaseImpl) GetAddressByID(ctx context.Context, addressID string) (res usersModel.AddressResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.addressesRepository.GetAddressByID(ctx, addressID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("alamat tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAddressByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	res = usersModel.AddressResp{
		ID:            resRepo.ID,
		AddressTitle:  resRepo.AddressTitle,
		RecipientName: resRepo.RecipientName,
		PhoneNumber:   resRepo.PhoneNumber,
		FullAddress:   resRepo.FullAddress,
	}

	return res, nil
}

func (alc *UsersUseCaseImpl) AddUserAddress(ctx context.Context, userID string, data usersModel.AddressReqCreate) (res uint, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	resRepo, errRepo := alc.addressesRepository.CreateAddress(ctx, userID, entity.Address{
		AddressTitle:  data.AddressTitle,
		RecipientName: data.RecipientName,
		PhoneNumber:   data.PhoneNumber,
		FullAddress:   data.FullAddress,
	})
	if errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at CreateAddress: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return resRepo, nil
}

func (alc *UsersUseCaseImpl) UpdateAddressByID(ctx context.Context, userID string, addressID string, data usersModel.AddressReqUpdate) (res string, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	errRepo := alc.addressesRepository.UpdateAddressByID(ctx, userID, addressID, entity.Address{
		RecipientName: data.RecipientName,
		PhoneNumber:   data.PhoneNumber,
		FullAddress:   data.FullAddress,
	})
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("alamat tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at UpdateAddressByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return "updated", nil
}

func (alc *UsersUseCaseImpl) DeleteAddressByID(ctx context.Context, userID string, addressID string) (res string, err *helper.ErrorStruct) {
	errRepo := alc.addressesRepository.DeleteAddressByID(ctx, userID, addressID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("alamat tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at DeleteAddressByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return "deleted", nil
}
