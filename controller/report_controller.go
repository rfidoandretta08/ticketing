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

func (r *ReportController) GetSummaryReport(c *gin.Context) {
	// Ambil data summary report
	report, err := r.reportService.GetSummaryReport(r.db) // Lakukan query jika perlu
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get summary report"})
		return
	}
	c.JSON(200, report)
}

// Endpoint untuk mendapatkan event reports
func (r *ReportController) GetEventReports(c *gin.Context) {
	// Ambil data event report
	reports, err := r.reportService.GetEventReports(r.db) // Lakukan query jika perlu
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get event reports"})
		return
	}
	c.JSON(200, reports)
}

// Method untuk generate summary report dalam format PDF
func (ctrl *ReportController) GenerateSummaryReportPDF(c *gin.Context) {
	pdfBytes, err := ctrl.reportService.GenerateSummaryReportPDF(ctrl.db)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/pdf", pdfBytes)
}

// Method untuk generate event report dalam format PDF
func (ctrl *ReportController) GenerateEventReportPDF(c *gin.Context) {
	// Menghasilkan laporan event dalam format PDF
	err := ctrl.reportService.GenerateEventReportPDF(ctrl.db)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Event report PDF generated successfully"})
}
