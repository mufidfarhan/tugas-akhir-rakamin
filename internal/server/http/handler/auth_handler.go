package handler

import (
	"github.com/gofiber/fiber/v2"

	authcontroller "backend-evermos/internal/pkg/controller"
	"backend-evermos/internal/pkg/usecase"
)

func AuthRoute(r fiber.Router, AuthUsc usecase.AuthUseCase) {
	controller := authcontroller.NewAuthController(AuthUsc)

	booksAPI := r.Group("/auth")
	booksAPI.Post("/register", controller.Register)
	booksAPI.Post("/login", controller.Login)
}
