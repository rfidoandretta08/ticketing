package repository

import (
	"ticketing/model"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ticket *model.Ticket) error
	FindAll(page, limit int, userID uint) ([]model.Ticket, int64, error)
	FindAllTickets(page, limit int) ([]model.Ticket, int64, error)
	FindByID(id uint) (*model.Ticket, error)
	Update(ticket *model.Ticket) error
	Cancel(id uint) error
	UpdatePaymentStatus(ticketID uint, status model.PaymentStatus, ticketStatus model.TicketStatus) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *model.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) FindAll(page, limit int, userID uint) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	query := r.db.Model(&model.Ticket{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).
		Preload("Event").
		Find(&tickets).Error
	return tickets, total, err
}

func (r *ticketRepository) FindByID(id uint) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Preload("Event").Preload("User").First(&ticket, id).Error
	return &ticket, err
}

func (r *ticketRepository) Update(ticket *model.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *ticketRepository) Cancel(id uint) error {
	return r.db.Model(&model.Ticket{}).
		Where("id = ?", id).
		Update("status", model.Cancelled).Error
}

func (r *ticketRepository) UpdatePaymentStatus(ticketID uint, paymentStatus model.PaymentStatus, status model.TicketStatus) error {
	return r.db.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"payment_status": paymentStatus,
			"status":         status,
		}).Error
}

func (r *ticketRepository) FindAllTickets(page, limit int) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	offset := (page - 1) * limit
	query := r.db.Preload("Event").Preload("User")

	if err := query.Offset(offset).Limit(limit).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&model.Ticket{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}
