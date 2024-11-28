package model

import model "ticket-tracker/internal/http/dtos/tcdd"

type NewTicketRequest struct {
	Id                 string                   `json:"id"`
	SearchTrainRequest model.SearchTrainRequest `json:"search_train_request"`
	DepartureStation   string                   `json:"departure_station"`
	ArrivalStation     string                   `json:"arrival_station"`
}
