package main

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/infrastructure/container"
	"fmt"

	rest "backend-evermos/internal/server/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	containerConf := container.InitContainer()
	// defer mysql.CloseDatabaseConnection(containerConf.Mysqldb)

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	rest.HTTPRouteInit(app, containerConf)

	port := fmt.Sprintf("%s:%d", containerConf.Apps.Host, containerConf.Apps.HttpPort)
	if err := app.Listen(port); err != nil {
		helper.Logger(helper.LoggerLevelFatal, "error", err)
	}
}
