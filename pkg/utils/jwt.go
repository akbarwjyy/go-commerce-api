package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims custom claims untuk JWT
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService untuk operasi JWT
type JWTService struct {
	secretKey  string
	expireHour int
}

// NewJWTService membuat instance JWTService
func NewJWTService(secretKey string, expireHour int) *JWTService {
	return &JWTService{
		secretKey:  secretKey,
		expireHour: expireHour,
	}
}

// GenerateToken membuat JWT token baru
func (j *JWTService) GenerateToken(userID uint, email, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken memvalidasi dan parse JWT token
func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenExpiry mengembalikan durasi expiry token
func (j *JWTService) GetTokenExpiry() time.Duration {
	return time.Duration(j.expireHour) * time.Hour
}
