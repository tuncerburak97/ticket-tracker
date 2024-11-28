package http

import (
	"github.com/gofiber/fiber/v2"
)

func ResponseWithStatusCode(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}

func JsonResponse(c *fiber.Ctx, data interface{}) error {
	return ResponseWithStatusCode(c, fiber.StatusOK, data)
}

func FailResponse(c *fiber.Ctx, error string) error {
	return ResponseWithStatusCode(c, fiber.StatusBadRequest, Error{
		Message: error,
	})
}

func FailResponseUnauthorized(c *fiber.Ctx, errors ...string) error {
	return ResponseWithStatusCode(c, fiber.StatusUnauthorized, Errors{
		Errors: errors,
	})
}

func FailResponseNotFound(c *fiber.Ctx, errors ...string) error {
	return ResponseWithStatusCode(c, fiber.StatusNotFound, Errors{
		Errors: errors,
	})
}

func DataResponse(c *fiber.Ctx, data interface{}) error {
	return ResponseWithStatusCode(c, fiber.StatusOK, data)
}

func DataResponseCreated(c *fiber.Ctx, data interface{}) error {
	return ResponseWithStatusCode(c, fiber.StatusCreated, data)
}
