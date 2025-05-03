package repository

import (
	"ticketing/dto"
	"ticketing/model"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetSummaryReport(db *gorm.DB) (dto.SummaryReportResponse, error)
	GetEventReports(db *gorm.DB) ([]dto.EventReportResponse, error)
}

type reportRepository struct {
	db *gorm.DB
}

// Konstruktor yang mengembalikan pointer ke reportRepository, bukan interface
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetSummaryReport(db *gorm.DB) (dto.SummaryReportResponse, error) {
	var summary dto.SummaryReportResponse
	var totalRevenue float64

	var totalEvents int64
	db.Model(&model.Event{}).Count(&totalEvents)
	summary.TotalEvents = int(totalEvents)

	var totalTickets int64
	db.Model(&model.Ticket{}).Count(&totalTickets)
	summary.TotalTickets = int(totalTickets)

	err := db.Table("tickets").
		Select("COALESCE(SUM(events.price), 0)").
		Joins("JOIN events ON tickets.event_id = events.id").
		Where("tickets.status = ?", model.Booked).
		Scan(&totalRevenue).Error
	if err != nil {
		return summary, err
	}
	summary.TotalRevenue = totalRevenue

	var upcoming, ongoing, completed int64
	db.Model(&model.Event{}).Where("status = ?", model.Upcoming).Count(&upcoming)
	db.Model(&model.Event{}).Where("status = ?", model.Ongoing).Count(&ongoing)
	db.Model(&model.Event{}).Where("status = ?", model.Completed).Count(&completed)
	summary.Upcoming = int(upcoming)
	summary.Ongoing = int(ongoing)
	summary.Completed = int(completed)

	return summary, nil
}

func (r *reportRepository) GetEventReports(db *gorm.DB) ([]dto.EventReportResponse, error) {
	var events []model.Event
	var reports []dto.EventReportResponse

	err := db.Preload("Tickets").Find(&events).Error
	if err != nil {
		return nil, err
	}

	for _, e := range events {
		ticketsSold := 0

		for _, t := range e.Tickets {
			if t.Status == model.Booked {
				ticketsSold++
			}
		}

		revenue := float64(ticketsSold) * e.Price
		occupancyRate := 0.0
		if e.Capacity > 0 {
			occupancyRate = float64(ticketsSold) / float64(e.Capacity) * 100
		}

		report := dto.EventReportResponse{
			EventName:     e.Name,
			TotalCapacity: e.Capacity,
			TicketsSold:   ticketsSold,
			Revenue:       revenue,
			OccupancyRate: occupancyRate,
		}
		reports = append(reports, report)
	}

	return reports, nil
}
