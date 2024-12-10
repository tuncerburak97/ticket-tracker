package v2

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
	model "ticket-tracker/internal/controller/dtos/tcdd"
	v2 "ticket-tracker/internal/domain/tcdd/v2"
	"ticket-tracker/pkg/utils/http"
)

type Handler struct {
	s *v2.TcddServiceV2
}

func NewHandler() *Handler {
	return &Handler{s: v2.NewService()}
}

func (h *Handler) LoadStationsV2(c *fiber.Ctx) error {
	stations, err := h.s.LoadAllStationV2()
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.DataResponse(c, stations)
}

func (h *Handler) QueryTrainV2(c *fiber.Ctx) error {
	var req model.QueryTrainRequest
	if err := c.BodyParser(&req); err != nil {
		return http.FailResponse(c, err.Error())
	}

	recipe, err := h.s.QueryTrainV2(&req)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.DataResponse(c, recipe)
}

func (h *Handler) AddSearchRequestV2(c *fiber.Ctx) error {
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
