package http

import (
	route "backend-evermos/internal/server/http/handler"

	"backend-evermos/internal/infrastructure/container"

	"github.com/gofiber/fiber/v2"
)

func HTTPRouteInit(r *fiber.App, containerConf *container.Container) {
	api := r.Group("/api/v1") // /api

	route.AuthRoute(api, containerConf.AuthUsc)
	route.UsersRoute(api, containerConf.UsersUsc)
	route.ShopsRoute(api, containerConf.ShopsUsc)
	route.ProductsRoute(api, containerConf.ProductsUsc)
	route.CategoriesRoute(api, containerConf.CategoriesUsc)
	route.TrxRoute(api, containerConf.TrxUsc)
	route.ProvcityRoute(api, containerConf.ProvcityUsc)
}
