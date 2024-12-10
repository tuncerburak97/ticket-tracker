package v2

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strings"
	"sync"
	v2 "ticket-tracker/internal/client/tcdd/v2"
	request2 "ticket-tracker/internal/client/tcdd/v2/model/request"
	"ticket-tracker/internal/client/tcdd/v2/model/response"
	apiModel "ticket-tracker/internal/controller/dtos/tcdd"
	"ticket-tracker/internal/domain"
	"ticket-tracker/internal/domain/ticket_request"
	"ticket-tracker/pkg/logger"
	"ticket-tracker/pkg/utils/validation/http"
	"time"
)

type TcddServiceV2 struct {
	tccdClientV2 *v2.HttpClient
	log          *logrus.Logger
	once         sync.Once
	stations     *[]response.StationLoadResponse
}

func NewService() *TcddServiceV2 {
	return &TcddServiceV2{
		tccdClientV2: v2.GetTcddHttpClientInstance(),
		log:          logger.GetLogger(),
	}
}

func (ts *TcddServiceV2) LoadAllStationV2() (*apiModel.StationInformation, error) {
	stations, err := ts.LoadStationsOnce()

	if err != nil {
		return &apiModel.StationInformation{
			Message:  "Error loading stations",
			Success:  false,
			Response: make([]apiModel.LoadStationResponse, 0),
		}, err

	}

	var stationList []apiModel.LoadStationResponse
	for _, station := range *stations {
		responseData := apiModel.LoadStationResponse{
			StationID:         station.Id,
			StationName:       station.Name,
			StationCode:       station.StationCode,
			StationTrainTypes: station.StationTrainTypes,
			StationViewName:   station.Name,
			ToStationList:     getToStationListFromPairId(stations, station.Pairs),
		}
		stationList = append(stationList, responseData)
	}

	return &apiModel.StationInformation{
		Message:  "Stations loaded successfully",
		Success:  true,
		Response: stationList,
	}, nil

}

func (ts *TcddServiceV2) QueryTrainV2(request *apiModel.QueryTrainRequest) (*apiModel.QueryTrainResponse, error) {
	clientTrainAvailabilityRequest := setClientTrainAvailabilityRequest(request)

	searchTrain, err := ts.tccdClientV2.TrainAvailability(clientTrainAvailabilityRequest)
	if err != nil {
		ts.log.Errorf("error [tcdd_service][QueryTrain]: %v\n", err)
		return nil, err
	}

	var queryTrainResponseDetail []apiModel.QueryTrainResponseDetail
	var mu sync.Mutex
	var wg sync.WaitGroup
	var globalError error

	for _, trainLeg := range searchTrain.TrainLegs {
		for _, trainAvailability := range trainLeg.TrainAvailabilities {
			for _, train := range trainAvailability.Trains {
				wg.Add(1)
				go func(train response.Train, trainLeg response.TrainLeg, trainAvailability response.TrainAvailability) {
					defer wg.Done()
					detail, err := ts.processTrain(request, train, trainLeg, trainAvailability)
					if err != nil {
						ts.log.Errorf("error [tcdd_service][QueryTrain]: %v\n", err)
						mu.Lock()
						if globalError == nil {
							globalError = err
						}
						mu.Unlock()
						return
					}

					mu.Lock()
					queryTrainResponseDetail = append(queryTrainResponseDetail, *detail)
					mu.Unlock()
				}(train, trainLeg, trainAvailability)
			}
		}
	}

	wg.Wait()

	if globalError != nil {
		return nil, globalError
	}

	orderByArrivalDate(queryTrainResponseDetail)

	return &apiModel.QueryTrainResponse{
		Details: queryTrainResponseDetail,
	}, nil
}

func (ts *TcddServiceV2) AddSearchRequest(requests *apiModel.SearchTrainRequest) (*apiModel.SearchTrainResponse, error) {

	err := http.ValidateRequest(requests)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	for _, request := range requests.Request {

		if stations, err := ts.LoadAllStationV2(); err != nil {
			return nil, fmt.Errorf("error getting stations: %v", err)
		} else {

			if !checkStationIDIsValid(request.DepartureStationID, stations.Response) || !checkStationIDIsValid(request.ArrivalStationID, stations.Response) {
				return nil, errors.New("invalid arrival or departure station id")
			}
			departureStation, err := GetStationByStationID(stations.Response, request.DepartureStationID)
			if err != nil {
				return nil, fmt.Errorf("error getting departure station: %v", err)
			}
			found := false
			for _, toStation := range departureStation.ToStationList {
				if toStation.ToStationID == request.ArrivalStationID {
					found = true
				}
			}
			if !found {
				return nil, errors.New("arrival station is not reachable from departure station")
			}

			arrivalStation, _ := GetStationByStationID(stations.Response, request.ArrivalStationID)

			err = checkEmailRequestExceedThreshold(request.Email, *requests)
			if err != nil {
				return nil, err
			}

			validIdentityNo, err := ts.tccdClientV2.VerifyIdentityNumber(
				&request2.VerifyIdentityNumberRequest{
					IdentityNumber: request.IdentityNumber,
					Name:           request.Name,
					LastName:       request.LastName,
					BirthDate:      request.BirthDate[:4],
				})
			if err != nil {
				return nil, fmt.Errorf("error verifying identity number: %v", err)
			}
			if !validIdentityNo {
				return nil, errors.New("invalid identity number")
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
				Name:                request.Name,
				LastName:            request.LastName,
				Phone:               request.Phone,
				BirthDate:           request.BirthDate,
				IdentityNumber:      request.IdentityNumber,
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
func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	validationResult := regexp.MustCompile(emailRegex).MatchString(email)
	return validationResult
}
func orderByArrivalDate(details []apiModel.QueryTrainResponseDetail) {
	sort.Slice(details, func(i, j int) bool {
		iTime, errI := time.Parse("2006-01-02 15:04:05", details[i].ArrivalDate)
		jTime, errJ := time.Parse("2006-01-02 15:04:05", details[j].ArrivalDate)

		if errI != nil {
			fmt.Printf("Error parsing ArrivalDate for index %d: %v\n", i, errI)
			return false
		}
		if errJ != nil {
			fmt.Printf("Error parsing ArrivalDate for index %d: %v\n", j, errJ)
			return true
		}
		return iTime.Before(jTime)
	})
}

func checkStationIDIsValid(stationID int64, stations []apiModel.LoadStationResponse) bool {
	for _, station := range stations {
		if station.StationID == stationID {
			return true
		}
	}
	return false

}

func (ts *TcddServiceV2) processTrain(
	request *apiModel.QueryTrainRequest,
	train response.Train,
	trainLeg response.TrainLeg,
	trainAvailability response.TrainAvailability,
) (*apiModel.QueryTrainResponseDetail, error) {
	seatMapByTrainID, err := ts.tccdClientV2.SeatMapByTrain(
		&request2.SeatMapByTrainRequest{
			FromStationID: int64(train.DepartureStationId),
			ToStationId:   int64(train.ArrivalStationId),
			TrainId:       int64(train.ID),
			LegIndex:      0,
		})
	if err != nil {
		return nil, err
	}

	availableSeatCount := findAvailableSeatCount(seatMapByTrainID)
	disablePersonCount := findDisablePersonCount(*seatMapByTrainID)
	departureDate, arrivalDate := findDepartureAndArrivalDateFromTrainsSegment(train)

	apiResponse := &apiModel.QueryTrainResponseDetail{
		TrainID:          int64(train.ID),
		TrainName:        train.CommercialName,
		ArrivalStation:   ts.findStationById(int64(train.ArrivalStationId)),
		DepartureStation: ts.findStationById(int64(train.DepartureStationId)),
		DepartureDate:    departureDate,
		ArrivalDate:      arrivalDate,
		TotalTripTime:    secondsToHoursString(trainAvailability.TotalTripTime),
		MinPrice:         trainAvailability.MinPrice,
		EmptyPlace: apiModel.EmptyPlace{
			DisabledPlaceCount:          int64(disablePersonCount),
			TotalEmptyPlaceCount:        int64(availableSeatCount),
			NormalPeopleEmptyPlaceCount: int64(availableSeatCount) - int64(disablePersonCount),
		},
		DepartureStationID: int64(train.DepartureStationId),
		ArrivalStationID:   int64(train.ArrivalStationId),
	}

	return apiResponse, nil
}

func findDepartureAndArrivalDateFromTrainsSegment(train response.Train) (string, string) {
	var departureDate, arrivalDate string

	// Convert Unix timestamps (milliseconds) to human-readable format
	departureTimeMillis := train.Segments[0].DepartureTime
	arrivalTimeMillis := train.Segments[len(train.Segments)-1].ArrivalTime

	// Convert milliseconds to seconds and parse time
	departureTime := time.UnixMilli(departureTimeMillis)
	arrivalTime := time.UnixMilli(arrivalTimeMillis)

	// Format the time as a human-readable string (e.g., "2006-01-02 15:04:05")
	departureDate = departureTime.Format("2006-01-02 15:04:05")
	arrivalDate = arrivalTime.Format("2006-01-02 15:04:05")

	return departureDate, arrivalDate
}

func (ts *TcddServiceV2) findStationById(id int64) string {
	if ts.stations == nil {
		_, err := ts.LoadAllStationV2()
		if err != nil {
			ts.log.Errorf("error [tcdd_service][findStationById]: %v\n", err)
			return ""
		}
	}
	stations := ts.stations
	for _, station := range *stations {
		if station.Id == id {
			return station.Name
		}
	}
	return ""
}

func findAvailableSeatCount(seatMapResponse *response.SeatMapByTrainRequestResponse) int {
	var availableSeatCount int
	for _, seatMap := range seatMapResponse.SeatMaps {
		availableSeatCount += seatMap.AvailableSeatCount
	}
	return availableSeatCount
}

func setClientTrainAvailabilityRequest(request *apiModel.QueryTrainRequest) *request2.TrainAvailabilityRequest {

	var passengerTypeCounts []request2.PassengerTypeCount
	passengerTypeCount := request2.PassengerTypeCount{
		ID:    0,
		Count: 1,
	}
	passengerTypeCounts = append(passengerTypeCounts, passengerTypeCount)

	return &request2.TrainAvailabilityRequest{
		PassengerTypeCounts: passengerTypeCounts,
		SearchRoutes: []request2.SearchRoute{
			{
				DepartureStationId:   int(request.DepartureStationID),
				DepartureStationName: request.DepartureStationName,
				ArrivalStationId:     int(request.ArrivalStationID),
				ArrivalStationName:   request.ArrivalStationName,
				DepartureDate:        request.DepartureDate,
			},
		},
		SearchReservation: false,
	}
}

func (ts *TcddServiceV2) LoadStationsOnce() (*[]response.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		ts.stations, err = ts.tccdClientV2.LoadAllStations()
	})
	return ts.stations, err
}

func getToStationListFromPairId(stations *[]response.StationLoadResponse, pairs []int64) []apiModel.ToStationList {
	var toStationList []apiModel.ToStationList
	for _, pair := range pairs {
		for _, station := range *stations {
			if station.Id == pair {
				toStationList = append(toStationList, apiModel.ToStationList{
					ToStationID:   station.Id,
					ToStationName: station.Name,
				})
			}
		}
	}
	return toStationList
}

func findDisablePersonCount(response response.SeatMapByTrainRequestResponse) int {
	found := 0

	for _, seatMapObject := range response.SeatMaps {
		for _, seatMap := range seatMapObject.SeatMapTemplate.SeatMaps {
			if strings.Contains(seatMap.SeatNumber, "h") || strings.Contains(seatMap.SeatNumber, "H") {
				existsInAllocation := false

				// AllocationSeats üzerinde SeatNumber kontrolü
				for _, allocation := range seatMapObject.AllocationSeats {
					if allocation.SeatNumber == seatMap.SeatNumber {
						existsInAllocation = true
						break
					}
				}

				// SeatMap üzerinde varsa ama AllocationSeats üzerinde yoksa found'u artır
				if !existsInAllocation {
					found++
				}
			}
		}
	}

	return found
}

func secondsToHoursString(seconds int64) string {
	hours := float64(seconds) / 3600
	return fmt.Sprintf("%.2f Hour", hours)
}

// GetStationByStationID
func GetStationByStationID(stations []apiModel.LoadStationResponse, stationID int64) (*apiModel.LoadStationResponse, error) {
	for _, station := range stations {
		if station.StationID == stationID {
			return &station, nil
		}
	}
	return nil, errors.New("station not found")
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
