package main

import (
	"log"
	"os"

	"assos/analytics-service/internal/database"
	"assos/analytics-service/internal/handlers"
	"assos/analytics-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Configuration
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Connect to database
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services and handlers
	analyticsService := services.NewAnalyticsService(db)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Create Fiber app
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "assos-analytics-service",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	analytics := api.Group("/analytics")
	analytics.Get("/dashboard", analyticsHandler.GetDashboard)
	analytics.Get("/performance/:video_id", analyticsHandler.GetPerformance)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting ASSOS Analytics Service on port %s", port)
	log.Fatal(app.Listen(":" + port))
}