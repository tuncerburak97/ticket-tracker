package v2

import (
	"ticket-tracker/internal/client/tcdd/v2/model/request"
	"ticket-tracker/internal/client/tcdd/v2/model/response"
)

type Interface interface {
	LoadAllStations() ([]response.StationLoadResponse, error)
	TrainAvailability(req *request.TrainAvailabilityRequest) (*response.TrainAvailabilityResponse, error)
	SeatMapByTrain(req *request.SeatMapByTrainRequest) (*response.SeatMapByTrainRequestResponse, error)
	SelectSeat(req *request.SelectSeatRequest) (*response.SelectSeatResponse, error)
}
