package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse adalah struktur standar untuk semua response API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success mengirim response sukses
func Success(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error mengirim response error
func Error(ctx *gin.Context, statusCode int, message string, err interface{}) {
	ctx.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// OK adalah shorthand untuk response 200
func OK(ctx *gin.Context, message string, data interface{}) {
	Success(ctx, http.StatusOK, message, data)
}

// Created adalah shorthand untuk response 201
func Created(ctx *gin.Context, message string, data interface{}) {
	Success(ctx, http.StatusCreated, message, data)
}

// BadRequest adalah shorthand untuk response 400
func BadRequest(ctx *gin.Context, message string, err interface{}) {
	Error(ctx, http.StatusBadRequest, message, err)
}

// Unauthorized adalah shorthand untuk response 401
func Unauthorized(ctx *gin.Context, message string) {
	Error(ctx, http.StatusUnauthorized, message, nil)
}

// Forbidden adalah shorthand untuk response 403
func Forbidden(ctx *gin.Context, message string) {
	Error(ctx, http.StatusForbidden, message, nil)
}

// NotFound adalah shorthand untuk response 404
func NotFound(ctx *gin.Context, message string) {
	Error(ctx, http.StatusNotFound, message, nil)
}

// InternalServerError adalah shorthand untuk response 500
func InternalServerError(ctx *gin.Context, message string, err interface{}) {
	Error(ctx, http.StatusInternalServerError, message, err)
}
