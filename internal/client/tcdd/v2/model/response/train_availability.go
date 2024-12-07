package response

type TrainLeg struct {
	TrainAvailabilities []TrainAvailability `json:"trainAvailabilities"`
}

type TrainAvailability struct {
	Trains        []Train `json:"trains"`
	TotalTripTime int64   `json:"totalTripTime"`
	MinPrice      float64 `json:"minPrice"`
}

type Train struct {
	ID                 int       `json:"id"`
	Number             string    `json:"number"`
	Name               string    `json:"name"`
	CommercialName     string    `json:"commercialName"`
	Type               string    `json:"type"`
	DepartureStationId int       `json:"departureStationId"`
	ArrivalStationId   int       `json:"arrivalStationId"`
	TotalDistance      float64   `json:"totalDistance"`
	Segments           []Segment `json:"segments"`
}

type Segment struct {
	ID            int     `json:"id"`
	DepartureTime int64   `json:"departureTime"`
	ArrivalTime   int64   `json:"arrivalTime"`
	Stops         bool    `json:"stops"`
	Duration      int     `json:"duration"`
	Distance      float64 `json:"distance"`
}

type TrainAvailabilityResponse struct {
	TrainLegs []TrainLeg `json:"trainLegs"`
}
