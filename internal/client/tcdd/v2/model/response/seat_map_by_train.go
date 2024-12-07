package response

type SeatAllocation struct {
	SeatNumber string   `json:"seatNumber"`
	Genders    []string `json:"genders"`
}

type SeatMapTemplate struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SeatMaps    []SeatMap `json:"seatMaps"`
}

type SeatMap struct {
	Id         int         `json:"id"`
	SeatNumber string      `json:"seatNumber"`
	Item       SeatMapItem `json:"item"`
}

type SeatMapItem struct {
	Name string `json:"name"`
}

type SeatMapObject struct {
	AvailableSeatCount int              `json:"availableSeatCount"`
	TrainCarID         int              `json:"trainCarId"`
	TrainCarIndex      int              `json:"trainCarIndex"`
	AllocationSeats    []SeatAllocation `json:"allocationSeats"`
	SeatMapTemplate    SeatMapTemplate  `json:"seatMapTemplate"`
}

type SeatMapByTrainRequestResponse struct {
	SeatMaps []SeatMapObject `json:"seatMaps"`
}
