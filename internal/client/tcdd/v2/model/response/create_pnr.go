package response

type CreatePnrResponse struct {
	Locator string `json:"locator"`
}

type Price struct {
	Type          *string `json:"type"`
	PriceAmount   float64 `json:"priceAmount"`
	PriceCurrency string  `json:"priceCurrency"`
}

type User struct {
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	EnrollmentNo *string `json:"enrollmentNo"`
	Email        *string `json:"email"`
	Gender       *string `json:"gender"`
}

type Passenger struct {
	PassengerId                           int                     `json:"passengerId"`
	Name                                  string                  `json:"name"`
	LastName                              string                  `json:"lastName"`
	Gender                                string                  `json:"gender"`
	Type                                  PassengerType           `json:"type"`
	BirthDate                             string                  `json:"birthDate"`
	Email                                 string                  `json:"email"`
	IdentityNumber                        string                  `json:"identityNumber"`
	PassportNumber                        *string                 `json:"passportNumber"`
	PhoneCountryCode                      string                  `json:"phoneCountryCode"`
	PhoneAreaCode                         string                  `json:"phoneAreaCode"`
	PhoneNumber                           string                  `json:"phoneNumber"`
	LoyaltyNumber                         string                  `json:"loyaltyNumber"`
	TariffCardNumber                      *string                 `json:"tariffCardNumber"`
	CountryId                             int                     `json:"countryId"`
	Contact                               bool                    `json:"contact"`
	CustomerNumber                        *string                 `json:"customerNumber"`
	DisabilityStatus                      *string                 `json:"disabilityStatus"`
	PassengerMultiLegSeatSelectionDetails [][]SeatSelectionDetail `json:"passengerMultiLegSeatSelectionDetails"`
}

type PassengerType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type SeatSelectionDetail struct {
	SeatNumber           string        `json:"seatNumber"`
	Route                Route         `json:"route"`
	Prices               []Price       `json:"prices"`
	TotalPrice           TotalPrice    `json:"totalPrice"`
	TotalTicketAmount    float64       `json:"totalTicketAmount"`
	TotalSsrFee          float64       `json:"totalSsrFee"`
	TotalPenaltyAmount   float64       `json:"totalPenaltyAmount"`
	TotalDeductionAmount float64       `json:"totalDeductionAmount"`
	BookingClass         BookingClass  `json:"bookingClass"`
	TicketNumber         *string       `json:"ticketNumber"`
	SsrIds               *string       `json:"ssrIds"`
	History              []interface{} `json:"history"`
	PostponedTrain       bool          `json:"postponedTrain"`
	AgencyCombinedTrain  bool          `json:"agencyCombinedTrain"`
}

type Route struct {
	Id                int       `json:"id"`
	Ticket            Ticket    `json:"ticket"`
	Train             TrainData `json:"train"`
	TrainCar          TrainCar  `json:"trainCar"`
	DepartureStation  Station   `json:"departureStation"`
	ArrivalStation    Station   `json:"arrivalStation"`
	DepartureTime     string    `json:"departureTime"`
	ArrivalTime       string    `json:"arrivalTime"`
	DepartureDateTime string    `json:"departureDateTime"`
	ArrivalDateTime   string    `json:"arrivalDateTime"`
	OpenTicketModel   *string   `json:"openTicketModel"`
	Mco               *string   `json:"mco"`
	DiscountPackage   *string   `json:"discountPackage"`
	Subscription      *string   `json:"subscription"`
}

type Ticket struct {
	Status                  *string `json:"status"`
	AlterationSource        *string `json:"alterationSource"`
	CreateTransactionDetail *string `json:"createTransactionDetail"`
	CheckTransactionDetail  *string `json:"checkTransactionDetail"`
	ChangeTransactionDetail *string `json:"changeTransactionDetail"`
	CancelTransactionDetail *string `json:"cancelTransactionDetail"`
	DivideTicketNumber      string  `json:"divideTicketNumber"`
	MasterTicketNumber      *string `json:"masterTicketNumber"`
	Number                  *string `json:"number"`
}

type TrainData struct {
	Id                         int    `json:"id"`
	TrainDepartureIsAnotherDay bool   `json:"trainDepartureIsAnotherDay"`
	DepartureDate              string `json:"departureDate"`
	Name                       string `json:"name"`
	CommercialName             string `json:"commercialName"`
	TrainNumber                string `json:"trainNumber"`
	Type                       string `json:"type"`
	Line                       Line   `json:"line"`
}

type Line struct {
	Id               int     `json:"id"`
	Name             *string `json:"name"`
	DepartureStation Station `json:"departureStation"`
	ArrivalStation   Station `json:"arrivalStation"`
}

type Station struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TrainCar struct {
	Index      int    `json:"index"`
	Name       string `json:"name"`
	TrainCarId int    `json:"trainCarId"`
}

type TotalPrice struct {
	Type          *string `json:"type"`
	PriceAmount   float64 `json:"priceAmount"`
	PriceCurrency string  `json:"priceCurrency"`
}

type BookingClass struct {
	Id         int        `json:"id"`
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	CabinClass CabinClass `json:"cabinClass"`
	FareFamily FareFamily `json:"fareFamily"`
}

type CabinClass struct {
	Id                      int    `json:"id"`
	Code                    string `json:"code"`
	Name                    string `json:"name"`
	ShowAvailabilityOnQuery bool   `json:"showAvailabilityOnQuery"`
}

type FareFamily struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
