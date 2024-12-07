package request

type SearchRoute struct {
	DepartureStationId   int    `json:"departureStationId"`
	DepartureStationName string `json:"departureStationName"`
	ArrivalStationId     int    `json:"arrivalStationId"`
	ArrivalStationName   string `json:"arrivalStationName"`
	DepartureDate        string `json:"departureDate"`
}

type PassengerTypeCount struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type TrainAvailabilityRequest struct {
	SearchRoutes        []SearchRoute        `json:"searchRoutes"`
	PassengerTypeCounts []PassengerTypeCount `json:"passengerTypeCounts"`
	SearchReservation   bool                 `json:"searchReservation"`
}
