// Package response provides standard API response helpers
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard API response structure
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success mengirim response sukses dengan data
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

// OK mengirim response sukses 200
func OK(ctx *gin.Context, message string, data interface{}) {
	Success(ctx, http.StatusOK, message, data)
}

// Created mengirim response sukses 201
func Created(ctx *gin.Context, message string, data interface{}) {
	Success(ctx, http.StatusCreated, message, data)
}

// BadRequest mengirim response error 400
func BadRequest(ctx *gin.Context, message string, err interface{}) {
	Error(ctx, http.StatusBadRequest, message, err)
}

// Unauthorized mengirim response error 401
func Unauthorized(ctx *gin.Context, message string) {
	Error(ctx, http.StatusUnauthorized, message, nil)
}

// Forbidden mengirim response error 403
func Forbidden(ctx *gin.Context, message string) {
	Error(ctx, http.StatusForbidden, message, nil)
}

// NotFound mengirim response error 404
func NotFound(ctx *gin.Context, message string) {
	Error(ctx, http.StatusNotFound, message, nil)
}

// InternalServerError mengirim response error 500
func InternalServerError(ctx *gin.Context, message string, err interface{}) {
	Error(ctx, http.StatusInternalServerError, message, err)
}
