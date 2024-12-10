package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"ticket-tracker/internal/controller/router"
)

func Init() error {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/", "./static")
	router.RegisterRoutes(app)
	app.Get("/metrics", monitor.New())

	err := app.Listen(":" + "8080")
	if err != nil {
		return err
	}
	return nil

}
