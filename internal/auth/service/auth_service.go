package service

import (
	"context"
	"errors"
	"time"

	"github.com/akbarwjyy/go-commerce-api/internal/auth/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/auth/entity"
	"github.com/akbarwjyy/go-commerce-api/internal/auth/repository"
	"github.com/akbarwjyy/go-commerce-api/pkg/utils"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthService interface untuk business logic authentication
type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(token string) error
	IsTokenBlacklisted(token string) bool
	GetUserByID(id uint) (*entity.User, error)
}

// authService implementasi AuthService
type authService struct {
	userRepo    repository.UserRepository
	jwtService  *utils.JWTService
	redisClient *redis.Client
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	jwtService *utils.JWTService,
	redisClient *redis.Client,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtService:  jwtService,
		redisClient: redisClient,
	}
}

// Register mendaftarkan user baru
func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set default role jika tidak diisi
	role := req.Role
	if role == "" || !entity.IsValidRole(role) {
		role = entity.RoleUser
	}

	// Buat user baru
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Token: token,
	}, nil
}

// Login melakukan autentikasi user
func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Token: token,
	}, nil
}

// Logout menambahkan token ke blacklist di Redis
func (s *authService) Logout(token string) error {
	if s.redisClient == nil {
		return nil // Skip jika Redis tidak tersedia
	}

	ctx := context.Background()
	// Simpan token di Redis dengan TTL sesuai expiry token
	return s.redisClient.Set(ctx, "blacklist:"+token, "1", s.jwtService.GetTokenExpiry()).Err()
}

// IsTokenBlacklisted mengecek apakah token ada di blacklist
func (s *authService) IsTokenBlacklisted(token string) bool {
	if s.redisClient == nil {
		return false // Skip jika Redis tidak tersedia
	}

	ctx := context.Background()
	result, err := s.redisClient.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false
	}
	return result == "1"
}

// GetUserByID mengambil user berdasarkan ID
func (s *authService) GetUserByID(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
}

// hashPassword helper untuk hash password (tidak diexport)
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash helper untuk verifikasi password (tidak diexport)
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetTokenRemainingTime menghitung sisa waktu token (untuk TTL blacklist)
func GetTokenRemainingTime(expireAt time.Time) time.Duration {
	remaining := time.Until(expireAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
