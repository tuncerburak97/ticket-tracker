package ticket_request

import "time"

type RetrieveTicketRequest struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DepartureStation string    `json:"departure_station"`
	DepartureDate    string    `json:"departure_date"`
	ArrivalStation   string    `json:"arrival_station"`
	ArrivalDate      string    `json:"arrival_date"`
	Status           string    `json:"status"`
	TotalAttempt     int       `json:"total_attempt"`
	Gender           string    `json:"gender"`
	Email            string    `json:"email"`
}
