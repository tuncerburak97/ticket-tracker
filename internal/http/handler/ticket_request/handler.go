package ticket_request

import (
	"github.com/gofiber/fiber/v2"
	"ticket-tracker/internal/domain/ticket_request"
	"ticket-tracker/pkg/utils/http"
)

type Handler struct {
	s *ticket_request.Service
}

func NewHandler() *Handler {
	return &Handler{
		s: ticket_request.GetService(),
	}
}

func (h *Handler) FindRequestById(c *fiber.Ctx) error {
	id := c.Params("id")
	response, err := h.s.FindById(id)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.JsonResponse(c, response)
}

func (h *Handler) FindAllRequest(c *fiber.Ctx) error {
	response, err := h.s.FindAll()
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.JsonResponse(c, response)
}

func (h *Handler) FindRequestByMail(c *fiber.Ctx) error {
	mail := c.Params("mail")
	response, err := h.s.FindByMail(mail)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.JsonResponse(c, response)
}

func (h *Handler) FindRequestByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	response, err := h.s.FindByStatus(status)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.JsonResponse(c, response)
}

func (h *Handler) FindRequestByMailAndStatus(c *fiber.Ctx) error {
	mail := c.Params("mail")
	status := c.Params("status")
	response, err := h.s.FindByMailAndStatus(mail, status)
	if err != nil {
		return http.FailResponse(c, err.Error())
	}
	return http.JsonResponse(c, response)
}
