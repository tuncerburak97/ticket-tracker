package model

type QueryTrainResponse struct {
	Details []QueryTrainResponseDetail `json:"details"`
}
type QueryTrainResponseDetail struct {
	TrainID            int64      `json:"trainID"`
	TrainName          string     `json:"trainName"`
	DepartureDate      string     `json:"departureDate"`
	ArrivalDate        string     `json:"arrivalDate"`
	EmptyPlace         EmptyPlace `json:"emptyPlace"`
	ArrivalStation     string     `json:"arrivalStation"`
	DepartureStation   string     `json:"departureStation"`
	DepartureStationID int64      `json:"departureStationID"`
	ArrivalStationID   int64      `json:"arrivalStationID"`
	TotalTripTime      string     `json:"totalTripTime"`
	MinPrice           float64    `json:"minPrice"`
}

type EmptyPlace struct {
	DisabledPlaceCount          int64 `json:"disabledPlaceCount"`
	TotalEmptyPlaceCount        int64 `json:"totalEmptyPlaceCount"`
	NormalPeopleEmptyPlaceCount int64 `json:"normalPeopleEmptyPlaceCount"`
}

type SearchTrainResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type StationInformation struct {
	Response []LoadStationResponse `json:"response"`
	Message  string                `json:"message"`
	Success  bool                  `json:"success"`
}

type LoadStationResponse struct {
	StationName       string          `json:"stationName"`
	StationID         int64           `json:"stationID"`
	StationCode       string          `json:"stationCode"`
	StationTrainTypes []string        `json:"stationTrainTypes"`
	StationViewName   string          `json:"stationViewName"`
	ToStationList     []ToStationList `json:"toStationList"`
}

type ToStationList struct {
	ToStationID   int64  `json:"toStationId"`
	ToStationName string `json:"toStationName"`
}
