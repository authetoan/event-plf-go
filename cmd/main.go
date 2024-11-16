package main

import (
	"booking-service/configs"
	"booking-service/internal/controllers"
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"booking-service/internal/services"
	"booking-service/pkg/db"
	"booking-service/pkg/kafka"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config := configs.LoadConfig("configs")
	log.Printf("Starting app on port %s in %s mode", config.App.Port, config.App.Env)

	// Initialize database connection
	database, err := db.Connect(
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate database models
	err = database.AutoMigrate(
		&models.Booking{},
		&models.Ticket{},
		&models.Event{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schemas: %v", err)
	}
	log.Println("Database auto-migration completed.")

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer(config.Kafka.Broker)

	// Initialize repositories
	bookingRepo := repositories.NewBookingRepository(database)
	ticketRepo := repositories.NewTicketRepository(database)

	// Initialize services
	ticketService := services.NewTicketService(ticketRepo)
	bookingService := services.NewBookingService(bookingRepo, ticketService, kafkaProducer)

	// Initialize controllers
	bookingController := controllers.NewBookingController(bookingService)
	ticketController := controllers.NewTicketController(ticketService)

	// Set up Gin router
	router := gin.Default()
	apiRoutes := router.Group("/api")

	// Register routes
	controllers.RegisterBookingRoutes(apiRoutes, bookingController)
	controllers.RegisterTicketRoutes(apiRoutes, ticketController)

	// Start the server
	if err := router.Run(":" + config.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
