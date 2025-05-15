package repository

import (
	"backend-evermos/internal/infrastructure/restclient"
	"backend-evermos/internal/pkg/entity"
	"fmt"
)

type ProvcityRepository interface {
	GetProvinces() (res []entity.Province, err error)
	GetProvinceByID(provinceID string) (res entity.Province, err error)
	GetCitiesByProvID(provinceID string) (res []entity.City, err error)
	GetCityByID(cityID string) (res entity.City, err error)
}

type ProvcityRepositoryImpl struct {
	client *restclient.RestClient
}

func NewProvcityRepository(client *restclient.RestClient) ProvcityRepository {
	return &ProvcityRepositoryImpl{
		client: client,
	}
}

func (r *ProvcityRepositoryImpl) GetProvinces() (res []entity.Province, err error) {
	var provinces []entity.Province
	url := "https://emsifa.github.io/api-wilayah-indonesia/api/provinces.json"

	err = r.client.Get(url, &provinces)
	if err != nil {
		return res, err
	}

	return provinces, err
}

func (r *ProvcityRepositoryImpl) GetProvinceByID(provinceID string) (res entity.Province, err error) {
	var province entity.Province
	url := fmt.Sprintf("https://emsifa.github.io/api-wilayah-indonesia/api/province/%s.json", provinceID)

	err = r.client.Get(url, &province)
	if err != nil {
		return res, err
	}

	return province, nil
}

func (r *ProvcityRepositoryImpl) GetCitiesByProvID(provinceID string) (res []entity.City, err error) {
	var cities []entity.City
	url := fmt.Sprintf("https://emsifa.github.io/api-wilayah-indonesia/api/regencies/%s.json", provinceID)

	err = r.client.Get(url, &cities)
	if err != nil {
		return res, err
	}

	return cities, err
}

func (r *ProvcityRepositoryImpl) GetCityByID(cityID string) (res entity.City, err error) {
	var city entity.City
	url := fmt.Sprintf("https://emsifa.github.io/api-wilayah-indonesia/api/regency/%s.json", cityID)

	err = r.client.Get(url, &city)
	if err != nil {
		return res, err
	}

	return city, nil
}
