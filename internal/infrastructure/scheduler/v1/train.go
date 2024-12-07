package v1

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"ticket-tracker/internal/client/notification/mail"
	"ticket-tracker/internal/client/notification/mail/model"
	"ticket-tracker/internal/client/tcdd/v1"
	"ticket-tracker/internal/client/tcdd/v1/model/common"
	request2 "ticket-tracker/internal/client/tcdd/v1/model/request"
	response2 "ticket-tracker/internal/client/tcdd/v1/model/response"
	"ticket-tracker/internal/domain"
	"ticket-tracker/internal/domain/ticket_request"
	"ticket-tracker/pkg/logger"
	"time"
)

type TrainScheduler struct {
	tcddClient          *v1.HttpClient
	mailClient          *mail.MailHttpClient
	stations            *response2.StationLoadResponse
	once                sync.Once
	mu                  sync.Mutex
	isZeroRequestLogged bool
	log                 *logrus.Logger
}

var trainSchedulerInstance *TrainScheduler

func GetTrainSchedulerInstance() *TrainScheduler {
	if trainSchedulerInstance == nil {
		trainSchedulerInstance = NewTrainScheduler(v1.GetTcddHttpClientInstance(),
			mail.GetMailHttpClientInstance())
		trainSchedulerInstance.isZeroRequestLogged = false
	}
	return trainSchedulerInstance

}

func NewTrainScheduler(tcddClient *v1.HttpClient,
	mailClient *mail.MailHttpClient,

) *TrainScheduler {
	return &TrainScheduler{
		tcddClient: tcddClient,
		mailClient: mailClient,
		log:        logger.GetLogger(),
	}
}

func (ts *TrainScheduler) getStations() (*response2.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		stationLoadRequest := request2.StationLoadRequest{
			Language:    0,
			ChannelCode: "3",
			Date:        "Nov 10, 2011 12:00:00 AM",
			SalesQuery:  true,
		}
		ts.stations, err = ts.tcddClient.LoadAllStation(stationLoadRequest)
	})
	return ts.stations, err
}

func (ts *TrainScheduler) Run() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var ticketRequestRepository = ticket_request.GetRepository()
	var pendingRequests, err = ticketRequestRepository.FindByStatus("PENDING")
	if err != nil {
		ts.log.Error("Error getting pending requests: ", err)
		return
	}

	if len(pendingRequests) == 0 && !ts.isZeroRequestLogged {
		ts.log.Info("No pending requests found")
		ts.isZeroRequestLogged = true
		return
	}

	if len(pendingRequests) == 0 {
		return
	}

	ts.log.Info("Running scheduler with pending requests: ", len(pendingRequests))

	if _, err := ts.getStations(); err != nil {
		ts.log.Error("Error getting stations: ", err)
		return
	}
	var foundedRequestIDList = make([]string, 0)

	for _, searchTrainRequest := range pendingRequests {
		foundedRequestIDS := ts.processRequest(searchTrainRequest)
		if foundedRequestIDS != "" {
			foundedRequestIDList = append(foundedRequestIDList, foundedRequestIDS)
		}
	}

	var filterFoundedRequests = make([]domain.TicketRequest, 0)
	for _, foundedRequestID := range foundedRequestIDList {
		for _, request := range pendingRequests {
			if request.ID == foundedRequestID {
				filterFoundedRequests = append(filterFoundedRequests, request)
			}
		}
	}

	ts.UpdateTicketRequestStatusToFound(filterFoundedRequests)
}

func (ts *TrainScheduler) processRequest(request domain.TicketRequest) (requestID string) {

	criteria := request2.Criteria{
		SalesChannel:       3,
		DepartureStation:   request.DepartureStation,
		IsMapDeparture:     false,
		ArrivalStation:     request.ArrivalStation,
		IsMapArrival:       false,
		DepartureDate:      request.DepartureDate,
		IsRegional:         false,
		OperationType:      0,
		PassengerCount:     1,
		IsTransfer:         true,
		DepartureStationID: request.DepartureStationID,
		ArrivalStationID:   request.ArrivalStationID,
		TravelType:         1,
	}

	search, err := ts.tcddClient.TripSearch(request2.TripSearchRequest{
		ChannelCode: 3,
		Language:    0,
		Criteria:    criteria,
	})
	if err != nil {
		ts.log.Error("Error searching trip: ", err)
		return
	}

	b := search.TripSearchResponseInfo.ResponseCode != "000"
	if b {
		ts.log.Error("Error searching trip: ", search.TripSearchResponseInfo.ResponseMsg)
		return
	}
	remainingDisabledNumber, found := ts.findTrip(search, request.TourID)
	if found {
		var requestId = ts.handleFoundTrip(request, int(remainingDisabledNumber), search.SearchResult[0].ArrivalDate)
		if requestId != "" {
			return requestId
		}
		return ts.handleNotFoundTrip(request)

	}

	return ts.handleNotFoundTrip(request)

}
func (ts *TrainScheduler) handleFoundTrip(request domain.TicketRequest, remainingDisabledNumber int, arrivalDate string) (requestID string) {

	placeSearch, err := ts.tcddClient.StationEmptyPlaceSearch(request2.StationEmptyPlaceSearchRequest{
		ChannelCode:   "3",
		Language:      0,
		TourTitleID:   request.TourID,
		DepartureStID: request.DepartureStationID,
		ArrivalStID:   int(request.ArrivalStationID),
	})
	if err != nil {
		ts.log.Error("Error searching empty place: ", err)
		return
	}

	totalEmptyPlace := calculateTotalEmptyPlace(placeSearch.EmptyPlaceList)
	availablePlace := totalEmptyPlace - remainingDisabledNumber
	if availablePlace > 0 {

		locationSelectionWagonRequestList := getLocationSelectionWagonRequestList(placeSearch.EmptyPlaceList, request)
		reservedSeats := ts.reserveSeat(locationSelectionWagonRequestList, request)

		departureValidation := true
		departureDateFormat, err := time.Parse("Jan 02, 2006 03:04:05 PM", request.DepartureDate)
		if err != nil {
			ts.log.Error("Departure Date parse edilemedi:", err)
			departureValidation = false
		}

		arrivalValidation := true
		arrivalDateFormat, err := time.Parse("Jan 02, 2006 03:04:05 PM", request.ArrivalDate)
		if err != nil {
			ts.log.Error("Arrival Date parse edilemedi:", err)
			arrivalValidation = false
		}

		// Türkçe tarih formatını oluşturma ve yazdırma
		var departureDateStr string
		if departureValidation {
			departureDateStr = formatTurkishDate(departureDateFormat)
		} else {
			departureDateStr = request.DepartureDate
		}

		var arrivalDateStr string
		if arrivalValidation {
			arrivalDateStr = formatTurkishDate(arrivalDateFormat)
		} else {
			arrivalDateStr = request.ArrivalDate
		}

		ts.log.WithFields(logrus.Fields{
			"ID":    request.ID,
			"Email": request.Email,
			"Date":  request.DepartureDate,
			"From":  request.DepartureStation,
			"To":    request.ArrivalStation,
		}).Info("Found trip for request")

		ts.sendEmail(
			request.Email,
			availablePlace,
			departureDateStr,
			arrivalDateStr,
			request.DepartureStation,
			request.ArrivalStation,
			reservedSeats)
		return request.ID
	}

	return ""
}

func (ts *TrainScheduler) reserveSeat(
	locationSelectionWagonRequestList []request2.LocationSelectionWagonRequest,
	request domain.TicketRequest,
) []common.ReserveSeatDetail {

	reservedSeats := make([]common.ReserveSeatDetail, 0)
	totalReservedSeat := 0

	for _, locationSelectionWagonRequest := range locationSelectionWagonRequestList {
		seats := ts.processWagonRequest(locationSelectionWagonRequest, request, &totalReservedSeat)
		if seats != nil {
			reservedSeats = append(reservedSeats, seats...)
		}
	}

	return reservedSeats
}

func (ts *TrainScheduler) processWagonRequest(
	locationSelectionWagonRequest request2.LocationSelectionWagonRequest,
	request domain.TicketRequest,
	totalReservedSeat *int,
) []common.ReserveSeatDetail {

	reservedSeats := make([]common.ReserveSeatDetail, 0)

	locationSelectionWagonResponse, err := ts.tcddClient.LocationSelectionWagon(locationSelectionWagonRequest)
	if err != nil {
		ts.log.Error("Error selecting wagon: ", err)
		return nil
	}
	if locationSelectionWagonResponse.ResponseInfo.ResponseCode != "000" {
		ts.log.Error("Error selecting wagon: ", locationSelectionWagonResponse.ResponseInfo.ResponseMsg)
		return nil
	}
	for _, locationSelectionWagon := range locationSelectionWagonResponse.LocationSelectionWagonResponseData.SeatInformationList {
		if locationSelectionWagon.Status == 0 {
			if *totalReservedSeat >= 3 {
				break
			}

			checkSeatRequest := request2.CheckSeatRequest{
				ChannelCode:             "3",
				Language:                0,
				SelectedSeatWagonNumber: locationSelectionWagon.WagonOrderNo,
				SelectedSeatNumber:      locationSelectionWagon.SeatNo,
				TourId:                  strconv.FormatInt(request.TourID, 10),
			}
			checkSeatResponse, err := ts.tcddClient.CheckSeat(checkSeatRequest)
			if err != nil {
				ts.log.Error("Error reserving seat: ", err)
				return nil
			}
			if checkSeatResponse.ResponseInfo.ResponseCode != "000" {
				ts.log.Error("Error reserving seat: ", checkSeatResponse.ResponseInfo.ResponseMsg)
				return nil
			}

			reserveSeatRequest := request2.ReserveSeatRequest{
				ChannelCode:        "3",
				Language:           0,
				TourID:             int(request.TourID),
				WagonOrder:         locationSelectionWagon.WagonOrderNo,
				SeatNo:             locationSelectionWagon.SeatNo,
				Gender:             "M",
				ArrivalStationID:   int(request.ArrivalStationID),
				DepartureStationID: int(request.DepartureStationID),
				Minute:             10,
				Huawei:             false,
			}

			reserveSeatResponse, err := ts.tcddClient.ReserveSeat(reserveSeatRequest)
			if err != nil {
				ts.log.Error("Error reserving seat: ", err)
				return nil
			}
			if reserveSeatResponse.ResponseInfo.ResponseCode != "000" {
				ts.log.Error("Error reserving seat: ", reserveSeatResponse.ResponseInfo.ResponseMsg)
				return nil
			}

			ts.log.WithFields(logrus.Fields{
				"ID":    request.ID,
				"Email": request.Email,
				"Date":  request.DepartureDate,
				"From":  request.DepartureStation,
				"To":    request.ArrivalStation,
			}).Info("Seat reserved for request")

			reservedSeats = append(reservedSeats, common.ReserveSeatDetail{
				SeatNo:       locationSelectionWagon.SeatNo,
				WagonOrderNo: locationSelectionWagon.WagonOrderNo,
			})
			*totalReservedSeat++
		}
	}

	return reservedSeats
}

func (ts *TrainScheduler) findTrip(search *response2.TripSearchResponse, tourID int64) (int64, bool) {
	for _, trip := range search.SearchResult {
		if trip.TourID == tourID {
			if len(trip.WagonTypesEmptyPlace) > 0 {
				return trip.WagonTypesEmptyPlace[0].RemainingDisabledNumber, true
			}
		}
	}
	return 0, false
}

func calculateTotalEmptyPlace(emptyPlaceList []response2.EmptyPlace) int {
	totalEmptyPlace := 0
	for _, emptyPlace := range emptyPlaceList {
		totalEmptyPlace += emptyPlace.EmptyPlace
	}
	return totalEmptyPlace
}

func (ts *TrainScheduler) sendEmail(recipient string,
	availablePlace int,
	departureDate string,
	arrivalDate string,
	departureStation string,
	arrivalStation string,
	reservedSeats []common.ReserveSeatDetail,
) {

	{
		body := fmt.Sprintf(`
<html>
<head>
<style>
table {
  font-family: Arial, sans-serif;
  border-collapse: collapse;
  width: 100%%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}

.margin-top {
  margin-top: 20px;
}
</style>
</head>
<body>
<p>Merhaba,</p>
<p>Aradığınız trende boş yer bulundu. &#128522;</p>
<table>
  <tr>
    <th>Kalan Boş Yer Sayısı</th>
    <th>Kalkış Zamanı</th>
    <th>Varış Zamanı</th>
    <th>Kalkış İstasyonu</th>
    <th>Varış İstasyonu</th>
  </tr>
  <tr>
    <td>%d</td>
    <td>%s</td>
    <td>%s</td>
    <td>%s</td>
    <td>%s</td>
  </tr>
</table>

<div class="margin-top">
<p>Sizin için aşağıdaki koltuklar rezerv edilmiştir. 10 dakika boyunca koltuk diğer kullanıcılar için görünür olmayacaktır. Bu maili aldıktan 10 dakika sonra koltuk kilidi kalkmış olacaktır. İlgili koltuğu 10 dakika sonra kontrol edebilirsiniz!</p>
<table>
  <tr>
    <th>Vagon No</th>
    <th>Koltuk No</th>
  </tr>
`, availablePlace, departureDate, arrivalDate, departureStation, arrivalStation)

		for _, seat := range reservedSeats {

			if strings.Contains(seat.SeatNo, "h") {
				continue
			}

			body += fmt.Sprintf(`
  <tr>
    <td>%d</td>
    <td>%s</td>
  </tr>
`, seat.WagonOrderNo, seat.SeatNo)
		}

		body += `
</table>
</div>

<p>Tekrardan bu yolculuga dair bildirimleri takip etmek isterseniz uygulama üzerinden aynı talebi oluşturabilirsiniz</p>
<p>İyi yolculuklar dileriz!</p>
</body>
</html>`

		email := model.Email{
			To:      recipient,
			Subject: "Tren Bilet Uyarısı",
			Body:    body,
		}

		// Send the email
		err := trainSchedulerInstance.mailClient.SendEmail(email)
		if err != nil {
			ts.log.Error("Error sending email: ", err)
		}

	}
}

func (ts *TrainScheduler) UpdateTicketRequestStatusToFound(foundedRequests []domain.TicketRequest) {

	var ticketRequestRepository = ticket_request.GetRepository()
	now := time.Now()
	for _, request := range foundedRequests {
		request.Status = "FOUND"
		request.UpdatedAt = now
		totalAttempt := request.TotalAttempt
		request.TotalAttempt = totalAttempt + 1
		err := ticketRequestRepository.Update(&request)
		if err != nil {
			ts.log.Error("Error updating ticket request: ", err)
		}
	}
}

func getLocationSelectionWagonRequestList(emptyPlaceList []response2.EmptyPlace, request domain.TicketRequest) []request2.LocationSelectionWagonRequest {
	response := make([]request2.LocationSelectionWagonRequest, 0)
	for _, emptyPlace := range emptyPlaceList {
		if emptyPlace.EmptyPlace > 0 {
			response = append(response, request2.LocationSelectionWagonRequest{
				ChannelCode:          "3",
				Language:             0,
				TourTitleID:          strconv.FormatInt(request.TourID, 10),
				WagonOrderNo:         emptyPlace.WagonOrderNo,
				DepartureStationName: request.DepartureStation,
				ArrivalStationName:   request.ArrivalStation,
			})
		}
	}
	return response
}

func formatTurkishDate(t time.Time) string {
	// Ay isimlerini Türkçe karşılıklarıyla değiştirin
	months := map[string]string{
		"January":   "Ocak",
		"February":  "Şubat",
		"March":     "Mart",
		"April":     "Nisan",
		"May":       "Mayıs",
		"June":      "Haziran",
		"July":      "Temmuz",
		"August":    "Ağustos",
		"September": "Eylül",
		"October":   "Ekim",
		"November":  "Kasım",
		"December":  "Aralık",
	}

	// Günü, ayı, yılı, saati ve dakikayı formatlayın
	day := t.Day()
	month := months[t.Month().String()]
	year := t.Year()
	hour := t.Hour()
	minute := t.Minute()

	// Türkçe formatta string oluşturma
	return fmt.Sprintf("%02d-%s-%d %02d:%02d", day, month, year, hour, minute)
}

func (ts *TrainScheduler) handleNotFoundTrip(request domain.TicketRequest) (requestID string) {

	/*
		log.Printf("Trip not found for request: %s and email: %s date: %s from: %s to: %s",
			request.ID,
			request.Email,
			request.DepartureDate,
			request.DepartureStation,
			request.ArrivalStation)


	*/
	totalAttempt := request.TotalAttempt
	request.TotalAttempt = totalAttempt + 1
	now := time.Now()
	request.UpdatedAt = now

	var ticketRequestRepository = ticket_request.GetRepository()
	err := ticketRequestRepository.Update(&request)
	if err != nil {
		ts.log.Error("Error updating ticket request: ", err)
		return ""
	}

	return ""
}
