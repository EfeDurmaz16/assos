package main

import (
	"log"
	"os"

	"assos/api-gateway/internal/config"
	"assos/api-gateway/internal/database"
	"assos/api-gateway/internal/handlers"
	"assos/api-gateway/internal/middleware"
	"assos/api-gateway/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient := database.ConnectRedis(cfg.RedisURL)

	// Initialize NATS
	natsConn, err := database.ConnectNATS(cfg.NatsURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer natsConn.Close()

	// Initialize services
	userService := services.NewUserService(db, redisClient)
	channelService := services.NewChannelService(db, redisClient)
	videoService := services.NewVideoService(db, redisClient, natsConn)
	aiService := services.NewAIService(db, redisClient, natsConn)
	analyticsService := services.NewAnalyticsService(cfg.AnalyticsServiceURL)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, cfg.JWTSecret)
	channelHandler := handlers.NewChannelHandler(channelService)
	videoHandler := handlers.NewVideoHandler(videoService)
	aiHandler := handlers.NewAIHandler(aiService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60, // 1 minute
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "assos-api-gateway",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("/", middleware.JWTAuth(cfg.JWTSecret))

	// User routes
	users := protected.Group("/users")
	users.Get("/me", authHandler.GetProfile)
	users.Put("/me", authHandler.UpdateProfile)

	// Channel routes
	channels := protected.Group("/channels")
	channels.Get("/", channelHandler.GetChannels)
	channels.Post("/", channelHandler.CreateChannel)
	channels.Get("/:id", channelHandler.GetChannel)
	channels.Put("/:id", channelHandler.UpdateChannel)
	channels.Delete("/:id", channelHandler.DeleteChannel)

	// Video routes
	videos := protected.Group("/videos")
	videos.Get("/", videoHandler.GetVideos)
	videos.Post("/", videoHandler.CreateVideo)
	videos.Get("/:id", videoHandler.GetVideo)
	videos.Put("/:id", videoHandler.UpdateVideo)
	videos.Delete("/:id", videoHandler.DeleteVideo)
	videos.Post("/:id/process", videoHandler.ProcessVideo)

	// AI Agent routes
	ai := protected.Group("/ai")
	ai.Get("/agents", aiHandler.GetAgents)
	ai.Post("/agents/:id/task", aiHandler.CreateTask)
	ai.Get("/tasks", aiHandler.GetTasks)
	ai.Get("/tasks/:id", aiHandler.GetTask)

	// Analytics routes
	analytics := protected.Group("/analytics")
	analytics.Get("/dashboard", analyticsHandler.GetDashboard)
	analytics.Get("/performance/:video_id", analyticsHandler.GetPerformance)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting ASSOS API Gateway on port %s", port)
	log.Fatal(app.Listen(":" + port))
}