package model

type QueryTrainRequest struct {
	DepartureStationID   int64  `json:"binisIstasyonId"`
	DepartureStationName string `json:"binisIstasyon"`
	ArrivalStationID     int64  `json:"inisIstasyonId"`
	ArrivalStationName   string `json:"inisIstasyonu"`
	DepartureDate        string `json:"gidisTarih"`
}

type SearchTrainRequest struct {
	Request []SearchTrainRequestDetail `json:"request"`
}

type SearchTrainRequestDetail struct {
	RequestID           string              `json:"requestID"`
	DepartureDate       string              `json:"gidisTarih"`
	DepartureStationID  int64               `json:"binisIstasyonId"`
	ArrivalStationID    int64               `json:"inisIstasyonId"`
	ArrivalDate         string              `json:"inisTarih"`
	TourID              int64               `json:"tourID"`
	TrainID             int64               `json:"trainID"`
	Email               string              `json:"email"`
	IsEmailNotification bool                `json:"emailNotification"`
	ExternalInformation ExternalInformation `json:"externalInformation"`
}

type ExternalInformation struct {
	DepartureStation string `json:"departureStation"`
	ArrivalStation   string `json:"arrivalStation"`
	DepartureDate    string `json:"departureDate"`
	ArrivalDate      string `json:"arrivalDate"`
}
