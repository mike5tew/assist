package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// Init initializes the global logger. If level is empty, falls back to LOG_LEVEL env or "info".
func Init(level string) error {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	lvl := level
	if lvl == "" {
		if v := os.Getenv("LOG_LEVEL"); v != "" {
			lvl = v
		} else {
			lvl = "info"
		}
	}

	parsed, err := logrus.ParseLevel(strings.ToLower(lvl))
	if err != nil {
		Logger.SetLevel(logrus.InfoLevel)
		return err
	}

	Logger.SetLevel(parsed)
	return nil
}

// Writer returns an io.Writer that forwards writes to logrus at appropriate levels.
// It tries to detect level hints in the message (e.g., "[DEBUG]") and route accordingly.
func Writer() io.Writer {
	return logrusWriter{}
}

type logrusWriter struct{}

func (w logrusWriter) Write(p []byte) (n int, err error) {
	if Logger == nil {
		// Fallback to stdout
		return os.Stdout.Write(p)
	}
	msg := strings.TrimSpace(string(p))
	// rudimentary level detection
	switch {
	case strings.HasPrefix(msg, "[DEBUG]") || strings.Contains(msg, "DEBUG"):
		Logger.Debug(strings.TrimPrefix(msg, "[DEBUG]"))
	case strings.HasPrefix(msg, "[WARN]") || strings.Contains(msg, "WARN") || strings.Contains(msg, "WARNING"):
		Logger.Warn(msg)
	case strings.HasPrefix(msg, "[ERROR]") || strings.Contains(msg, "ERROR"):
		Logger.Error(msg)
	default:
		Logger.Info(msg)
	}
	return len(p), nil
}

// Convenience wrappers
func Debugf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, args...)
	}
}
func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, args...)
	}
}
func Warnf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, args...)
	}
}
func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, args...)
	}
}

// WithFields convenience
func WithFields(fields map[string]interface{}) *logrus.Entry {
	if Logger == nil {
		Logger = logrus.New()
	}
	return Logger.WithFields(fields)
}
