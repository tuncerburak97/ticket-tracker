package tcdd

import (
	"github.com/gofiber/fiber/v2"
	"ticket-tracker/internal/http/handler/tcdd"
)

func Router(router fiber.Router) {
	var recipeHandler = tcdd.NewHandler()
	router.Post("/add", recipeHandler.AddSearchRequest)
	router.Get("/load", recipeHandler.LoadStations)
	router.Post("/query", recipeHandler.QueryTrain)
}
