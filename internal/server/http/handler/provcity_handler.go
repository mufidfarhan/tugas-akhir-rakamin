package handler

import (
	"backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

func ProvcityRoute(r fiber.Router, ProvcityUsc usecase.ProvcityUseCase) {
	controller := controller.NewProvcityController(ProvcityUsc)

	provCityAPI := r.Group("/provcity")
	provCityAPI.Get("/listprovincies", controller.GetListProvinces)
	provCityAPI.Get("/detailprovince/:id", controller.GetProvinceDetails)
	provCityAPI.Get("/listcities/:id", controller.GetListCities)
	provCityAPI.Get("/detailcity/:id", controller.GetCityDetails)
}
