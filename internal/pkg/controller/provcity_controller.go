package controller

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProvcityController interface {
	GetListProvinces(ctx *fiber.Ctx) error
	GetProvinceDetails(ctx *fiber.Ctx) error
	GetListCities(ctx *fiber.Ctx) error
	GetCityDetails(ctx *fiber.Ctx) error
}

type ProvcityControllerImpl struct {
	ProvcityUseCase usecase.ProvcityUseCase
}

func NewProvcityController(ProvcityUseCase usecase.ProvcityUseCase) ProvcityController {
	return &ProvcityControllerImpl{
		ProvcityUseCase: ProvcityUseCase,
	}
}

func (uc *ProvcityControllerImpl) GetListProvinces(ctx *fiber.Ctx) error {
	res, err := uc.ProvcityUseCase.GetListProvinces()
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Err.Error()},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *ProvcityControllerImpl) GetProvinceDetails(ctx *fiber.Ctx) error {
	provinceID := ctx.Params("id")

	res, err := uc.ProvcityUseCase.GetProvinceDetails(provinceID)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Err.Error()},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *ProvcityControllerImpl) GetListCities(ctx *fiber.Ctx) error {
	provinceID := ctx.Params("id")

	res, err := uc.ProvcityUseCase.GetListCities(provinceID)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Err.Error()},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *ProvcityControllerImpl) GetCityDetails(ctx *fiber.Ctx) error {
	cityID := ctx.Params("id")

	res, err := uc.ProvcityUseCase.GetCityDetails(cityID)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Err.Error()},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    res,
	})
}
