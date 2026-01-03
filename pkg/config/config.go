package config

import (
	"os"
)

// Config menyimpan konfigurasi aplikasi
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// AppConfig untuk konfigurasi aplikasi
type AppConfig struct {
	Name string
	Env  string
	Port string
}

// DatabaseConfig untuk konfigurasi PostgreSQL
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig untuk konfigurasi Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig untuk konfigurasi JWT
type JWTConfig struct {
	Secret     string
	ExpireHour int
}

// Load membaca konfigurasi dari environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-commerce-api"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "go_commerce"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireHour: 24,
		},
	}
}

// getEnv membaca env variable dengan default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
