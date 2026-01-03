package main

import (
	"log"
	"net/http"

	"github.com/akbar/go-commerce-api/internal/common/response"
	"github.com/akbar/go-commerce-api/pkg/config"
	"github.com/akbar/go-commerce-api/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connections
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		// Redis bersifat opsional untuk development
	}

	// Untuk menghindari warning unused variable sementara
	_ = db
	_ = redisClient

	// Setup Gin router
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		response.OK(ctx, "Service is healthy", gin.H{
			"app":    cfg.App.Name,
			"env":    cfg.App.Env,
			"status": "running",
		})
	})

	// API v1 routes group
	v1 := router.Group("/api/v1")
	{
		// Auth routes (akan diimplementasikan)
		v1.GET("/auth", func(ctx *gin.Context) {
			response.OK(ctx, "Auth module ready", nil)
		})

		// Product routes (akan diimplementasikan)
		v1.GET("/products", func(ctx *gin.Context) {
			response.OK(ctx, "Product module ready", nil)
		})

		// Order routes (akan diimplementasikan)
		v1.GET("/orders", func(ctx *gin.Context) {
			response.OK(ctx, "Order module ready", nil)
		})

		// Payment routes (akan diimplementasikan)
		v1.GET("/payments", func(ctx *gin.Context) {
			response.OK(ctx, "Payment module ready", nil)
		})
	}

	// Handle 404
	router.NoRoute(func(ctx *gin.Context) {
		response.NotFound(ctx, "Route not found")
	})

	// Start server
	serverAddr := ":" + cfg.App.Port
	log.Printf("Starting %s server on %s (env: %s)", cfg.App.Name, serverAddr, cfg.App.Env)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
