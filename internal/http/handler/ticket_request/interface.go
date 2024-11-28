package ticket_request

import "github.com/gofiber/fiber/v2"

type Interface interface {
	FindRequestById(c *fiber.Ctx) error
	FindAllRequest(c *fiber.Ctx) error
	FindRequestByMail(c *fiber.Ctx) error
	FindRequestByStatus(c *fiber.Ctx) error
	FindRequestByMailAndStatus(c *fiber.Ctx) error
}
