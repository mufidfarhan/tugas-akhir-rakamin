package controller

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type ShopsController interface {
	GetMyShop(ctx *fiber.Ctx) error
	UpdateShopByID(ctx *fiber.Ctx) error
	GetShopByID(ctx *fiber.Ctx) error
	GetAllShops(ctx *fiber.Ctx) error
}

type ShopsControllerImpl struct {
	shopsUseCase usecase.ShopsUseCase
}

func NewShopsController(shopsUseCase usecase.ShopsUseCase) ShopsController {
	return &ShopsControllerImpl{
		shopsUseCase: shopsUseCase,
	}
}

func (uc *ShopsControllerImpl) GetMyShop(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	res, err := uc.shopsUseCase.GetMyShop(c, userID)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *ShopsControllerImpl) UpdateShopByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)
	shopID := ctx.Params("id")

	data := new(model.ShopReqUpdate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	fileHeader, errFile := ctx.FormFile("photo")
	if errFile != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{errFile.Error()},
		})
	}

	res, err := uc.shopsUseCase.UpdateShopByID(c, shopID, userID, *data, fileHeader)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to PUT data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *ShopsControllerImpl) GetShopByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	shopID := ctx.Params("id")

	res, err := uc.shopsUseCase.GetShopByID(c, shopID)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to GET data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *ShopsControllerImpl) GetAllShops(ctx *fiber.Ctx) error {
	c := ctx.Context()

	filter := new(model.ShopsFilter)
	if err := ctx.QueryParser(filter); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Error()},
		})
	}

	res, err := uc.shopsUseCase.GetAllShops(c, model.ShopsFilter{
		Limit:    filter.Limit,
		Page:     filter.Page,
		ShopName: filter.ShopName,
	})
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
