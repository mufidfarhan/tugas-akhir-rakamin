package handler

import (
	productscontroller "backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

func ProductsRoute(r fiber.Router, ProductUsc usecase.ProductsUseCase) {
	controller := productscontroller.NewProductsController(ProductUsc)

	ProductsAPI := r.Group("/product")
	ProductsAPI.Post("", MiddlewareAuth, controller.AddProduct)
	ProductsAPI.Get("", controller.GetAllProducts)
	ProductsAPI.Get("/:id", controller.GetProductByID)
	ProductsAPI.Put("/:id", MiddlewareAuth, controller.UpdateProductByID)
	ProductsAPI.Delete("/:id", MiddlewareAuth, controller.DeleteProductByID)
}
