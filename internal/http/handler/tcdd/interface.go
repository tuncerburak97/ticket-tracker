package tcdd

import "github.com/gofiber/fiber/v2"

type HandlerInterface interface {
	AddSearchRequest(c *fiber.Ctx) error
	LoadStations(c *fiber.Ctx) error
	QueryTrain(c *fiber.Ctx) error
}
