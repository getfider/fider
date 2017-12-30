package log

import (
	"fmt"
	"time"
)

// Level defines all possible log levels
type Level uint8

const (
	// DEBUG for verbose logs
	DEBUG Level = iota + 1
	// INFO for WARN+ERROR+INFO logs
	INFO
	// WARN for WARN+ERROR logs
	WARN
	// ERROR for ERROR only logs
	ERROR
)

// Logger defines the logging interface.
type Logger interface {
	SetLevel(level Level)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Error(err error)
	Write(p []byte) (int, error)
}

// ConsoleLogger output messages to console
type ConsoleLogger struct {
	level Level
}

// NewConsoleLogger creates a new ConsoleLogger
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

// SetLevel increases/decreases current log level
func (l *ConsoleLogger) SetLevel(level Level) {
	l.level = level
}

// Debugf logs a DEBUG message
func (l *ConsoleLogger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Infof logs a INFO message
func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warnf logs a WARN message
func (l *ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Errorf logs a ERROR message
func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Error logs a ERROR message
func (l *ConsoleLogger) Error(err error) {
	if err != nil {
		l.log(ERROR, err.Error())
	} else {
		l.log(ERROR, "nil")
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *ConsoleLogger) Write(p []byte) (int, error) {
	l.Debugf("%s", p)
	return len(p), nil
}

func (l *ConsoleLogger) log(level Level, format string, args ...interface{}) {
	if level >= l.level {
		message := ""
		if format == "" {
			message = fmt.Sprint(args...)
		} else {
			message = fmt.Sprintf(format, args...)
		}

		fmt.Printf("%s [%s] %s\n", levelString(level), time.Now().Format(time.RFC3339), message)
	}
}

func levelString(level Level) string {
	switch level {
	case DEBUG:
		return Bold(Magenta("DEBUG"))
	case INFO:
		return Bold(Blue("INFO"))
	case WARN:
		return Bold(Yellow("WARN"))
	case ERROR:
		return Bold(Red("ERROR"))
	default:
		return Bold(Red("???"))
	}
}
