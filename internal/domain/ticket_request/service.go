package ticket_request

import apiResponse "ticket-tracker/internal/controller/dtos/ticket_request"

type Service struct {
	ticketRequestRepository Repository
}

type Interface interface {
	FindById(id string) (apiResponse.RetrieveTicketRequest, error)
	FindAll() ([]apiResponse.RetrieveTicketRequest, error)
	FindByMail(mail string) ([]apiResponse.RetrieveTicketRequest, error)
	FindByStatus(status string) ([]apiResponse.RetrieveTicketRequest, error)
	FindByMailAndStatus(status, mail string) ([]apiResponse.RetrieveTicketRequest, error)
}

var ticketRequestService *Service

func NewService() *Service {
	ticketRequestService = &Service{
		ticketRequestRepository: GetRepository(),
	}
	return ticketRequestService
}

func GetService() *Service {
	if ticketRequestService == nil {
		return NewService()
	}
	return ticketRequestService
}

func (service *Service) FindById(id string) (apiResponse.RetrieveTicketRequest, error) {
	entity, err := service.ticketRequestRepository.FindById(id)
	if err != nil {
		return apiResponse.RetrieveTicketRequest{}, err
	}

	dto := apiResponse.RetrieveTicketRequest{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		DepartureStation: entity.DepartureStation,
		DepartureDate:    entity.DepartureDate,
		ArrivalStation:   entity.ArrivalStation,
		ArrivalDate:      entity.ArrivalDate,
		Email:            entity.Email,
		Status:           entity.Status,
		TotalAttempt:     entity.TotalAttempt,
		Gender:           entity.Gender,
	}

	return dto, nil
}

func (service *Service) FindAll() ([]apiResponse.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindAll()
	if err != nil {
		return nil, err
	}

	dtoList := []apiResponse.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := apiResponse.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
			Gender:           entity.Gender,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (service *Service) FindByMail(mail string) ([]apiResponse.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindByMail(mail)
	if err != nil {
		return nil, err
	}

	dtoList := []apiResponse.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := apiResponse.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
			Gender:           entity.Gender,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (service *Service) FindByStatus(status string) ([]apiResponse.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	dtoList := []apiResponse.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := apiResponse.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
			Gender:           entity.Gender,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (service *Service) FindByMailAndStatus(mail, status string) ([]apiResponse.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindByMailAndStatus(mail, status)
	if err != nil {
		return nil, err
	}

	dtoList := []apiResponse.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := apiResponse.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
			Gender:           entity.Gender,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}
