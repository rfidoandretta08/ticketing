package controller

import (
	"ticketing/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportController struct {
	reportService service.ReportService
	db            *gorm.DB
}

// Konstruktor untuk ReportController
func NewReportController(reportService service.ReportService, db *gorm.DB) *ReportController {
	return &ReportController{reportService: reportService, db: db}
}

// Method untuk mendapatkan summary report
func (ctrl *ReportController) GetSummaryReport(c *gin.Context) {
	// Implementasi untuk mendapatkan summary report
	summaryReport, err := ctrl.reportService.GetSummaryReport(ctrl.db)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, summaryReport)
}

// Method untuk mendapatkan event reports
func (ctrl *ReportController) GetEventReports(c *gin.Context) {
	// Implementasi untuk mendapatkan event reports
	eventReports, err := ctrl.reportService.GetEventReports(ctrl.db)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, eventReports)
}
