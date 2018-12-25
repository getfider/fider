package noop

import (
	"github.com/getfider/fider/app/pkg/log"
)

// Logger doesn't log anything
type Logger struct {
	enabled bool
}

// NewLogger creates a new Logger
func NewLogger() *Logger {
	return &Logger{}
}

// Disable logs for this logger
func (l *Logger) Disable() {
}

// Enable logs for this logger
func (l *Logger) Enable() {
}

// SetLevel increases/decreases current log level
func (l *Logger) SetLevel(level log.Level) {
}

// IsEnabled returns true if given level is enabled
func (l *Logger) IsEnabled(level log.Level) bool {
	return true
}

// Debug logs a DEBUG message
func (l *Logger) Debug(message string) {
}

// Debugf logs a DEBUG message
func (l *Logger) Debugf(message string, props log.Props) {
}

// Info logs a INFO message
func (l *Logger) Info(message string) {
}

// Infof logs a INFO message
func (l *Logger) Infof(message string, props log.Props) {
}

// Warn logs a WARN message
func (l *Logger) Warn(message string) {
}

// Warnf logs a WARN message
func (l *Logger) Warnf(message string, props log.Props) {
}

// Error logs a ERROR message
func (l *Logger) Error(err error) {
}

// Errorf logs a ERROR message
func (l *Logger) Errorf(message string, props log.Props) {
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *Logger) Write(p []byte) (int, error) {
	return 0, nil
}

// New returns a copy of current logger with empty context
func (l *Logger) New() log.Logger {
	return NewLogger()
}

// SetProperty with given key:value into current logger context
func (l *Logger) SetProperty(key string, value interface{}) {
}

func (l *Logger) log(level log.Level, format string, args ...interface{}) {
}
