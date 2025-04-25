package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initialize sets up the global logger with the specified configuration
func Initialize(level string, prettyPrint bool) error {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %v", level, err)
	}
	zerolog.SetGlobalLevel(logLevel)

	// Configure logger output
	if prettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	return nil
}

// Logger interface defines the logging methods
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
}

// AppLogger implements the Logger interface using zerolog
type AppLogger struct {
	logger zerolog.Logger
}

// NewLogger creates a new AppLogger instance
func NewLogger() *AppLogger {
	return &AppLogger{
		logger: log.With().Caller().Logger(),
	}
}

// Info logs an info level message
func (l *AppLogger) Info(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs an error level message
func (l *AppLogger) Error(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Debug logs a debug level message
func (l *AppLogger) Debug(msg string, fields map[string]interface{}) {
	event := l.logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs a warning level message
func (l *AppLogger) Warn(msg string, fields map[string]interface{}) {
	event := l.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}
