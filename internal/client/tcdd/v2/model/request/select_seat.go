package request

type SelectSeatRequest struct {
	TrainCarID          int    `json:"trainCarId"`
	FromStationID       int    `json:"fromStationId"`
	ToStationID         int    `json:"toStationId"`
	Gender              string `json:"gender"`
	SeatNumber          string `json:"seatNumber"`
	PassengerTypeID     int    `json:"passengerTypeId"`
	TotalPassengerCount int    `json:"totalPassengerCount"`
	FareFamilyID        int    `json:"fareFamilyId"`
}
