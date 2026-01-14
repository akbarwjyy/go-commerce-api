package service

import (
	"testing"

	"github.com/akbarwjyy/go-commerce-api/internal/auth/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/auth/entity"
	"github.com/stretchr/testify/assert"
)

// Test Register Request Validation
func TestRegisterRequest(t *testing.T) {
	req := &dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "Password123",
		Role:     "user",
	}

	assert.Equal(t, "Test User", req.Name)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Password123", req.Password)
	assert.Equal(t, "user", req.Role)
}

// Test Login Request Validation
func TestLoginRequest(t *testing.T) {
	req := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123",
	}

	assert.NotEmpty(t, req.Email)
	assert.NotEmpty(t, req.Password)
}

// Test User Entity IsAdmin
func TestUserEntity_IsAdmin(t *testing.T) {
	user := &entity.User{
		ID:   1,
		Role: entity.RoleAdmin,
	}

	assert.True(t, user.IsAdmin())
	assert.False(t, user.IsSeller())
}

// Test User Entity IsSeller
func TestUserEntity_IsSeller(t *testing.T) {
	user := &entity.User{
		ID:   1,
		Role: entity.RoleSeller,
	}

	assert.True(t, user.IsSeller())
	assert.False(t, user.IsAdmin())
}

// Test User Entity Role Validation
func TestUserEntity_IsValidRole(t *testing.T) {
	assert.True(t, entity.IsValidRole("admin"))
	assert.True(t, entity.IsValidRole("seller"))
	assert.True(t, entity.IsValidRole("user"))
	assert.False(t, entity.IsValidRole("invalid"))
	assert.False(t, entity.IsValidRole(""))
}

// Test User Response
func TestUserResponse(t *testing.T) {
	user := &entity.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Test User", response.Name)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "user", response.Role)
}

// Test Auth Response
func TestAuthResponse(t *testing.T) {
	response := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
			Role:  "user",
		},
		Token: "jwt-token-here",
	}

	assert.Equal(t, "jwt-token-here", response.Token)
	assert.Equal(t, uint(1), response.User.ID)
}

// Test Error Constants
func TestAuthErrors(t *testing.T) {
	assert.NotNil(t, ErrEmailAlreadyExists)
	assert.NotNil(t, ErrInvalidCredentials)
	assert.NotNil(t, ErrUserNotFound)
}

// Test Role Constants
func TestRoleConstants(t *testing.T) {
	assert.Equal(t, "admin", entity.RoleAdmin)
	assert.Equal(t, "seller", entity.RoleSeller)
	assert.Equal(t, "user", entity.RoleUser)
}
