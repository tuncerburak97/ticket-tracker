package ticket_request

import (
	"github.com/gofiber/fiber/v2"
	"ticket-tracker/internal/controller/handler/ticket_request"
)

func Router(router fiber.Router) {
	var ticketRequestHandler = ticket_request.NewHandler()
	router.Get("/:id", ticketRequestHandler.FindRequestById)
	router.Get("", ticketRequestHandler.FindAllRequest)
	router.Get("/mail/:mail", ticketRequestHandler.FindRequestByMail)
	router.Get("/status/:status", ticketRequestHandler.FindRequestByStatus)
	router.Get("/mail/:mail/status/:status", ticketRequestHandler.FindRequestByMailAndStatus)
}
