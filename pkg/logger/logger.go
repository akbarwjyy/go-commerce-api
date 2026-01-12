package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global logger
func Init(env string) {
	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Set log level based on environment
	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		// JSON output for production
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		// Pretty console output for development
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Logger()
	}
}

// GinLogger returns a gin middleware for structured logging
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		status := c.Writer.Status()

		event := log.Info()
		if status >= 400 && status < 500 {
			event = log.Warn()
		} else if status >= 500 {
			event = log.Error()
		}

		event.
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP Request")
	}
}

// Info logs info level message
func Info() *zerolog.Event {
	return log.Info()
}

// Debug logs debug level message
func Debug() *zerolog.Event {
	return log.Debug()
}

// Warn logs warning level message
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs error level message
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal logs fatal level message and exits
func Fatal() *zerolog.Event {
	return log.Fatal()
}
