package handler

import (
	userscontroller "backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

func UsersRoute(r fiber.Router, UserUsc usecase.UsersUseCase) {
	controller := userscontroller.NewUsersController(UserUsc)

	usersAPI := r.Group("/user")
	usersAPI.Get("", MiddlewareAuth, controller.GetMyProfile)
	usersAPI.Put("", MiddlewareAuth, controller.UpdateMyProfile)
	usersAPI.Get("/alamat", MiddlewareAuth, controller.GetMyAddresses)
	usersAPI.Get("/alamat/:id", MiddlewareAuth, controller.GetAddressByID)
	usersAPI.Post("/alamat", MiddlewareAuth, controller.AddAddress)
	usersAPI.Put("/alamat/:id", MiddlewareAuth, controller.UpdateAddressByID)
	usersAPI.Delete("/alamat/:id", MiddlewareAuth, controller.DeleteAddressByID)
}
