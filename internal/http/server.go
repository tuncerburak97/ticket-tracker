package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"ticket-tracker/internal/http/router"
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

	err := app.Listen(":" + "8080")
	if err != nil {
		return err
	}
	return nil

}
