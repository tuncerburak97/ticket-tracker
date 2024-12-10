package v1

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"sort"
	"sync"
	"ticket-tracker/internal/client/tcdd/v1"
	"ticket-tracker/internal/client/tcdd/v1/model/request"
	"ticket-tracker/internal/client/tcdd/v1/model/response"
	apiModel "ticket-tracker/internal/controller/dtos/tcdd"
	"ticket-tracker/internal/domain"
	"ticket-tracker/internal/domain/ticket_request"
	"time"
)

type TccdService struct {
	tcddClientV1 *v1.HttpClient
	stations     *response.StationLoadResponse
	once         sync.Once
}

type TccdServiceInterface interface {
	GetStations() (*response.StationLoadResponse, error)
	LoadStations() (*apiModel.StationInformation, error)
	AddSearchRequest(request *apiModel.SearchTrainRequest) (*apiModel.SearchTrainResponse, error)
	QueryTrain(request *apiModel.QueryTrainRequest) (*apiModel.QueryTrainResponse, error)
}

func NewService() *TccdService {
	return &TccdService{
		tcddClientV1: v1.GetTcddHttpClientInstance(),
	}
}

func (ts *TccdService) GetStations() (*response.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		stationLoadRequest := request.StationLoadRequest{
			Language:    0,
			ChannelCode: "3",
			Date:        "Nov 10, 2011 12:00:00 AM",
			SalesQuery:  true,
		}
		ts.stations, err = ts.tcddClientV1.LoadAllStation(stationLoadRequest)
	})
	return ts.stations, err
}

func (ts *TccdService) LoadStations() (*apiModel.StationInformation, error) {
	stations, err := ts.GetStations()
	if err != nil {
		return &apiModel.StationInformation{
			Message:  "Error loading stations",
			Success:  false,
			Response: make([]apiModel.LoadStationResponse, 0),
		}, err
	}
	var stationList []apiModel.LoadStationResponse

	for _, station := range stations.StationInformation {
		isYht := false
		for _, stationTrainType := range station.StationTrainTypes {
			if stationTrainType == "YHT" {
				isYht = true
			}
		}
		if !isYht {
			continue
		}

		var toStationList []apiModel.ToStationList
		for _, toStation := range station.ToStationIDs {
			toStationData, _ := GetStationByStationID(stations.StationInformation, toStation)
			if toStationData == nil {
				continue
			}
			toStationList = append(toStationList, apiModel.ToStationList{
				ToStationID:   toStation,
				ToStationName: toStationData.StationName,
			})
			sortToStationListByName(toStationList)
		}
		stationList = append(stationList, apiModel.LoadStationResponse{
			StationID:         station.StationID,
			StationCode:       station.StationCode,
			StationName:       station.StationName,
			StationViewName:   station.StationViewName,
			StationTrainTypes: station.StationTrainTypes,
			ToStationList:     toStationList,
		})

	}

	sortStationsByStationName(stationList)

	return &apiModel.StationInformation{
		Message:  "Stations loaded",
		Success:  true,
		Response: stationList,
	}, nil
}
func sortStationsByStationName(loadStationResponse []apiModel.LoadStationResponse) {
	sort.Slice(loadStationResponse, func(i, j int) bool {
		return loadStationResponse[i].StationName < loadStationResponse[j].StationName
	})
}

func sortToStationListByName(toStationList []apiModel.ToStationList) {
	sort.Slice(toStationList, func(i, j int) bool {
		return toStationList[i].ToStationName < toStationList[j].ToStationName
	})

}

func (ts *TccdService) AddSearchRequest(requests *apiModel.SearchTrainRequest) (*apiModel.SearchTrainResponse, error) {
	for _, request := range requests.Request {
		parsedTime, err := time.Parse("Jan 2, 2006 03:04:05 PM", request.DepartureDate)
		if err != nil {
			return nil, fmt.Errorf("invalid departure date: %v", err)
		}
		var now = time.Now()
		if now.After(parsedTime) {
			return nil, errors.New("past departure date")
		}

		if !validateEmail(request.Email) {
			return nil, errors.New("invalid email format")
		}
		if stations, err := ts.GetStations(); err != nil {
			return nil, fmt.Errorf("error getting stations: %v", err)
		} else {

			if !checkStationIDIsValid(request.DepartureStationID, stations.StationInformation) || !checkStationIDIsValid(request.ArrivalStationID, stations.StationInformation) {
				return nil, errors.New("invalid arrival or departure station id")
			}

			departureStation, err := GetStationByStationID(stations.StationInformation, request.DepartureStationID)
			if err != nil {
				return nil, fmt.Errorf("error getting departure station: %v", err)
			}
			found := false
			for _, toStationID := range departureStation.ToStationIDs {
				if toStationID == request.ArrivalStationID {
					found = true
				}
			}
			if !found {
				return nil, errors.New("arrival station is not reachable from departure station")
			}

			arrivalStation, _ := GetStationByStationID(stations.StationInformation, request.ArrivalStationID)

			err = checkEmailRequestExceedThreshold(request.Email, *requests)
			if err != nil {
				return nil, err
			}

			uuidWithHyphen := uuid.New()
			ticketRequestEntity := domain.TicketRequest{
				ID:                  uuidWithHyphen.String(),
				DepartureDate:       request.DepartureDate,
				DepartureStationID:  request.DepartureStationID,
				DepartureStation:    departureStation.StationName,
				ArrivalDate:         request.ArrivalDate,
				ArrivalStationID:    request.ArrivalStationID,
				ArrivalStation:      arrivalStation.StationName,
				TourID:              request.TourID,
				TrainID:             request.TrainID,
				Email:               request.Email,
				IsEmailNotification: request.IsEmailNotification,
				Gender:              request.Gender,
				Status:              "PENDING",
				TotalAttempt:        0,
			}

			ticketRequestRepository := ticket_request.GetRepository()
			err = ticketRequestRepository.Create(&ticketRequestEntity)
			if err != nil {
				return nil, fmt.Errorf("error creating ticket request: %v", err)
			}
		}
	}
	return &apiModel.SearchTrainResponse{
		Message: "Request added to scheduler",
		Success: true,
	}, nil

}

func (ts *TccdService) QueryTrain(request *apiModel.QueryTrainRequest) (*apiModel.QueryTrainResponse, error) {

	return &apiModel.QueryTrainResponse{
		nil,
	}, nil
}

func orderByArrivalDate(Details []apiModel.QueryTrainResponseDetail) {
	sort.Slice(Details, func(i, j int) bool {
		iTime, _ := time.Parse("Jan 2, 2006 03:04:05 PM", Details[i].ArrivalDate)
		jTime, _ := time.Parse("Jan 2, 2006 03:04:05 PM", Details[j].ArrivalDate)
		return iTime.Before(jTime)
	})
}

func (ts *TccdService) processTripSearchResult(
	wg *sync.WaitGroup,
	detailsChan chan<- apiModel.QueryTrainResponseDetail,
	errChan chan<- error,
	trip response.SearchResult,
	request *apiModel.QueryTrainRequest,
	tripSearchResponse *response.TripSearchResponse,
) {

}

func calculateTotalEmptyPlace(emptyPlaceList []response.EmptyPlace) int {
	totalEmptyPlace := 0
	for _, emptyPlace := range emptyPlaceList {
		totalEmptyPlace += emptyPlace.EmptyPlace
	}
	return totalEmptyPlace
}
func findTrip(search *response.TripSearchResponse, tourID int64) (int64, bool) {
	for _, trip := range search.SearchResult {
		if trip.TourID == tourID {
			if len(trip.WagonTypesEmptyPlace) > 0 {
				return trip.WagonTypesEmptyPlace[0].RemainingDisabledNumber, true
			}
		}
	}
	return 0, false
}

// commons
func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	validationResult := regexp.MustCompile(emailRegex).MatchString(email)
	return validationResult
}
func GetStationByStationID(stations []response.StationInformation, stationID int64) (*response.StationInformation, error) {
	for _, station := range stations {
		if station.StationID == stationID {
			return &station, nil
		}
	}
	return nil, fmt.Errorf("no station found with ID: %v", stationID)
}

func checkStationIDIsValid(stationID int64, stations []response.StationInformation) bool {
	for _, station := range stations {
		if station.StationID == stationID {
			return true
		}
	}
	return false

}

func checkEmailRequestExceedThreshold(email string, requests apiModel.SearchTrainRequest) error {
	foundedCount := 0
	for _, request := range requests.Request {
		if request.Email == email {
			foundedCount++
		}
	}
	if foundedCount > 5 {
		return errors.New("exceed threshold")
	}
	return nil
}
