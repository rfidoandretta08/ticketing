package service

import (
	"errors"

	"ticketing/dto"
	"ticketing/model"
	"ticketing/repository"
)

type EventService interface {
	CreateEvent(req dto.EventRequest) (*dto.EventResponse, error)
	GetAllEvents(page, limit int, search string) ([]dto.EventResponse, *dto.Pagination, error)
	GetEventByID(id uint) (*dto.EventResponse, error)
	UpdateEvent(id uint, req dto.EventRequest) (*dto.EventResponse, error)
	DeleteEvent(id uint) error
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) CreateEvent(req dto.EventRequest) (*dto.EventResponse, error) {
	event := &model.Event{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		DateTime:    req.DateTime,
		Capacity:    req.Capacity,
		Price:       req.Price,
		Status:      model.Upcoming,
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	// Assuming available tickets is calculated from repository function.
	available, err := s.eventRepo.GetAvailableTickets(event.ID)
	if err != nil {
		return nil, err
	}

	// Now pass the available tickets to mapEventToResponse
	return s.mapEventToResponse(event, available), nil
}

func (s *eventService) GetAllEvents(page, limit int, search string) ([]dto.EventResponse, *dto.Pagination, error) {
	events, total, err := s.eventRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, nil, err
	}

	var responses []dto.EventResponse
	for _, event := range events {
		available, err := s.eventRepo.GetAvailableTickets(event.ID)
		if err != nil {
			return nil, nil, err
		}

		// Pass available tickets to mapEventToResponse
		response := s.mapEventToResponse(&event, available)
		responses = append(responses, *response)
	}

	pagination := &dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}

	return responses, pagination, nil
}

func (s *eventService) GetEventByID(id uint) (*dto.EventResponse, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	available, err := s.eventRepo.GetAvailableTickets(event.ID)
	if err != nil {
		return nil, err
	}

	// Pass available tickets to mapEventToResponse
	return s.mapEventToResponse(event, available), nil
}

func (s *eventService) UpdateEvent(id uint, req dto.EventRequest) (*dto.EventResponse, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if event has started or completed
	if event.Status != model.Upcoming {
		return nil, errors.New("cannot update event that is not upcoming")
	}

	event.Name = req.Name
	event.Description = req.Description
	event.Location = req.Location
	event.DateTime = req.DateTime
	event.Capacity = req.Capacity
	event.Price = req.Price

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	available, err := s.eventRepo.GetAvailableTickets(event.ID)
	if err != nil {
		return nil, err
	}

	// Pass available tickets to mapEventToResponse
	return s.mapEventToResponse(event, available), nil
}

func (s *eventService) DeleteEvent(id uint) error {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if event has tickets
	if len(event.Tickets) > 0 {
		return errors.New("cannot delete event with existing tickets")
	}

	return s.eventRepo.Delete(id)
}

func (s *eventService) mapEventToResponse(event *model.Event, available int) *dto.EventResponse {
	return &dto.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		DateTime:    event.DateTime,
		Capacity:    event.Capacity,
		Available:   available,
		Price:       event.Price,
		Status:      string(event.Status),
	}
}
