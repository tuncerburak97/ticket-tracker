package v2

import (
	"github.com/gofiber/fiber/v2"
	v2 "ticket-tracker/internal/http/handler/tcdd/v2"
)

func Router(router fiber.Router) {
	var recipeHandler = v2.NewHandler()
	router.Get("/load", recipeHandler.LoadStationsV2)
	router.Post("/query", recipeHandler.QueryTrainV2)
	router.Post("/add", recipeHandler.AddSearchRequestV2)
}
