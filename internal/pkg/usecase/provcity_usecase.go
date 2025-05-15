package usecase

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/repository"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProvcityUseCase interface {
	GetListProvinces() (res []model.ProvinceResp, err *helper.ErrorStruct)
	GetProvinceDetails(provinceID string) (res model.ProvinceResp, err *helper.ErrorStruct)
	GetListCities(provinceID string) (res []model.CityResp, err *helper.ErrorStruct)
	GetCityDetails(cityID string) (res model.CityResp, err *helper.ErrorStruct)
}

type ProvcityUseCaseImpl struct {
	provCityRepository repository.ProvcityRepository
}

func NewProvcityUseCase(provCityRepository repository.ProvcityRepository) ProvcityUseCase {
	return &ProvcityUseCaseImpl{
		provCityRepository: provCityRepository,
	}
}

func (alc *ProvcityUseCaseImpl) GetListProvinces() (res []model.ProvinceResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.provCityRepository.GetProvinces()
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("list provinsi tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetProvinces: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		res = append(res, model.ProvinceResp{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	return res, nil
}

func (alc *ProvcityUseCaseImpl) GetProvinceDetails(provinceID string) (res model.ProvinceResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.provCityRepository.GetProvinceByID(provinceID)
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

	res = model.ProvinceResp{
		ID:   resRepo.ID,
		Name: resRepo.Name,
	}

	return res, nil
}

func (alc *ProvcityUseCaseImpl) GetListCities(provinceID string) (res []model.CityResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.provCityRepository.GetCitiesByProvID(provinceID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("list kota tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetCitiesByProvID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		res = append(res, model.CityResp{
			ID:         v.ID,
			ProvinceID: v.ProvinceID,
			Name:       v.Name,
		})
	}

	return res, nil
}

func (alc *ProvcityUseCaseImpl) GetCityDetails(cityID string) (res model.CityResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.provCityRepository.GetCityByID(cityID)
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

	res = model.CityResp{
		ID:         resRepo.ID,
		ProvinceID: resRepo.ProvinceID,
		Name:       resRepo.Name,
	}

	return res, nil
}
