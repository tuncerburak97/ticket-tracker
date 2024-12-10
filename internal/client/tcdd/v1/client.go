package v1

import (
	"encoding/json"
	"log"
	http2 "net/http"
	request2 "ticket-tracker/internal/client/tcdd/v1/model/request"
	response2 "ticket-tracker/internal/client/tcdd/v1/model/response"
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

func (c *HttpClient) LoadAllStation(loadRequest request2.StationLoadRequest) (*response2.StationLoadResponse, error) {

	httpClientInstance := http.GetRestClient()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/istasyon/istasyonYukle",
		Body:    loadRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var stationLoadResponse response2.StationLoadResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][LoadAllStation]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &stationLoadResponse)

	return &stationLoadResponse, nil
}

func (c *HttpClient) TripSearch(tripSearchRequest request2.TripSearchRequest) (*response2.TripSearchResponse, error) {

	httpClientInstance := http.GetRestClient()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/sefer/seferSorgula",
		Body:    tripSearchRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var tripSearchResponse response2.TripSearchResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][TripSearch]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &tripSearchResponse)
	return &tripSearchResponse, nil
}

func (c *HttpClient) StationEmptyPlaceSearch(stationEmptyPlaceSearchRequest request2.StationEmptyPlaceSearchRequest) (*response2.StationEmptyPlaceSearchResponse, error) {

	httpClientInstance := http.GetRestClient()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/vagon/vagonBosYerSorgula",
		Body:    stationEmptyPlaceSearchRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var stationEmptyPlaceSearchResponse response2.StationEmptyPlaceSearchResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][StationEmptyPlaceSearch]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &stationEmptyPlaceSearchResponse)
	return &stationEmptyPlaceSearchResponse, nil
}

func (c *HttpClient) CheckSeat(reserveSeatRequest request2.CheckSeatRequest) (*response2.CheckSeatResponse, error) {
	httpClientInstance := http.GetRestClient()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/koltuk/klCheck",
		Body:    reserveSeatRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var reserveSeatResponse response2.CheckSeatResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][CheckSeat]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &reserveSeatResponse)
	return &reserveSeatResponse, nil
}

func (c *HttpClient) LocationSelectionWagon(locationSelectionWagonRequest request2.LocationSelectionWagonRequest) (*response2.LocationSelectionWagonResponse, error) {
	httpClientInstance := http.GetRestClient()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/vagon/vagonHaritasindanYerSecimi",
		Body:    locationSelectionWagonRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var locationSelectionWagonResponse response2.LocationSelectionWagonResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][LocationSelectionWagon]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &locationSelectionWagonResponse)
	return &locationSelectionWagonResponse, nil

}

func (c *HttpClient) ReserveSeat(reserveSeatRequest request2.ReserveSeatRequest) (*response2.ReserveSeatResponse, error) {
	httpClientInstance := http.GetRestClient()

	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://api-yebsp.tcddtasimacilik.gov.tr/koltuk/klSec",
		Body:    reserveSeatRequest,
		Headers: map[string]interface{}{"Content-Type": "application/json", "Authorization": "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"},
	}
	var reserveSeatResponse response2.ReserveSeatResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		log.Printf("error [tcdd_client][ReserveSeat]: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(resp, &reserveSeatResponse)
	return &reserveSeatResponse, nil

}
