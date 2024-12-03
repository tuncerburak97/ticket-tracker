package ticket_request

import (
	"gorm.io/gorm"
	"ticket-tracker/internal/domain"
	"ticket-tracker/pkg/db"
)

type Repository interface {
	Create(ticketRequest *domain.TicketRequest) error
	FindAll() ([]domain.TicketRequest, error)
	FindById(id string) (*domain.TicketRequest, error)
	FindByStatus(status string) ([]domain.TicketRequest, error)
	Update(ticketRequest *domain.TicketRequest) error
	FindByMail(mail string) ([]domain.TicketRequest, error)
	FindByMailAndStatus(mail, status string) ([]domain.TicketRequest, error)
}

type ticketRequestRepository struct {
	db *gorm.DB
}

var ticketRequestRepositoryInstance Repository

func NewTicketRequestRepository() Repository {
	return &ticketRequestRepository{db.GetDb()}
}

func GetRepository() Repository {

	if ticketRequestRepositoryInstance == nil {
		ticketRequestRepositoryInstance = NewTicketRequestRepository()
	}
	return ticketRequestRepositoryInstance
}

func (r *ticketRequestRepository) Create(ticketRequest *domain.TicketRequest) error {
	return r.db.Create(ticketRequest).Error
}
func (r *ticketRequestRepository) FindAll() ([]domain.TicketRequest, error) {
	var ticketRequests []domain.TicketRequest
	if err := r.db.
		Order("created_at desc").
		Find(&ticketRequests).
		Error; err != nil {
		return nil, err
	}
	return ticketRequests, nil
}

func (r *ticketRequestRepository) FindById(id string) (ticketRequest *domain.TicketRequest, err error) {
	var ticketRequestResponse domain.TicketRequest

	if err := r.db.
		Take(&ticketRequestResponse, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticketRequestResponse, nil
}

func (r *ticketRequestRepository) FindByStatus(status string) ([]domain.TicketRequest, error) {
	var ticketRequests []domain.TicketRequest
	if err := r.db.
		Order("created_at desc").
		Find(&ticketRequests, "status = ?", status).Error; err != nil {
		return nil, err
	}
	return ticketRequests, nil
}

func (r *ticketRequestRepository) Update(ticketRequest *domain.TicketRequest) error {
	return r.db.Save(ticketRequest).Error
}

func (r *ticketRequestRepository) FindByMail(mail string) ([]domain.TicketRequest, error) {
	var ticketRequests []domain.TicketRequest
	if err := r.db.
		Order("created_at desc").
		Find(&ticketRequests, "email = ?", mail).Error; err != nil {
		return nil, err
	}
	return ticketRequests, nil
}

func (r *ticketRequestRepository) FindByMailAndStatus(mail, status string) ([]domain.TicketRequest, error) {
	var ticketRequests []domain.TicketRequest
	if err := r.db.
		Order("created_at desc").
		Find(&ticketRequests, "email = ? AND status = ?", mail, status).Error; err != nil {
		return nil, err
	}
	return ticketRequests, nil
}
