package routes

import (
	"log"
	"os"
	"ticketing/config"
	"ticketing/controller"
	"ticketing/middleware"
	"ticketing/model"
	"ticketing/repository"
	"ticketing/service"

	"github.com/gin-gonic/gin"
)

func Run() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Connect to database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Perform auto-migration for models
	if err := db.AutoMigrate(&model.User{}, &model.Event{}, &model.Ticket{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	eventService := service.NewEventService(eventRepo)
	ticketService := service.NewTicketService(ticketRepo, eventRepo)
	reportService := service.NewReportService(reportRepo)
	userService := service.NewUserService(userRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	eventController := controller.NewEventController(eventService)
	ticketController := controller.NewTicketController(ticketService)
	reportController := controller.NewReportController(reportService, db)
	userController := controller.NewUserController(userService)

	// Create Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.ErrorHandler()) // Global error handler

	// Register routes
	SetupRoutes(router, db, authController, userController, eventController, ticketController, reportController)

	// Set server port, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server running on port %s", port)
	log.Fatal(router.Run(":" + port))
}
