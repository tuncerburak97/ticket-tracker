package tcdd

import (
	"encoding/json"
	"log"
	http2 "net/http"
	"ticket-tracker/internal/client/tcdd/model/request"
	"ticket-tracker/internal/client/tcdd/model/response"
	"ticket-tracker/pkg/http"
)

type HttpClient struct {
}

var httpClientInstance *HttpClient

func GetTcddHttpClientInstance() *HttpClient {
	if httpClientInstance == nil {
		httpClientInstance = NewTcddHttpClient()
	}
	return httpClientInstance
}

func NewTcddHttpClient() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) LoadAllStation(loadRequest request.StationLoadRequest) (*response.StationLoadResponse, error) {

	httpClientInstance := http.GetHttpClientInstance()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/istasyon/istasyonYukle",
		Body:    loadRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var stationLoadResponse response.StationLoadResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][LoadAllStation]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &stationLoadResponse)

	return &stationLoadResponse, nil
}

func (c *HttpClient) TripSearch(tripSearchRequest request.TripSearchRequest) (*response.TripSearchResponse, error) {

	httpClientInstance := http.GetHttpClientInstance()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/sefer/seferSorgula",
		Body:    tripSearchRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var tripSearchResponse response.TripSearchResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][TripSearch]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &tripSearchResponse)
	return &tripSearchResponse, nil
}

func (c *HttpClient) StationEmptyPlaceSearch(stationEmptyPlaceSearchRequest request.StationEmptyPlaceSearchRequest) (*response.StationEmptyPlaceSearchResponse, error) {

	httpClientInstance := http.GetHttpClientInstance()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/vagon/vagonBosYerSorgula",
		Body:    stationEmptyPlaceSearchRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var stationEmptyPlaceSearchResponse response.StationEmptyPlaceSearchResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][StationEmptyPlaceSearch]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &stationEmptyPlaceSearchResponse)
	return &stationEmptyPlaceSearchResponse, nil
}

func (c *HttpClient) CheckSeat(reserveSeatRequest request.CheckSeatRequest) (*response.CheckSeatResponse, error) {
	httpClientInstance := http.GetHttpClientInstance()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/koltuk/klCheck",
		Body:    reserveSeatRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var reserveSeatResponse response.CheckSeatResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][CheckSeat]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &reserveSeatResponse)
	return &reserveSeatResponse, nil
}

func (c *HttpClient) LocationSelectionWagon(locationSelectionWagonRequest request.LocationSelectionWagonRequest) (*response.LocationSelectionWagonResponse, error) {
	httpClientInstance := http.GetHttpClientInstance()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/vagon/vagonHaritasindanYerSecimi",
		Body:    locationSelectionWagonRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var locationSelectionWagonResponse response.LocationSelectionWagonResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][LocationSelectionWagon]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &locationSelectionWagonResponse)
	return &locationSelectionWagonResponse, nil

}

func (c *HttpClient) ReserveSeat(reserveSeatRequest request.ReserveSeatRequest) (*response.ReserveSeatResponse, error) {
	httpClientInstance := http.GetHttpClientInstance()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/koltuk/klSec",
		Body:    reserveSeatRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var reserveSeatResponse response.ReserveSeatResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][ReserveSeat]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &reserveSeatResponse)
	return &reserveSeatResponse, nil

}
