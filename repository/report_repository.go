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

	// Menggunakan transaksi untuk memastikan semua query dilakukan bersamaan
	tx := db.Begin()
	if tx.Error != nil {
		return summary, tx.Error
	}

	// Total Events
	var totalEvents int64
	if err := tx.Model(&model.Event{}).Count(&totalEvents).Error; err != nil {
		tx.Rollback()
		return summary, err
	}
	summary.TotalEvents = int(totalEvents)

	// Total Tickets
	var totalTickets int64
	if err := tx.Model(&model.Ticket{}).Count(&totalTickets).Error; err != nil {
		tx.Rollback()
		return summary, err
	}
	summary.TotalTickets = int(totalTickets)

	// Total Revenue
	var totalRevenue float64
	if err := tx.Table("tickets").
		Select("COALESCE(SUM(events.price), 0)").
		Joins("JOIN events ON tickets.event_id = events.id").
		Where("tickets.status = ?", model.Booked).
		Scan(&totalRevenue).Error; err != nil {
		tx.Rollback()
		return summary, err
	}
	summary.TotalRevenue = totalRevenue

	// Event Status Counts (Upcoming, Ongoing, Completed)
	var upcoming, ongoing, completed int64
	if err := tx.Model(&model.Event{}).Where("status = ?", model.Upcoming).Count(&upcoming).Error; err != nil {
		tx.Rollback()
		return summary, err
	}
	if err := tx.Model(&model.Event{}).Where("status = ?", model.Ongoing).Count(&ongoing).Error; err != nil {
		tx.Rollback()
		return summary, err
	}
	if err := tx.Model(&model.Event{}).Where("status = ?", model.Completed).Count(&completed).Error; err != nil {
		tx.Rollback()
		return summary, err
	}

	summary.Upcoming = int(upcoming)
	summary.Ongoing = int(ongoing)
	summary.Completed = int(completed)

	// Commit transaksi setelah semua query berhasil
	if err := tx.Commit().Error; err != nil {
		return summary, err
	}

	return summary, nil
}

func (r *reportRepository) GetEventReports(db *gorm.DB) ([]dto.EventReportResponse, error) {
	var reports []dto.EventReportResponse
	var events []model.Event

	// Preload tiket untuk mendapatkan semua data event dan tiket yang terkait
	if err := db.Preload("Tickets").Find(&events).Error; err != nil {
		return nil, err
	}

	// Menghitung laporan untuk setiap event
	for _, e := range events {
		// Menghitung tiket yang terjual
		ticketsSold := 0
		for _, t := range e.Tickets {
			if t.Status == model.Booked {
				ticketsSold++
			}
		}

		// Menghitung pendapatan dan tingkat hunian
		revenue := float64(ticketsSold) * e.Price
		occupancyRate := 0.0
		if e.Capacity > 0 {
			occupancyRate = float64(ticketsSold) / float64(e.Capacity) * 100
		}

		// Menyusun laporan untuk event ini
		report := dto.EventReportResponse{
			EventName:     e.Name,
			TotalCapacity: e.Capacity,
			TicketsSold:   ticketsSold,
			Revenue:       revenue,
			OccupancyRate: occupancyRate,
		}

		// Menambahkan laporan ke dalam list
		reports = append(reports, report)
	}

	return reports, nil
}
