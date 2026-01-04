package main

import (
	"log"
	"net/http"

	"github.com/akbar/go-commerce-api/internal/auth/entity"
	authHandler "github.com/akbar/go-commerce-api/internal/auth/handler"
	authMiddleware "github.com/akbar/go-commerce-api/internal/auth/middleware"
	authRepo "github.com/akbar/go-commerce-api/internal/auth/repository"
	authService "github.com/akbar/go-commerce-api/internal/auth/service"
	"github.com/akbar/go-commerce-api/internal/common/response"
	"github.com/akbar/go-commerce-api/pkg/config"
	"github.com/akbar/go-commerce-api/pkg/database"
	"github.com/akbar/go-commerce-api/pkg/utils"
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

	// Auto migrate (hanya untuk development)
	if cfg.App.Env == "development" {
		if err := database.AutoMigrate(db, &entity.User{}); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	}

	// Initialize Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Token blacklist feature will be disabled")
		redisClient = nil
	}

	// ========================================
	// Dependency Injection Setup
	// ========================================

	// JWT Service
	jwtService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	// Auth Module
	userRepository := authRepo.NewUserRepository(db)
	authSvc := authService.NewAuthService(userRepository, jwtService, redisClient)
	authHdl := authHandler.NewAuthHandler(authSvc)

	// ========================================
	// Setup Gin Router
	// ========================================
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

	// ========================================
	// API v1 Routes
	// ========================================
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHdl.Register)
			auth.POST("/login", authHdl.Login)
			auth.POST("/logout", authHdl.Logout)

			// Protected route - requires authentication
			auth.GET("/me", authMiddleware.AuthMiddleware(jwtService, authSvc), authHdl.GetProfile)
		}

		// Protected routes group (requires authentication)
		protected := v1.Group("")
		protected.Use(authMiddleware.AuthMiddleware(jwtService, authSvc))
		{
			// Product routes (akan diimplementasikan)
			protected.GET("/products", func(ctx *gin.Context) {
				response.OK(ctx, "Product module ready", nil)
			})

			// Order routes (akan diimplementasikan)
			protected.GET("/orders", func(ctx *gin.Context) {
				response.OK(ctx, "Order module ready", nil)
			})

			// Payment routes (akan diimplementasikan)
			protected.GET("/payments", func(ctx *gin.Context) {
				response.OK(ctx, "Payment module ready", nil)
			})

			// Admin only routes
			admin := protected.Group("/admin")
			admin.Use(authMiddleware.RoleMiddleware(entity.RoleAdmin))
			{
				admin.GET("/dashboard", func(ctx *gin.Context) {
					response.OK(ctx, "Admin dashboard", nil)
				})
			}

			// Seller only routes
			seller := protected.Group("/seller")
			seller.Use(authMiddleware.RoleMiddleware(entity.RoleSeller, entity.RoleAdmin))
			{
				seller.GET("/dashboard", func(ctx *gin.Context) {
					response.OK(ctx, "Seller dashboard", nil)
				})
			}
		}
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
