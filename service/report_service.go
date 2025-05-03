package service

import (
	"ticketing/dto"
	"ticketing/repository"

	"gorm.io/gorm"
)

type ReportService interface {
	GetSummaryReport(db *gorm.DB) (dto.SummaryReportResponse, error)
	GetEventReports(db *gorm.DB) ([]dto.EventReportResponse, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}

func (s *reportService) GetSummaryReport(db *gorm.DB) (dto.SummaryReportResponse, error) {
	return s.reportRepo.GetSummaryReport(db)
}

func (s *reportService) GetEventReports(db *gorm.DB) ([]dto.EventReportResponse, error) {
	return s.reportRepo.GetEventReports(db)
}
