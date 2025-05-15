package controller

import (
	"backend-evermos/internal/helper"
	model "backend-evermos/internal/pkg/model"
	authUsc "backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type AuthControllerImpl struct {
	authUsc authUsc.AuthUseCase
}

func NewAuthController(authUsc authUsc.AuthUseCase) AuthController {
	return &AuthControllerImpl{
		authUsc: authUsc,
	}
}

func (uc *AuthControllerImpl) Login(ctx *fiber.Ctx) error {
	c := ctx.Context()

	data := new(model.Login)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.authUsc.Login(c, *data)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to POST data",
		Errors:  nil,
		Data:    res,
	})
}

func (uc *AuthControllerImpl) Register(ctx *fiber.Ctx) error {
	c := ctx.Context()

	data := new(model.UserReqCreate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	_, err := uc.authUsc.CreateUser(c, *data)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to POST data",
		Errors:  nil,
		Data:    "Register Succeed",
	})
}
