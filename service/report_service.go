package service

import (
	"bytes"
	"log"
	"strconv"
	"ticketing/dto"
	"ticketing/repository"
	"ticketing/utils"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ReportService interface {
	GetSummaryReport(db *gorm.DB) (dto.SummaryReportResponse, error)
	GetEventReports(db *gorm.DB) ([]dto.EventReportResponse, error)
	GenerateSummaryReportExcel(db *gorm.DB) error
	GenerateEventReportExcel(db *gorm.DB) error
	GenerateSummaryReportPDF(db *gorm.DB) ([]byte, error)
	GenerateEventReportPDF(db *gorm.DB) error
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

func (s *reportService) GenerateSummaryReportExcel(db *gorm.DB) error {
	// Ambil data laporan
	summaryReport, err := s.GetSummaryReport(db)
	if err != nil {
		return err
	}

	// Membuat file Excel baru
	f := excelize.NewFile()
	sheet := "Summary Report"
	index, err := f.NewSheet(sheet) // Correctly handle both values returned by NewSheet
	if err != nil {
		log.Fatalf("failed to create new sheet: %v", err)
		return err
	}

	// Menulis header
	headers := []string{"Total Events", "Total Tickets", "Total Revenue", "Upcoming Events", "Ongoing Events", "Completed Events"}
	for col, header := range headers {
		cell, _ := excelize.ColumnNumberToName(col + 1) // Start from 1
		f.SetCellValue(sheet, cell+"1", header)
	}

	// Menulis data
	f.SetCellValue(sheet, "A2", summaryReport.TotalEvents)
	f.SetCellValue(sheet, "B2", summaryReport.TotalTickets)
	f.SetCellValue(sheet, "C2", summaryReport.TotalRevenue)
	f.SetCellValue(sheet, "D2", summaryReport.Upcoming)
	f.SetCellValue(sheet, "E2", summaryReport.Ongoing)
	f.SetCellValue(sheet, "F2", summaryReport.Completed)

	// Set sheet as active
	f.SetActiveSheet(index)

	// Simpan file
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		log.Printf("failed to write Excel to buffer: %v", err)
		return err
	}

	// Simpan ke folder report
	err = utils.SaveFileToReportFolder("SummaryReport.xlsx", buf.Bytes())
	if err != nil {
		log.Printf("failed to save Excel file: %v", err)
		return err
	}

	return nil
}

func (s *reportService) GenerateEventReportExcel(db *gorm.DB) error {
	// Ambil data laporan
	eventReports, err := s.GetEventReports(db)
	if err != nil {
		return err
	}

	// Membuat file Excel baru
	f := excelize.NewFile()
	sheet := "Event Report"
	index, err := f.NewSheet(sheet) // Correctly handle both values returned by NewSheet
	if err != nil {
		log.Fatalf("failed to create new sheet: %v", err)
		return err
	}

	// Menulis header
	headers := []string{"Event Name", "Total Capacity", "Tickets Sold", "Revenue", "Occupancy Rate"}
	for col, header := range headers {
		cell, _ := excelize.ColumnNumberToName(col + 1) // Start from 1
		f.SetCellValue(sheet, cell+"1", header)
	}

	// Menulis data
	for i, eventReport := range eventReports {
		f.SetCellValue(sheet, "A"+strconv.Itoa(i+2), eventReport.EventName)
		f.SetCellValue(sheet, "B"+strconv.Itoa(i+2), eventReport.TotalCapacity)
		f.SetCellValue(sheet, "C"+strconv.Itoa(i+2), eventReport.TicketsSold)
		f.SetCellValue(sheet, "D"+strconv.Itoa(i+2), eventReport.Revenue)
		f.SetCellValue(sheet, "E"+strconv.Itoa(i+2), eventReport.OccupancyRate)
	}

	// Set sheet as active
	f.SetActiveSheet(index)

	// Simpan file
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		log.Printf("failed to write Excel to buffer: %v", err)
		return err
	}

	err = utils.SaveFileToReportFolder("EventReport.xlsx", buf.Bytes())
	if err != nil {
		log.Printf("failed to save Excel file: %v", err)
		return err
	}
	return nil
}

func (s *reportService) GenerateSummaryReportPDF(db *gorm.DB) ([]byte, error) {
	summaryReport, err := s.GetSummaryReport(db)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Summary Report")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Total Events: "+strconv.Itoa(summaryReport.TotalEvents))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Total Tickets: "+strconv.Itoa(summaryReport.TotalTickets))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Total Revenue: $"+strconv.FormatFloat(summaryReport.TotalRevenue, 'f', 2, 64))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Upcoming Events: "+strconv.Itoa(summaryReport.Upcoming))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Ongoing Events: "+strconv.Itoa(summaryReport.Ongoing))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Completed Events: "+strconv.Itoa(summaryReport.Completed))

	// Menulis ke buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		log.Printf("failed to write PDF to buffer: %v", err)
		return nil, err
	}

	err = utils.SaveFileToReportFolder("SummaryReport.pdf", buf.Bytes())
	if err != nil {
		log.Printf("failed to save PDF file: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *reportService) GenerateEventReportPDF(db *gorm.DB) error {
	eventReports, err := s.GetEventReports(db)
	if err != nil {
		return err
	}

	// Membuat instance PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font untuk header
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Event Report")

	// Menulis header
	pdf.Ln(10)
	pdf.Cell(40, 10, "Event Name")
	pdf.Cell(40, 10, "Total Capacity")
	pdf.Cell(40, 10, "Tickets Sold")
	pdf.Cell(40, 10, "Revenue")
	pdf.Cell(40, 10, "Occupancy Rate")
	pdf.Ln(10)

	// Menulis data setiap event
	pdf.SetFont("Arial", "", 12)
	for _, eventReport := range eventReports {
		pdf.Cell(40, 10, eventReport.EventName)
		pdf.Cell(40, 10, strconv.Itoa(eventReport.TotalCapacity))
		pdf.Cell(40, 10, strconv.Itoa(eventReport.TicketsSold))
		pdf.Cell(40, 10, "$"+strconv.FormatFloat(eventReport.Revenue, 'f', 2, 64))
		pdf.Cell(40, 10, strconv.FormatFloat(eventReport.OccupancyRate, 'f', 2, 64)+"%")
		pdf.Ln(10)
	}

	// Menyimpan file PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		log.Printf("failed to write PDF to buffer: %v", err)
		return err
	}

	err = utils.SaveFileToReportFolder("EventReport.pdf", buf.Bytes())
	if err != nil {
		log.Printf("failed to save PDF file: %v", err)
		return err
	}
	return nil

}
