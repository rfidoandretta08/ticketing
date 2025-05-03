package routes

import (
	"ticketing/controller"
	"ticketing/middleware"

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
) {
	api := r.Group("/api")

	// AUTH routes (tanpa middleware)
	api.POST("/register", authController.Register)
	api.POST("/login", authController.Login)

	// USER routes
	userGroup := api.Group("/users")
	{
		userGroup.Use(middleware.AuthMiddleware("admin")) // hanya admin boleh
		userGroup.POST("/", userController.CreateUser)
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
	}

	// REPORT routes (admin only)
	reportGroup := api.Group("/reports")
	reportGroup.Use(middleware.AuthMiddleware("admin"))
	{
		reportGroup.GET("/summary", reportController.GetSummaryReport)
		reportGroup.GET("/events", reportController.GetEventReports)
	}
}
