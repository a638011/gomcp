package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the global logger instance
	Logger zerolog.Logger
)

// Config holds logger configuration
type Config struct {
	Level  string
	Output io.Writer
}

// Init initializes the global logger
func Init(level string) {
	initLogger(Config{
		Level:  level,
		Output: os.Stdout,
	})
}

// InitWithWriter initializes the logger with a custom writer
func InitWithWriter(level string, output io.Writer) {
	initLogger(Config{
		Level:  level,
		Output: output,
	})
}

func initLogger(cfg Config) {
	// Set up zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"

	// Parse log level
	logLevel := parseLogLevel(cfg.Level)

	// Create logger with JSON output
	Logger = zerolog.New(cfg.Output).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set global logger
	log.Logger = Logger
}

// parseLogLevel converts string log level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARNING", "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "CRITICAL", "FATAL":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// GetLogger returns the global logger
func GetLogger() *zerolog.Logger {
	return &Logger
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	Logger.Info().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

// Error logs an error message
func Error(msg string) {
	Logger.Error().Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string) {
	Logger.Fatal().Msg(msg)
}

// WithField returns a logger with an additional field
func WithField(key string, value interface{}) *zerolog.Logger {
	l := Logger.With().Interface(key, value).Logger()
	return &l
}

// WithFields returns a logger with multiple fields
func WithFields(fields map[string]interface{}) *zerolog.Logger {
	ctx := Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	l := ctx.Logger()
	return &l
}

// GetUvicornLogConfig returns a compatible log configuration
// This is kept for API compatibility but not directly used in Go
func GetUvicornLogConfig(level string) map[string]interface{} {
	return map[string]interface{}{
		"version":                  1,
		"disable_existing_loggers": false,
		"formatters": map[string]interface{}{
			"default": map[string]interface{}{
				"format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
			},
		},
		"handlers": map[string]interface{}{
			"default": map[string]interface{}{
				"class":     "logging.StreamHandler",
				"formatter": "default",
				"stream":    "ext://sys.stdout",
			},
		},
		"loggers": map[string]interface{}{
			"": map[string]interface{}{
				"handlers": []string{"default"},
				"level":    level,
			},
		},
	}
}

// ForceReconfigureAllLoggers reconfigures the logger with a new level
func ForceReconfigureAllLoggers(level string) {
	Init(level)
}
