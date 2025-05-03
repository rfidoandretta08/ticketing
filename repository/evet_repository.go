package repository

import (
	"ticketing/model"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
	FindAll(page, limit int, search string) ([]model.Event, int64, error)
	FindByID(id uint) (*model.Event, error)
	Update(event *model.Event) error
	Delete(id uint) error
	GetAvailableTickets(eventID uint) (int, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindAll(page, limit int, search string) ([]model.Event, int64, error) {
	var events []model.Event
	var total int64

	query := r.db.Model(&model.Event{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR location LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&events).Error
	return events, total, err
}

func (r *eventRepository) FindByID(id uint) (*model.Event, error) {
	var event model.Event
	err := r.db.Preload("Tickets").First(&event, id).Error
	return &event, err
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&model.Event{}, id).Error
}

func (r *eventRepository) GetAvailableTickets(eventID uint) (int, error) {
	var event model.Event
	if err := r.db.First(&event, eventID).Error; err != nil {
		return 0, err
	}

	var bookedTickets int64
	if err := r.db.Model(&model.Ticket{}).
		Where("event_id = ? AND status = ?", eventID, model.Booked).
		Count(&bookedTickets).Error; err != nil {
		return 0, err
	}

	return event.Capacity - int(bookedTickets), nil
}
