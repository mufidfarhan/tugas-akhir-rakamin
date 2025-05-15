package controller

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type TrxController interface {
	CreateTrx(ctx *fiber.Ctx) error
	GetTrxByID(ctx *fiber.Ctx) error
	GetAllTrx(ctx *fiber.Ctx) error
}

type TrxControllerImpl struct {
	trxUseCase usecase.TrxUseCase
}

func NewTrxController(trxUseCase usecase.TrxUseCase) TrxController {
	return &TrxControllerImpl{
		trxUseCase: trxUseCase,
	}
}

func (uc *TrxControllerImpl) CreateTrx(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	data := new(model.TrxReqCreate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.trxUseCase.CreateTrx(c, userID, *data)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to POST data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *TrxControllerImpl) GetTrxByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)
	trxID := ctx.Params("id")

	res, err := uc.trxUseCase.GetTrxByID(c, userID, trxID)
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

func (uc *TrxControllerImpl) GetAllTrx(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	filter := new(model.TrxFilter)
	if err := ctx.QueryParser(filter); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Error()},
		})
	}

	res, err := uc.trxUseCase.GetAllTrx(c, userID, model.TrxFilter{
		Limit:  filter.Limit,
		Page:   filter.Page,
		Search: filter.Search,
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
