package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	http2 "net/http"
	"ticket-tracker/internal/client/tcdd/v2/model/request"
	"ticket-tracker/internal/client/tcdd/v2/model/response"
	"ticket-tracker/pkg/http"
	"ticket-tracker/pkg/logger"
)

type HttpClient struct {
	log *logrus.Logger
}

var httpClientInstance *HttpClient
var authHeader = "Basic ZGl0cmF2b3llYnNwOmRpdHJhMzQhdm8u"

func GetTcddHttpClientInstance() *HttpClient {
	if httpClientInstance == nil {
		httpClientInstance = NewTcddHttpClient()
	}
	return httpClientInstance
}

func NewTcddHttpClient() *HttpClient {
	return &HttpClient{
		log: logger.GetLogger(),
	}
}

func (s *HttpClient) LoadAllStations() (*[]response.StationLoadResponse, error) {

	httpClientInstance := http.GetHttpClientInstance()
	httpRequest := http.Request{
		Method:  http2.MethodGet,
		URL:     "https://cdn-api-prod-ytp.tcddtasimacilik.gov.tr/datas/station-pairs-INTERNET.json?environment=dev&userId=1",
		Headers: map[string]interface{}{"Unit-Id": "3895", "Content-Type": "application/json", "Authorization": authHeader},
	}

	var stationLoadResponse []response.StationLoadResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		s.log.Errorf("error [tcdd_client][LoadAllStations]: %v\n", err)
		return nil, err
	}

	err = json.Unmarshal(resp, &stationLoadResponse)
	if err != nil {
		s.log.Errorf("error [tcdd_client][LoadAllStations]: %v\n", err)
		return nil, err
	}

	return &stationLoadResponse, nil
}

func (s *HttpClient) TrainAvailability(req *request.TrainAvailabilityRequest) (*response.TrainAvailabilityResponse, error) {
	httpClientInstance := http.GetHttpClientInstance()
	httpRequest := http.Request{
		Method:  http2.MethodPost,
		URL:     "https://gise-api-prod-ytp.tcddtasimacilik.gov.tr/tms/train/train-availability?environment=dev&userId=1",
		Headers: map[string]interface{}{"Unit-Id": "3895", "Content-Type": "application/json", "Authorization": authHeader},
		Body:    req,
	}

	var trainAvailabilityResponse response.TrainAvailabilityResponse
	resp, err := httpClientInstance.SendHttpRequest(httpRequest)
	if err != nil {
		s.log.Errorf("error [tcdd_client][TrainAvailability]: %v\n", err)
		return nil, err
	}

	err = json.Unmarshal(resp, &trainAvailabilityResponse)
	if err != nil {
		s.log.Errorf("error [tcdd_client][TrainAvailability]: %v\n", err)
		return nil, err
	}

	return &trainAvailabilityResponse, nil
}

func (s *HttpClient) SeatMapByTrain(request *request.SeatMapByTrainRequest) (*response.SeatMapByTrainRequestResponse, error) {
	url := "https://gise-api-prod-ytp.tcddtasimacilik.gov.tr/tms/seat-maps/load-by-train-id?environment=dev&userId=1"

	headers := map[string]string{
		"Unit-Id":       "3895",
		"Authorization": authHeader,
	}

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}

	req, err := http2.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")
	// HTTP isteğini gönderiyoruz
	client := &http2.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Yanıtı okuyoruz
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}
	// JSON verisini Response tipine unmarshall ediyoruz
	var responseData response.SeatMapByTrainRequestResponse
	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}

	return &responseData, nil

}

func (s *HttpClient) SelectSeat(request *request.SelectSeatRequest) (*response.SelectSeatResponse, error) {
	url := "https://gise-api-prod-ytp.tcddtasimacilik.gov.tr/tms/inventory/select-seat?environment=dev&userId=1"

	headers := map[string]string{
		"Unit-Id":       "3895",
		"Authorization": authHeader,
	}

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}

	req, err := http2.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")
	// HTTP isteğini gönderiyoruz
	client := &http2.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Yanıtı okuyoruz
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}
	// JSON verisini Response tipine unmarshall ediyoruz
	var responseData response.SelectSeatResponse
	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		s.log.Errorf("error [tcdd_client][SeatMapByTrain]: %v\n", err)
		return nil, err
	}

	return &responseData, nil
}
