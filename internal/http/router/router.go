package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	tcddRouter "ticket-tracker/internal/http/router/tcdd"
	tickerRequestRouter "ticket-tracker/internal/http/router/ticket_request"
)

func RegisterRoutes(app *fiber.App) {
	//metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	//tcdd
	tcdd := app.Group("/tcdd")
	tcddRouter.Router(tcdd)

	// ticket-request
	ticketRequest := app.Group("/ticket-request")
	tickerRequestRouter.Router(ticketRequest)

	// not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})
}
