package handler

import (
	shopscontroller "backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

func ShopsRoute(r fiber.Router, ShopUsc usecase.ShopsUseCase) {
	controller := shopscontroller.NewShopsController(ShopUsc)

	shopsAPI := r.Group("/toko")
	shopsAPI.Get("/my", MiddlewareAuth, controller.GetMyShop)
	shopsAPI.Put("/:id", MiddlewareAuth, controller.UpdateShopByID)
	shopsAPI.Get("/:id", MiddlewareAuth, controller.GetShopByID)
	shopsAPI.Get("", controller.GetAllShops)
}
