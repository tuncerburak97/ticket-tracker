package v2

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"ticket-tracker/internal/client/notification/mail"
	"ticket-tracker/internal/client/notification/mail/model"
	"ticket-tracker/internal/client/tcdd/v1/model/common"
	v2 "ticket-tracker/internal/client/tcdd/v2"
	request2 "ticket-tracker/internal/client/tcdd/v2/model/request"
	"ticket-tracker/internal/client/tcdd/v2/model/response"
	"ticket-tracker/internal/domain"
	"ticket-tracker/internal/domain/ticket_request"
	"ticket-tracker/pkg/logger"
	"time"
)

type TrainScheduler struct {
	tcddClient          *v2.HttpClient
	mailClient          *mail.MailHttpClient
	stations            *[]response.StationLoadResponse
	once                sync.Once
	mu                  sync.Mutex
	isZeroRequestLogged bool
	log                 *logrus.Logger
}

var trainSchedulerInstance *TrainScheduler

func GetTrainSchedulerInstance() *TrainScheduler {
	if trainSchedulerInstance == nil {
		trainSchedulerInstance = NewTrainScheduler(v2.GetTcddHttpClientInstance(),
			mail.GetMailHttpClientInstance())
		trainSchedulerInstance.isZeroRequestLogged = false
	}
	return trainSchedulerInstance

}

func NewTrainScheduler(tcddClient *v2.HttpClient,
	mailClient *mail.MailHttpClient,

) *TrainScheduler {
	return &TrainScheduler{
		tcddClient: tcddClient,
		mailClient: mailClient,
		log:        logger.GetLogger(),
	}
}

func (ts *TrainScheduler) GetStationsOnce() (*[]response.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		ts.stations, err = ts.tcddClient.LoadAllStations()
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
	ts.isZeroRequestLogged = false

	ts.log.Info("Running scheduler with pending requests: ", len(pendingRequests))

	if _, err := ts.GetStationsOnce(); err != nil {
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
	seatMapByTrainID, err := ts.tcddClient.SeatMapByTrain(
		&request2.SeatMapByTrainRequest{
			FromStationID: request.DepartureStationID,
			ToStationId:   request.ArrivalStationID,
			TrainId:       request.TrainID,
			LegIndex:      0,
		})
	if err != nil {
		ts.log.Error("Error getting seat map by train: ", err)
		return ""
	}

	found := false

	for _, seatMap := range seatMapByTrainID.SeatMaps {
		unAllocatedSeats := findUnallocatedSeatsByTrainCar(seatMap.SeatMapTemplate.SeatMaps, seatMap.AllocationSeats)
		if len(unAllocatedSeats) > 0 {
			reservedSeat, createdPnr, err := ts.makeReservation(request, unAllocatedSeats, seatMap)
			if err != nil {
				return ""
			}
			ts.sendReservationMail(
				request.Email,
				request.DepartureDate,
				request.ArrivalDate,
				request.DepartureStation,
				request.ArrivalStation,
				[]common.ReserveSeatDetail{
					{
						SeatNo:       reservedSeat,
						WagonOrderNo: seatMap.TrainCarIndex + 1,
					},
				},
				unAllocatedSeats[0].Item.Name,
				createdPnr,
			)
			found = true
			break // Exit the loop once an unallocated seat is found and handled
		}
	}

	// Call handleNotAllocated if no unallocated seats were found
	if !found {
		ts.handleNotAllocated(request)
		return ""
	}

	return request.ID

}

func (ts *TrainScheduler) handleNotAllocated(entity domain.TicketRequest) {
	var ticketRequestRepository = ticket_request.GetRepository()
	now := time.Now()
	entity.UpdatedAt = now
	totalAttempt := entity.TotalAttempt
	entity.TotalAttempt = totalAttempt + 1
	err := ticketRequestRepository.Update(&entity)
	if err != nil {
		ts.log.Error("Error updating ticket request: ", err)
		return
	}
}

func (ts *TrainScheduler) makeReservation(entity domain.TicketRequest, seatMap []response.SeatMap, seatMapObject response.SeatMapObject) (string, string, error) {
	reservedSeat := seatMap[0].SeatNumber
	selectSeatRequest := request2.SelectSeatRequest{
		TrainCarID:          seatMapObject.TrainCarID,
		FromStationID:       int(entity.DepartureStationID),
		ToStationID:         int(entity.ArrivalStationID),
		Gender:              entity.Gender,
		SeatNumber:          reservedSeat,
		PassengerTypeID:     0,
		FareFamilyID:        0,
		TotalPassengerCount: 1,
	}

	selectSeatResponse, err := ts.tcddClient.SelectSeat(&selectSeatRequest)
	if err != nil {
		ts.log.Error("Error selecting seat: ", err)
		return "", "", err
	}

	Tr := request2.SeatSelection{
		TrainCarId:             int64(seatMapObject.TrainCarID),
		SeatNumber:             reservedSeat,
		FromStationId:          int64(int(entity.DepartureStationID)),
		ToStationId:            int64(int(entity.ArrivalStationID)),
		LockForDate:            int64(selectSeatResponse.LockFor),
		SelectedBookingClassId: 1,
		SelectedFareFamilyId:   1,
		SelectedCabinClassId:   2,
	}

	passenger := request2.Passenger{
		Name:                            entity.Name,
		LastName:                        entity.LastName,
		BirthDate:                       entity.BirthDate,
		Contact:                         false,
		PhoneCountryCode:                "90",
		PhoneAreaCode:                   entity.Phone[0:3],
		PhoneNumber:                     entity.Phone[3:],
		Email:                           entity.Email,
		CountryId:                       1,
		IdentityNumber:                  entity.IdentityNumber,
		LoyaltyNumber:                   "",
		GDPR:                            false,
		PassengerMultiLegSeatSelections: [][]request2.SeatSelection{{Tr}},
	}

	createPnrResponse, err := ts.tcddClient.CreatePnr(&request2.CreatePnrRequest{
		AllocationId:   selectSeatResponse.AllocationID,
		PreReservation: false,
		Passengers:     []request2.Passenger{passenger},
	})

	if err != nil {
		ts.log.Error("Error creating PNR: ", err)
		return reservedSeat, "", err

	}
	return reservedSeat, createPnrResponse.Locator, nil

}

func (ts *TrainScheduler) sendReservationMail(recipient string,
	departureDate string,
	arrivalDate string,
	departureStation string,
	arrivalStation string,
	reservedSeats []common.ReserveSeatDetail,
	seatItem string,
	pnrNumber string) {

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
    <th>Kalkış Zamanı</th>
    <th>Varış Zamanı</th>
    <th>Kalkış İstasyonu</th>
    <th>Varış İstasyonu</th>
	<th>Koltuk Tipi</th>
	<th>Rezervasyon Numarası</th>
  </tr>
  <tr>
    <td>%s</td>
    <td>%s</td>
    <td>%s</td>
    <td>%s</td>
	<td>%s</td>
	<td>%s</td>
  </tr>
</table>

<div class="margin-top">
<p>Sizin için aşağıdaki koltuklar rezerv edilmiştir. 10 dakika boyunca koltuk diğer kullanıcılar için görünür olmayacaktır. Bu maili aldıktan itibaren 10 dakika içinde https://ebilet.tcddtasimacilik.gov.tr/ adresinden yukarıdaki Rezervasyon numarası ile Biletlerim sekmesinden ödeme adımına geçebilirsiniz.</p>
<table>
  <tr>
    <th>Vagon No</th>
    <th>Koltuk No</th>
  </tr>
`, departureDate, arrivalDate, departureStation, arrivalStation, seatItem, pnrNumber)

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

func findUnallocatedSeatsByTrainCar(seatMaps []response.SeatMap, allocations []response.SeatAllocation) []response.SeatMap {
	unallocatedSeats := []response.SeatMap{}

	// SeatAllocation'larda bulunan koltuk numaralarını bir set olarak tutuyoruz.
	allocatedSeatNumbers := make(map[string]struct{})
	for _, allocation := range allocations {
		allocatedSeatNumbers[allocation.SeatNumber] = struct{}{}
	}

	for _, seatMap := range seatMaps {
		if _, exists := allocatedSeatNumbers[seatMap.SeatNumber]; !exists {
			if seatMap.SeatNumber == "" {
				continue
			}
			if strings.Contains(seatMap.SeatNumber, "h") || strings.Contains(seatMap.SeatNumber, "H") {
				continue
			}
			unallocatedSeats = append(unallocatedSeats, seatMap)
		}
	}
	return unallocatedSeats
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
