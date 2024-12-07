package v1

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
	"ticket-tracker/internal/domain/tcdd/v1"
	v2 "ticket-tracker/internal/domain/tcdd/v2"
	model "ticket-tracker/internal/http/dtos/tcdd"
	"ticket-tracker/pkg/utils/http"
)

type Handler struct {
	s  *v1.TccdService
	s2 *v2.TcddServiceV2
}

func NewHandler() *Handler {
	return &Handler{s: v1.NewService(), s2: v2.NewService()}
}

func (h *Handler) AddSearchRequest(c *fiber.Ctx) error {
	var req model.SearchTrainRequest
	if err := c.BodyParser(&req); err != nil {
		return http.FailResponse(c, err.Error())
	}
	recipe, err := h.s.AddSearchRequest(&req)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.DataResponseCreated(c, recipe)
}

func (h *Handler) LoadStations(c *fiber.Ctx) error {
	stations, err := h.s.LoadStations()
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.DataResponse(c, stations)
}

func (h *Handler) QueryTrain(c *fiber.Ctx) error {
	var req model.QueryTrainRequest
	if err := c.BodyParser(&req); err != nil {
		return http.FailResponse(c, err.Error())
	}

	recipe, err := h.s.QueryTrain(&req)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.DataResponse(c, recipe)

}
