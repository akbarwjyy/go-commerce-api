// Go-Commerce API
//
// Backend E-Commerce REST API built with Golang (Gin) using Modular Monolith architecture.
//
//	@title			Go-Commerce API
//	@version		1.0
//	@description	E-Commerce REST API with modular monolith architecture
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	Go-Commerce API Support
//	@contact.email	support@go-commerce.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"
	"net/http"

	authEntity "github.com/akbarwjyy/go-commerce-api/internal/auth/entity"
	authHandler "github.com/akbarwjyy/go-commerce-api/internal/auth/handler"
	authMiddleware "github.com/akbarwjyy/go-commerce-api/internal/auth/middleware"
	authRepo "github.com/akbarwjyy/go-commerce-api/internal/auth/repository"
	authService "github.com/akbarwjyy/go-commerce-api/internal/auth/service"
	"github.com/akbarwjyy/go-commerce-api/internal/common/response"
	orderEntity "github.com/akbarwjyy/go-commerce-api/internal/order/entity"
	orderHandler "github.com/akbarwjyy/go-commerce-api/internal/order/handler"
	orderRepo "github.com/akbarwjyy/go-commerce-api/internal/order/repository"
	orderService "github.com/akbarwjyy/go-commerce-api/internal/order/service"
	paymentEntity "github.com/akbarwjyy/go-commerce-api/internal/payment/entity"
	paymentHandler "github.com/akbarwjyy/go-commerce-api/internal/payment/handler"
	paymentRepo "github.com/akbarwjyy/go-commerce-api/internal/payment/repository"
	paymentService "github.com/akbarwjyy/go-commerce-api/internal/payment/service"
	productEntity "github.com/akbarwjyy/go-commerce-api/internal/product/entity"
	productHandler "github.com/akbarwjyy/go-commerce-api/internal/product/handler"
	productRepo "github.com/akbarwjyy/go-commerce-api/internal/product/repository"
	productService "github.com/akbarwjyy/go-commerce-api/internal/product/service"
	"github.com/akbarwjyy/go-commerce-api/pkg/config"
	"github.com/akbarwjyy/go-commerce-api/pkg/database"
	"github.com/akbarwjyy/go-commerce-api/pkg/utils"
	"github.com/gin-gonic/gin"

	_ "github.com/akbarwjyy/go-commerce-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		if err := database.AutoMigrate(db,
			&authEntity.User{},
			&productEntity.Category{},
			&productEntity.Product{},
			&orderEntity.Order{},
			&orderEntity.OrderItem{},
			&paymentEntity.Payment{},
		); err != nil {
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

	// Product Module
	categoryRepository := productRepo.NewCategoryRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	productSvc := productService.NewProductService(productRepository, categoryRepository, db)
	productHdl := productHandler.NewProductHandler(productSvc)

	// Order Module
	orderRepository := orderRepo.NewOrderRepository(db)
	orderSvc := orderService.NewOrderService(orderRepository, productSvc, db)
	orderHdl := orderHandler.NewOrderHandler(orderSvc)

	// Payment Module
	paymentRepository := paymentRepo.NewPaymentRepository(db)
	paymentSvc := paymentService.NewPaymentService(paymentRepository, orderSvc, db)
	paymentHdl := paymentHandler.NewPaymentHandler(paymentSvc)

	// ========================================
	// Setup Gin Router
	// ========================================
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

		// Categories routes (public read, protected write)
		categories := v1.Group("/categories")
		{
			categories.GET("", productHdl.GetAllCategories)
			categories.GET("/:id", productHdl.GetCategory)

			// Admin only - create/update/delete categories
			categories.Use(authMiddleware.AuthMiddleware(jwtService, authSvc))
			categories.Use(authMiddleware.RoleMiddleware(authEntity.RoleAdmin))
			categories.POST("", productHdl.CreateCategory)
			categories.PUT("/:id", productHdl.UpdateCategory)
			categories.DELETE("/:id", productHdl.DeleteCategory)
		}

		// Products routes (public read)
		products := v1.Group("/products")
		{
			products.GET("", productHdl.GetAllProducts)
			products.GET("/:id", productHdl.GetProduct)
		}

		// Protected routes group (requires authentication)
		protected := v1.Group("")
		protected.Use(authMiddleware.AuthMiddleware(jwtService, authSvc))
		{
			// Product management (seller/admin only)
			protectedProducts := protected.Group("/products")
			protectedProducts.Use(authMiddleware.RoleMiddleware(authEntity.RoleSeller, authEntity.RoleAdmin))
			{
				protectedProducts.POST("", productHdl.CreateProduct)
				protectedProducts.PUT("/:id", productHdl.UpdateProduct)
				protectedProducts.DELETE("/:id", productHdl.DeleteProduct)
				protectedProducts.PATCH("/:id/stock", productHdl.UpdateStock)
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("/checkout", orderHdl.Checkout)
				orders.GET("", orderHdl.GetMyOrders)
				orders.GET("/:id", orderHdl.GetOrder)
				orders.PATCH("/:id/status", orderHdl.UpdateOrderStatus)
				orders.POST("/:id/cancel", orderHdl.CancelOrder)
				orders.GET("/:id/payment", paymentHdl.GetPaymentByOrder)
			}

			// Payment routes
			payments := protected.Group("/payments")
			{
				payments.POST("", paymentHdl.CreatePayment)
				payments.GET("", paymentHdl.GetMyPayments)
				payments.GET("/:id", paymentHdl.GetPayment)
				payments.POST("/callback", paymentHdl.PaymentCallback) // For testing
			}

			// Seller routes
			seller := protected.Group("/seller")
			seller.Use(authMiddleware.RoleMiddleware(authEntity.RoleSeller, authEntity.RoleAdmin))
			{
				seller.GET("/dashboard", func(ctx *gin.Context) {
					response.OK(ctx, "Seller dashboard", nil)
				})
				seller.GET("/products", productHdl.GetMyProducts)
			}

			// Admin only routes
			admin := protected.Group("/admin")
			admin.Use(authMiddleware.RoleMiddleware(authEntity.RoleAdmin))
			{
				admin.GET("/dashboard", func(ctx *gin.Context) {
					response.OK(ctx, "Admin dashboard", nil)
				})
				admin.GET("/orders", orderHdl.GetAllOrders)
				admin.GET("/payments", paymentHdl.GetAllPayments)
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
	log.Printf("Swagger docs available at http://localhost%s/swagger/index.html", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
