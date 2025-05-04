package routes

import (
	"ticketing/controller"
	"ticketing/middleware"
	"ticketing/service" // Import service untuk memanggil laporan PDF

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(
	r *gin.Engine,
	db *gorm.DB,
	authController *controller.AuthController,
	userController *controller.UserController,
	eventController *controller.EventController,
	ticketController *controller.TicketController,
	reportController *controller.ReportController,
	reportService service.ReportService, // Gunakan service untuk laporan
) {
	api := r.Group("/api")

	// AUTH routes (tanpa middleware)
	api.POST("/register", authController.Register)
	api.POST("/login", authController.Login)

	// USER routes
	userGroup := api.Group("/users")
	{
		userGroup.Use(middleware.AuthMiddleware("admin")) // hanya admin boleh
		userGroup.GET("/", userController.GetAllUsers)
		userGroup.GET("/:id", userController.GetUserByID)
	}

	// EVENT routes
	eventGroup := api.Group("/events")
	{
		eventGroup.GET("", eventController.GetAllEvents)     // publik
		eventGroup.GET("/:id", eventController.GetEventByID) // publik

		eventGroup.Use(middleware.AuthMiddleware("admin")) // hanya admin boleh buat, update, hapus
		eventGroup.POST("", eventController.CreateEvent)
		eventGroup.PUT("/:id", eventController.UpdateEvent)
		eventGroup.DELETE("/:id", eventController.DeleteEvent)
	}

	// TICKET routes (user)
	ticketGroup := api.Group("/tickets")
	ticketGroup.Use(middleware.AuthMiddleware("user"))
	{
		ticketGroup.POST("", ticketController.PurchaseTicket)
		ticketGroup.GET("", ticketController.GetUserTickets)
		ticketGroup.GET("/:id", ticketController.GetTicketByID)
		ticketGroup.PATCH("/:id", ticketController.CancelTicket)
		ticketGroup.PATCH("/:id/payment", ticketController.UpdatePayment)
		ticketGroup.PATCH("/:id/cancel-payment", ticketController.CancelPayment)
	}

	// REPORT routes (admin only)
	reportGroup := api.Group("/reports")
	reportGroup.Use(middleware.AuthMiddleware("admin"))
	{
		reportGroup.GET("/summary", reportController.GetSummaryReport)
		reportGroup.GET("/events", reportController.GetEventReports)
		reportGroup.GET("/ticket", ticketController.GetAllTickets)

		// Route untuk generate summary report PDF
		reportGroup.GET("/generate-summary-excel", func(c *gin.Context) {
			err := reportService.GenerateSummaryReportExcel(db) // Gunakan service untuk generate PDF
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "Summary report Excel generated successfully"})
		})
		reportGroup.GET("/generate-event-excel", func(c *gin.Context) {
			err := reportService.GenerateEventReportExcel(db) // Gunakan service untuk generate PDF
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "Event report Excel generated successfully"})
		})

		reportGroup.GET("/generate-summary-pdf", func(c *gin.Context) {
			// Mengambil dua nilai yang dikembalikan oleh GenerateSummaryReportPDF
			pdfData, err := reportService.GenerateSummaryReportPDF(db)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			// Mengirimkan PDF sebagai response (contoh: sebagai file download)
			c.Header("Content-Type", "application/pdf")
			c.Header("Content-Disposition", "attachment; filename=SummaryReport.pdf")
			c.Data(200, "application/pdf", pdfData) // Menyertakan PDF dalam response
		})

		// Route untuk generate event report PDF
		reportGroup.GET("/generate-event-pdf", func(c *gin.Context) {
			err := reportService.GenerateEventReportPDF(db) // Gunakan service untuk generate PDF
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "Event report PDF generated successfully"})
		})
	}
}
