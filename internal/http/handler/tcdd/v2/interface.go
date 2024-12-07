package v2

import "github.com/gofiber/fiber/v2"

type HandlerV2Interface interface {
	LoadStationsV2(c *fiber.Ctx) error
	QueryTrainV2(c *fiber.Ctx) error
	AddSearchRequestV2(c *fiber.Ctx) error
}
