package service

import (
	"errors"
	"time"

	"ticketing/dto"
	"ticketing/model"
	"ticketing/repository"
)

type TicketService interface {
	PurchaseTicket(userID uint, req dto.TicketRequest) (*dto.TicketResponse, error)
	GetUserTickets(userID uint, page, limit int) ([]dto.TicketResponse, *dto.Pagination, error)
	GetTicketByID(userID, ticketID uint) (*dto.TicketResponse, error)
	CancelTicket(userID, ticketID uint) error
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
}

func NewTicketService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
	}
}

func (s *ticketService) PurchaseTicket(userID uint, req dto.TicketRequest) (*dto.TicketResponse, error) {
	// Check event availability
	event, err := s.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	if event.Status != model.Upcoming {
		return nil, errors.New("event is not available for ticket purchase")
	}

	available, err := s.eventRepo.GetAvailableTickets(event.ID)
	if err != nil {
		return nil, err
	}

	if available <= 0 {
		return nil, errors.New("no available tickets for this event")
	}

	ticket := &model.Ticket{
		EventID:     event.ID,
		UserID:      userID,
		Status:      model.Booked,
		BookingDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	return s.mapTicketToResponse(ticket, event), nil
}

func (s *ticketService) GetUserTickets(userID uint, page, limit int) ([]dto.TicketResponse, *dto.Pagination, error) {
	tickets, total, err := s.ticketRepo.FindAll(page, limit, userID)
	if err != nil {
		return nil, nil, err
	}

	var responses []dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, *s.mapTicketToResponse(&ticket, &ticket.Event))
	}

	pagination := &dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}

	return responses, pagination, nil
}

func (s *ticketService) GetTicketByID(userID, ticketID uint) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if ticket.UserID != userID {
		return nil, errors.New("unauthorized to view this ticket")
	}

	return s.mapTicketToResponse(ticket, &ticket.Event), nil
}

func (s *ticketService) CancelTicket(userID, ticketID uint) error {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	if ticket.UserID != userID {
		return errors.New("unauthorized to cancel this ticket")
	}

	if ticket.Status != model.Booked {
		return errors.New("only booked tickets can be cancelled")
	}

	// Check if event has already started
	event, err := s.eventRepo.FindByID(ticket.EventID)
	if err != nil {
		return err
	}

	eventTime, err := time.Parse("2006-01-02 15:04:05", event.DateTime)
	if err != nil {
		return err
	}

	if time.Now().After(eventTime) {
		return errors.New("cannot cancel ticket for event that has already started")
	}

	return s.ticketRepo.Cancel(ticketID)
}

func (s *ticketService) mapTicketToResponse(ticket *model.Ticket, event *model.Event) *dto.TicketResponse {
	return &dto.TicketResponse{
		ID:          ticket.ID,
		EventName:   event.Name,
		EventDate:   event.DateTime,
		Location:    event.Location,
		Price:       event.Price,
		Status:      string(ticket.Status),
		BookingDate: ticket.BookingDate,
	}
}
