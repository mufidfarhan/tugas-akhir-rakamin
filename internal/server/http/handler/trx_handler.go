package handler

import (
	trxcontroller "backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

func TrxRoute(r fiber.Router, TrxUsc usecase.TrxUseCase) {
	controller := trxcontroller.NewTrxController(TrxUsc)

	trxAPI := r.Group("/trx")
	trxAPI.Post("", MiddlewareAuth, controller.CreateTrx)
	trxAPI.Get("/:id", MiddlewareAuth, controller.GetTrxByID)
	trxAPI.Get("", MiddlewareAuth, controller.GetAllTrx)
}
