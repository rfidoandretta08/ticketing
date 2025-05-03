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
	GetAllTickets(page, limit int) ([]dto.TicketResponse, *dto.Pagination, error)
	GetUserTickets(userID uint, page, limit int) ([]dto.TicketResponse, *dto.Pagination, error)
	GetTicketByID(userID, ticketID uint) (*dto.TicketResponse, error)
	CancelTicket(userID, ticketID uint) error
	UpdatePayment(userID, ticketID uint) (*dto.PaymentUpdateResponse, error)
	CancelPayment(userID, ticketID uint) (*dto.PaymentUpdateResponse, error)
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
	// Cek ketersediaan event
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

	// Hitung SubTotal berdasarkan harga tiket dan jumlah (qty)
	subTotal := event.Price * float64(req.Qty)

	// Buat tiket baru
	ticket := &model.Ticket{
		EventID:     event.ID,
		UserID:      userID,
		Status:      model.Available, // Status awal Available
		Qty:         req.Qty,
		SubTotal:    subTotal,
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
		Qty:         ticket.Qty,      // Menyertakan Quantity
		SubTotal:    ticket.SubTotal, // Menyertakan SubTotal
	}
}

func (s *ticketService) UpdatePayment(userID, ticketID uint) (*dto.PaymentUpdateResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}
	if ticket.UserID != userID {
		return nil, errors.New("unauthorized access to ticket")
	}
	if ticket.Status != model.Available || ticket.PaymentStatus != model.Pending {
		return nil, errors.New("ticket is not available for payment")
	}

	// ⬇️ Tambahkan validasi waktu event di sini
	event, err := s.eventRepo.FindByID(ticket.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}
	eventTime, err := time.Parse("2006-01-02 15:04:05", event.DateTime)
	if err != nil || time.Now().After(eventTime) {
		return nil, errors.New("event already started or invalid date")
	}

	// ⬇️ Update status pembayaran jika semua valid
	err = s.ticketRepo.UpdatePaymentStatus(ticketID, model.Success, model.Booked)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentUpdateResponse{
		ID:            ticket.ID,
		Status:        string(model.Booked),
		PaymentStatus: string(model.Success),
	}, nil
}

func (s *ticketService) CancelPayment(userID, ticketID uint) (*dto.PaymentUpdateResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil || ticket.UserID != userID {
		return nil, errors.New("unauthorized or ticket not found")
	}

	if err := s.ticketRepo.UpdatePaymentStatus(ticketID, model.Cancel, model.Cancelled); err != nil {
		return nil, err
	}

	return &dto.PaymentUpdateResponse{
		ID:            ticket.ID,
		Status:        string(model.Cancelled),
		PaymentStatus: string(model.Cancel),
	}, nil
}

func (s *ticketService) GetAllTickets(page, limit int) ([]dto.TicketResponse, *dto.Pagination, error) {
	tickets, total, err := s.ticketRepo.FindAllTickets(page, limit)
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
