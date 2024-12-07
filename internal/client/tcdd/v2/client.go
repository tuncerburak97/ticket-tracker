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
var authHeader = "Basic eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJlVFFicDhDMmpiakp1cnUzQVk2a0ZnV196U29MQXZIMmJ5bTJ2OUg5THhRIn0.eyJleHAiOjE3MjEzODQ0NzAsImlhdCI6MTcyMTM4NDQxMCwianRpIjoiYWFlNjVkNzgtNmRkZS00ZGY4LWEwZWYtYjRkNzZiYjZlODNjIiwiaXNzIjoiaHR0cDovL3l0cC1wcm9kLW1hc3RlcjEudGNkZHRhc2ltYWNpbGlrLmdvdi50cjo4MDgwL3JlYWxtcy9tYXN0ZXIiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiMDAzNDI3MmMtNTc2Yi00OTBlLWJhOTgtNTFkMzc1NWNhYjA3IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoidG1zIiwic2Vzc2lvbl9zdGF0ZSI6IjAwYzM4NTJiLTg1YjEtNDMxNS04OGIwLWQ0MWMxMTcyYzA0MSIsImFjciI6IjEiLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1tYXN0ZXIiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjAwYzM4NTJiLTg1YjEtNDMxNS04OGIwLWQ0MWMxMTcyYzA0MSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoid2ViIiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIn0.AIW_4Qws2wfwxyVg8dgHRT9jB3qNavob2C4mEQIQGl3urzW2jALPx-e51ZwHUb-TXB-X2RPHakonxKnWG6tDIP5aKhiidzXDcr6pDDoYU5DnQhMg1kywyOaMXsjLFjuYN5PAyGUMh6YSOVsg1PzNh-5GrJF44pS47JnB9zk03Pr08napjsZPoRB-5N4GQ49cnx7ePC82Y7YIc-gTew2baqKQPz9_v381Gbm2V38PZDH9KldlcWut7kqQYJFMJ7dkM_entPJn9lFk7R5h5j_06OlQEpWRMQTn9SQ1AYxxmZxBu5XYMKDkn4rzIIVCkdTPJNCt5PvjENjClKFeUA1DOg"

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

func (s *HttpClient) LoadAllStations() ([]response.StationLoadResponse, error) {

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

	return stationLoadResponse, nil
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
		"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJlVFFicDhDMmpiakp1cnUzQVk2a0ZnV196U29MQXZIMmJ5bTJ2OUg5THhRIn0.eyJleHAiOjE3MjEzODQ0NzAsImlhdCI6MTcyMTM4NDQxMCwianRpIjoiYWFlNjVkNzgtNmRkZS00ZGY4LWEwZWYtYjRkNzZiYjZlODNjIiwiaXNzIjoiaHR0cDovL3l0cC1wcm9kLW1hc3RlcjEudGNkZHRhc2ltYWNpbGlrLmdvdi50cjo4MDgwL3JlYWxtcy9tYXN0ZXIiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiMDAzNDI3MmMtNTc2Yi00OTBlLWJhOTgtNTFkMzc1NWNhYjA3IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoidG1zIiwic2Vzc2lvbl9zdGF0ZSI6IjAwYzM4NTJiLTg1YjEtNDMxNS04OGIwLWQ0MWMxMTcyYzA0MSIsImFjciI6IjEiLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1tYXN0ZXIiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjAwYzM4NTJiLTg1YjEtNDMxNS04OGIwLWQ0MWMxMTcyYzA0MSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoid2ViIiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIn0.AIW_4Qws2wfwxyVg8dgHRT9jB3qNavob2C4mEQIQGl3urzW2jALPx-e51ZwHUb-TXB-X2RPHakonxKnWG6tDIP5aKhiidzXDcr6pDDoYU5DnQhMg1kywyOaMXsjLFjuYN5PAyGUMh6YSOVsg1PzNh-5GrJF44pS47JnB9zk03Pr08napjsZPoRB-5N4GQ49cnx7ePC82Y7YIc-gTew2baqKQPz9_v381Gbm2V38PZDH9KldlcWut7kqQYJFMJ7dkM_entPJn9lFk7R5h5j_06OlQEpWRMQTn9SQ1AYxxmZxBu5XYMKDkn4rzIIVCkdTPJNCt5PvjENjClKFeUA1DOg", // Burada geçerli token'ı kullanmalısınız.
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
		"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJlVFFicDhDMmpiakp1cnUzQVk2a0ZnV196U29MQXZIMmJ5bTJ2OUg5THhRIn0.eyJleHAiOjE3MjEzODQ0NzAsImlhdCI6MTcyMTM4NDQxMCwianRpIjoiYWFlNjVkNzgtNmRkZS00ZGY4LWEwZWYtYjRkNzZiYjZlODNjIiwiaXNzIjoiaHR0cDovL3l0cC1wcm9kLW1hc3RlcjEudGNkZHRhc2ltYWNpbGlrLmdvdi50cjo4MDgwL3JlYWxtcy9tYXN0ZXIiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiMDAzNDI3MmMtNTc2Yi00OTBlLWJhOTgtNTFkMzc1NWNhYjA3IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoidG1zIiwic2Vzc2lvbl9zdGF0ZSI6IjAwYzM4NTJiLTg1YjEtNDMxNS04OGIwLWQ0MWMxMTcyYzA0MSIsImFjciI6IjEiLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1tYXN0ZXIiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjAwYzM4NTJiLTg1YjEtNDMxNS04OGIwLWQ0MWMxMTcyYzA0MSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoid2ViIiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIn0.AIW_4Qws2wfwxyVg8dgHRT9jB3qNavob2C4mEQIQGl3urzW2jALPx-e51ZwHUb-TXB-X2RPHakonxKnWG6tDIP5aKhiidzXDcr6pDDoYU5DnQhMg1kywyOaMXsjLFjuYN5PAyGUMh6YSOVsg1PzNh-5GrJF44pS47JnB9zk03Pr08napjsZPoRB-5N4GQ49cnx7ePC82Y7YIc-gTew2baqKQPz9_v381Gbm2V38PZDH9KldlcWut7kqQYJFMJ7dkM_entPJn9lFk7R5h5j_06OlQEpWRMQTn9SQ1AYxxmZxBu5XYMKDkn4rzIIVCkdTPJNCt5PvjENjClKFeUA1DOg", // Burada geçerli token'ı kullanmalısınız.
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
