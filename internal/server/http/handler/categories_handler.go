package handler

import (
	categoriescontroller "backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

func CategoriesRoute(r fiber.Router, CategoryUsc usecase.CategoriesUseCase) {
	controller := categoriescontroller.NewCategoriesController(CategoryUsc)

	CategoriesAPI := r.Group("/category")
	CategoriesAPI.Get("", controller.GetCategories)
	CategoriesAPI.Get("/:id", controller.GetCategoryByID)
	CategoriesAPI.Post("", MiddlewareAuth, MiddlewareAuthAdmin, controller.AddCategory)
	CategoriesAPI.Put("/:id", MiddlewareAuth, MiddlewareAuthAdmin, controller.UpdateCategoryByID)
	CategoriesAPI.Delete("/:id", MiddlewareAuth, MiddlewareAuthAdmin, controller.DeleteCategoryByID)
}
