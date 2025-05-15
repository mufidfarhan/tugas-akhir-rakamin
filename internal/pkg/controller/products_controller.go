package controller

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductsController interface {
	AddProduct(ctx *fiber.Ctx) error
	GetAllProducts(ctx *fiber.Ctx) error
	GetProductByID(ctx *fiber.Ctx) error
	UpdateProductByID(ctx *fiber.Ctx) error
	DeleteProductByID(ctx *fiber.Ctx) error
}

type ProductsControllerImpl struct {
	productsUseCase usecase.ProductsUseCase
}

func NewProductsController(productsUseCase usecase.ProductsUseCase) ProductsController {
	return &ProductsControllerImpl{
		productsUseCase: productsUseCase,
	}
}

func (uc *ProductsControllerImpl) AddProduct(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)

	data := new(model.ProductReqCreate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	form, errFile := ctx.MultipartForm()
	if errFile != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{errFile.Error()},
			Data:    nil,
		})
	}

	files := form.File["photos"]
	if len(files) != 2 {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{"Harus melengkapi foto produk"},
			Data:    nil,
		})
	}

	res, err := uc.productsUseCase.CreateProduct(c, userID, *data, files)
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

func (uc *ProductsControllerImpl) GetAllProducts(ctx *fiber.Ctx) error {
	c := ctx.Context()

	filter := new(model.ProductsFilter)
	if err := ctx.QueryParser(filter); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to GET data",
			Errors:  []string{err.Error()},
		})
	}

	res, err := uc.productsUseCase.GetAllProducts(c, model.ProductsFilter{
		Limit:       filter.Limit,
		Page:        filter.Page,
		ProductName: filter.ProductName,
		CategoryID:  filter.CategoryID,
		ShopID:      filter.ShopID,
		MaxPrice:    filter.MaxPrice,
		MinPrice:    filter.MinPrice,
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

func (uc *ProductsControllerImpl) GetProductByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	productID := ctx.Params("id")

	res, err := uc.productsUseCase.GetProductByID(c, productID)
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

func (uc *ProductsControllerImpl) UpdateProductByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)
	productID := ctx.Params("id")

	data := new(model.ProductReqUpdate)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to PUT data",
			Errors:  []string{err.Error()},
			Data:    nil,
		})
	}

	form, errFile := ctx.MultipartForm()
	if errFile != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to POST data",
			Errors:  []string{errFile.Error()},
			Data:    nil,
		})
	}

	files := form.File["photos"]
	res, err := uc.productsUseCase.UpdateProductByID(c, userID, productID, *data, files)
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

func (uc *ProductsControllerImpl) DeleteProductByID(ctx *fiber.Ctx) error {
	c := ctx.Context()
	userID := ctx.Locals("userid").(string)
	productID := ctx.Params("id")
	if productID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.Response{
			Status:  false,
			Message: "Failed to DELETE data",
			Errors:  []string{"Bad request"},
			Data:    nil,
		})
	}

	res, err := uc.productsUseCase.DeleteProductByID(c, userID, productID)
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
