package request

type Passenger struct {
	Name                            string            `json:"name"`
	LastName                        string            `json:"lastName"`
	BirthDate                       string            `json:"birthDate"`
	Contact                         bool              `json:"contact"`
	PhoneCountryCode                string            `json:"phoneCountryCode"`
	PhoneAreaCode                   string            `json:"phoneAreaCode"`
	PhoneNumber                     string            `json:"phoneNumber"`
	Email                           string            `json:"email"`
	CountryId                       int               `json:"countryId"`
	IdentityNumber                  string            `json:"identityNumber"`
	LoyaltyNumber                   string            `json:"loyaltyNumber"`
	GDPR                            bool              `json:"gdpr"`
	PassengerMultiLegSeatSelections [][]SeatSelection `json:"passengerMultiLegSeatSelections"`
}

type SeatSelection struct {
	TrainCarId             int64  `json:"trainCarId"`
	SeatNumber             string `json:"seatNumber"`
	FromStationId          int64  `json:"fromStationId"`
	ToStationId            int64  `json:"toStationId"`
	LockForDate            int64  `json:"lockForDate"`
	SelectedBookingClassId int    `json:"selectedBookingClassId"`
	SelectedFareFamilyId   int    `json:"selectedFareFamilyId"`
	SelectedCabinClassId   int    `json:"selectedCabinClassId"`
}

type CreatePnrRequest struct {
	AllocationId   string      `json:"allocationId"`
	Passengers     []Passenger `json:"passengers"`
	PreReservation bool        `json:"preReservation"`
}
