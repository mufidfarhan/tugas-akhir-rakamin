package controller

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type CategoriesController interface {
	AddCategory(ctx *fiber.Ctx) error
	GetCategories(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	UpdateCategoryByID(ctx *fiber.Ctx) error
	DeleteCategoryByID(ctx *fiber.Ctx) error
}

type CategoriesControllerImpl struct {
	categoriesUseCase usecase.CategoriesUseCase
}

func NewCategoriesController(categoriesUseCase usecase.CategoriesUseCase) CategoriesController {
	return &CategoriesControllerImpl{
		categoriesUseCase: categoriesUseCase,
	}
}

func (uc *CategoriesControllerImpl) AddCategory(ctx *fiber.Ctx) error {
	c := ctx.Context()

	data := new(model.CategoryReqCreate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.categoriesUseCase.CreateCategory(c, *data)
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

func (uc *CategoriesControllerImpl) GetCategories(ctx *fiber.Ctx) error {
	c := ctx.Context()

	res, err := uc.categoriesUseCase.GetCategories(c)
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

func (uc *CategoriesControllerImpl) GetCategoryByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	categoryID := ctx.Params("id")
	if categoryID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	res, err := uc.categoriesUseCase.GetCategoryByID(c, categoryID)
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

func (uc *CategoriesControllerImpl) UpdateCategoryByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	categoryID := ctx.Params("id")

	if categoryID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	data := new(model.CategoryReqUpdate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	res, err := uc.categoriesUseCase.UpdateCategoryByID(c, categoryID, *data)
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

func (uc *CategoriesControllerImpl) DeleteCategoryByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	categoryID := ctx.Params("id")

	if categoryID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to DELETE data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	res, err := uc.categoriesUseCase.DeleteCategoryByID(c, categoryID)
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
