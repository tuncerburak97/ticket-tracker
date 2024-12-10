package domain

import (
	"gorm.io/gorm"
	"time"
)

type TicketRequest struct {
	ID                  string `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DepartureDate       string `gorm:"column:departure_date"`
	DepartureStationID  int64  `gorm:"column:departure_station_id"`
	DepartureStation    string `gorm:"column:departure_station"`
	ArrivalDate         string `gorm:"column:arrival_date"`
	ArrivalStationID    int64  `gorm:"column:arrival_station_id"`
	ArrivalStation      string `gorm:"column:arrival_station"`
	TrainID             int64  `gorm:"column:train_id"`
	IsEmailNotification bool   `gorm:"column:is_email_notification"`
	Status              string `gorm:"column:status"`
	TotalAttempt        int    `gorm:"column:total_attempt"`
	TourID              int64  `gorm:"column:tour_id"`
	Gender              string `json:"gender"`
	Name                string `json:"name"`
	LastName            string `json:"lastName"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	BirthDate           string `json:"birthDate"`
	IdentityNumber      string `json:"identityNumber"`
}

func (entity TicketRequest) TableName() string {
	return "ticket_request"
}

func (entity *TicketRequest) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	entity.CreatedAt = now
	entity.UpdatedAt = now
	return nil
}
