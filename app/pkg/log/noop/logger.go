package noop

import "github.com/getfider/fider/app/pkg/log"

// Logger doesn't log anything
type Logger struct {
}

// NewLogger creates a new Logger
func NewLogger() *Logger {
	return &Logger{}
}

// SetLevel increases/decreases current log level
func (l *Logger) SetLevel(level log.Level) {
}

// IsEnabled returns true if given level is enabled
func (l *Logger) IsEnabled(level log.Level) bool {
	return true
}

// Debugf logs a DEBUG message
func (l *Logger) Debugf(format string, args ...interface{}) {
}

// Infof logs a INFO message
func (l *Logger) Infof(format string, args ...interface{}) {
}

// Warnf logs a WARN message
func (l *Logger) Warnf(format string, args ...interface{}) {
}

// Errorf logs a ERROR message
func (l *Logger) Errorf(format string, args ...interface{}) {
}

// Error logs a ERROR message
func (l *Logger) Error(err error) {
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
