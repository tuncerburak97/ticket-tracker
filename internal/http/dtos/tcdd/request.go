package model

type QueryTrainRequest struct {
	DepartureStationID   int64  `json:"departureStationID"`
	DepartureStationName string `json:"departureStation"`
	ArrivalStationID     int64  `json:"arrivalStationID"`
	ArrivalStationName   string `json:"arrivalStation"`
	DepartureDate        string `json:"departureDate"`
}

type SearchTrainRequest struct {
	Request []SearchTrainRequestDetail `json:"request"`
}

type SearchTrainRequestDetail struct {
	RequestID           string              `json:"requestID"`
	DepartureDate       string              `json:"departureDate"`
	DepartureStationID  int64               `json:"departureStationID"`
	ArrivalStationID    int64               `json:"arrivalStationID"`
	ArrivalDate         string              `json:"arrivalDate"`
	TourID              int64               `json:"tourID"`
	TrainID             int64               `json:"trainID"`
	Email               string              `json:"email"`
	Gender              string              `json:"gender"`
	IsEmailNotification bool                `json:"emailNotification"`
	ExternalInformation ExternalInformation `json:"externalInformation"`
}

type ExternalInformation struct {
	DepartureStation string `json:"departureStation"`
	ArrivalStation   string `json:"arrivalStation"`
	DepartureDate    string `json:"departureDate"`
	ArrivalDate      string `json:"arrivalDate"`
}
