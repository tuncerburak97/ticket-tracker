package v1

import (
	request2 "ticket-tracker/internal/client/tcdd/v1/model/request"
	response2 "ticket-tracker/internal/client/tcdd/v1/model/response"
)

type Interface interface {
	LoadAllStation(loadRequest request2.StationLoadRequest) (*response2.StationLoadResponse, error)
	TripSearch(tripSearchRequest request2.TripSearchRequest) (*response2.TripSearchResponse, error)
	StationEmptyPlaceSearch(stationEmptyPlaceSearchRequest request2.StationEmptyPlaceSearchRequest) (*response2.StationEmptyPlaceSearchResponse, error)
	CheckSeat(reserveSeatRequest request2.CheckSeatRequest) (*response2.CheckSeatResponse, error)
	LocationSelectionWagon(locationSelectionWagonRequest request2.LocationSelectionWagonRequest) (*response2.LocationSelectionWagonResponse, error)
	ReserveSeat(reserveSeatRequest request2.ReserveSeatRequest) (*response2.ReserveSeatResponse, error)
}
