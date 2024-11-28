package tcdd

import (
	"ticket-tracker/internal/api_clients/tcdd/model/request"
	"ticket-tracker/internal/api_clients/tcdd/model/response"
)

type Interface interface {
	LoadAllStation(loadRequest request.StationLoadRequest) (*response.StationLoadResponse, error)
	TripSearch(tripSearchRequest request.TripSearchRequest) (*response.TripSearchResponse, error)
	StationEmptyPlaceSearch(stationEmptyPlaceSearchRequest request.StationEmptyPlaceSearchRequest) (*response.StationEmptyPlaceSearchResponse, error)
	CheckSeat(reserveSeatRequest request.CheckSeatRequest) (*response.CheckSeatResponse, error)
	LocationSelectionWagon(locationSelectionWagonRequest request.LocationSelectionWagonRequest) (*response.LocationSelectionWagonResponse, error)
	ReserveSeat(reserveSeatRequest request.ReserveSeatRequest) (*response.ReserveSeatResponse, error)
}
