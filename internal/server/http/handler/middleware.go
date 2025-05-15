package handler

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func MiddlewareAuth(ctx *fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helper.Response{
			Status:  false,
			Message: fmt.Sprintf("Failed to %s data", ctx.Method()),
			Errors:  []string{"Unauthorized"},
			Data:    nil,
		})

		// return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 	"message": "unauthenticated",
		// })
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helper.Response{
			Status:  false,
			Message: fmt.Sprintf("Failed to %s data", ctx.Method()),
			Errors:  []string{"Unauthorized"},
			Data:    nil,
		})
		// return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 	"message": "unauthenticated",
		// })
	}

	ctx.Locals("userid", claims["id"])
	ctx.Locals("useremail", claims["email"])
	ctx.Locals("is_admin", claims["is_admin"])

	// Go to next middleware:
	return ctx.Next()
}

func MiddlewareAuthAdmin(ctx *fiber.Ctx) error {
	isAdmin := ctx.Locals("is_admin") // diasumsikan sudah di-set di middleware JWT
	if isAdmin != true {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helper.Response{
			Status:  false,
			Message: fmt.Sprintf("Failed to %s data", ctx.Method()),
			Errors:  []string{"Unauthorized"},
			Data:    nil,
		})
		// return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		// 	"error": "Access denied. Admins only.",
		// })
	}

	return ctx.Next()
}
