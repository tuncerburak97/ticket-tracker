package v1

import (
	"github.com/gofiber/fiber/v2"
	"ticket-tracker/internal/controller/handler/tcdd/v1"
)

func Router(router fiber.Router) {
	var recipeHandler = v1.NewHandler()
	router.Post("/add", recipeHandler.AddSearchRequest)
	router.Get("/load", recipeHandler.LoadStations)
	router.Post("/query", recipeHandler.QueryTrain)
}
