package handler

import (
	"net/http"
	"strings"

	"github.com/akbarwjyy/go-commerce-api/internal/auth/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/auth/service"
	"github.com/akbarwjyy/go-commerce-api/internal/common/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler menangani HTTP request untuk authentication
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register request"
// @Success      201 {object} response.APIResponse{data=dto.AuthResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      409 {object} response.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			response.Error(ctx, http.StatusConflict, "Email already registered", nil)
			return
		}
		response.InternalServerError(ctx, "Failed to register user", err.Error())
		return
	}

	response.Created(ctx, "User registered successfully", result)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login request"
// @Success      200 {object} response.APIResponse{data=dto.AuthResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Unauthorized(ctx, "Invalid email or password")
			return
		}
		response.InternalServerError(ctx, "Failed to login", err.Error())
		return
	}

	response.OK(ctx, "Login successful", result)
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout user and blacklist the token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// Ambil token dari header Authorization
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(ctx, "Authorization header required")
		return
	}

	// Parse Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Unauthorized(ctx, "Invalid authorization format")
		return
	}
	token := parts[1]

	if err := h.authService.Logout(token); err != nil {
		response.InternalServerError(ctx, "Failed to logout", err.Error())
		return
	}

	response.OK(ctx, "Logout successful", nil)
}

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Get the profile of the currently authenticated user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure      401 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /auth/me [get]
func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	// User ID diambil dari context (diset oleh middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		response.NotFound(ctx, "User not found")
		return
	}

	response.OK(ctx, "Profile retrieved successfully", dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}
