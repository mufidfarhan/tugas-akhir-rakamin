package controller

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type UsersController interface {
	// Profile
	GetMyProfile(ctx *fiber.Ctx) error
	UpdateMyProfile(ctx *fiber.Ctx) error

	// Address
	GetMyAddresses(ctx *fiber.Ctx) error
	GetAddressByID(ctx *fiber.Ctx) error
	AddAddress(ctx *fiber.Ctx) error
	UpdateAddressByID(ctx *fiber.Ctx) error
	DeleteAddressByID(ctx *fiber.Ctx) error
}

type UsersControllerImpl struct {
	usersUseCase usecase.UsersUseCase
}

func NewUsersController(usersUseCase usecase.UsersUseCase) UsersController {
	return &UsersControllerImpl{
		usersUseCase: usersUseCase,
	}
}

func (uc *UsersControllerImpl) GetMyProfile(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	res, err := uc.usersUseCase.GetMyProfile(c, userID)
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

func (uc *UsersControllerImpl) UpdateMyProfile(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	data := new(model.UserReqUpdate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.usersUseCase.UpdateMyProfile(c, userID, *data)
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

func (uc *UsersControllerImpl) GetMyAddresses(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	res, err := uc.usersUseCase.GetMyAddresses(c, userID)
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

func (uc *UsersControllerImpl) GetAddressByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	addressID := ctx.Params("id")
	if addressID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	res, err := uc.usersUseCase.GetAddressByID(c, addressID)
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

func (uc *UsersControllerImpl) AddAddress(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	data := new(model.AddressReqCreate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.usersUseCase.AddUserAddress(c, userID, *data)
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

func (uc *UsersControllerImpl) UpdateAddressByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)
	addressID := ctx.Params("id")
	if addressID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	data := new(model.AddressReqUpdate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.usersUseCase.UpdateAddressByID(c, userID, addressID, *data)
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

func (uc *UsersControllerImpl) DeleteAddressByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)
	addressID := ctx.Params("id")
	if addressID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to DELETE data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	res, err := uc.usersUseCase.DeleteAddressByID(c, userID, addressID)
	if err != nil {
		return ctx.Status(err.Code).JSON(helper.Response{
			Status:  false,
			Message: "Failed to DELETE data",
			Errors:  []string{err.Err.Error()},
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(helper.Response{
		Status:  true,
		Message: "Succeed to DELETE data",
		Errors:  nil,
		Data:    res,
	})
}
