package middleware

import (
	"strings"

	"github.com/akbarwjyy/go-commerce-api/internal/auth/service"
	"github.com/akbarwjyy/go-commerce-api/internal/common/response"
	"github.com/akbarwjyy/go-commerce-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware untuk proteksi route yang membutuhkan authentication
func AuthMiddleware(jwtService *utils.JWTService, authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(ctx, "Authorization header required")
			ctx.Abort()
			return
		}

		// Parse token - support both "Bearer <token>" and plain "<token>" format
		var token string
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			// Format: "Bearer <token>"
			token = parts[1]
		} else if len(parts) == 1 {
			// Format: plain token (for Swagger UI compatibility)
			token = parts[0]
		} else {
			response.Unauthorized(ctx, "Invalid authorization format. Use: Bearer <token>")
			ctx.Abort()
			return
		}

		// Cek apakah token ada di blacklist
		if authService.IsTokenBlacklisted(token) {
			response.Unauthorized(ctx, "Token has been revoked")
			ctx.Abort()
			return
		}

		// Validasi token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(ctx, "Invalid or expired token")
			ctx.Abort()
			return
		}

		// Set user info ke context untuk digunakan handler
		ctx.Set("userID", claims.UserID)
		ctx.Set("userEmail", claims.Email)
		ctx.Set("userRole", claims.Role)

		ctx.Next()
	}
}

// RoleMiddleware untuk membatasi akses berdasarkan role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("userRole")
		if !exists {
			response.Unauthorized(ctx, "User role not found")
			ctx.Abort()
			return
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				ctx.Next()
				return
			}
		}

		response.Forbidden(ctx, "You don't have permission to access this resource")
		ctx.Abort()
	}
}
