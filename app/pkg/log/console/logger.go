package console

import (
	"fmt"
	stdLog "log"
	"os"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

// Logger output messages to console
type Logger struct {
	logger *stdLog.Logger
	level  log.Level
	tag    string
}

// NewLogger creates a new Logger
func NewLogger(tag string) *Logger {
	level := strings.ToUpper(env.GetEnvOrDefault("LOG_LEVEL", ""))
	return &Logger{
		tag:    tag,
		logger: stdLog.New(os.Stdout, "", 0),
		level:  log.ParseLevel(level),
	}
}

// SetLevel increases/decreases current log level
func (l *Logger) SetLevel(level log.Level) {
	l.level = level
}

// IsEnabled returns true if given level is enabled
func (l *Logger) IsEnabled(level log.Level) bool {
	return level >= l.level
}

// Debugf logs a DEBUG message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(log.DEBUG, format, args...)
}

// Infof logs a INFO message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(log.INFO, format, args...)
}

// Warnf logs a WARN message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(log.WARN, format, args...)
}

// Errorf logs a ERROR message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(log.ERROR, format, args...)
}

// Error logs a ERROR message
func (l *Logger) Error(err error) {
	if err != nil {
		l.log(log.ERROR, err.Error())
	} else {
		l.log(log.ERROR, "nil")
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *Logger) Write(p []byte) (int, error) {
	l.Debugf("%s", p)
	return len(p), nil
}

// New returns a copy of current logger with empty context
func (l *Logger) New() log.Logger {
	return NewLogger(l.tag)
}

func (l *Logger) log(level log.Level, format string, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	message := ""
	if format == "" {
		message = fmt.Sprint(args...)
	} else {
		message = fmt.Sprintf(format, args...)
	}

	l.logger.Printf("%s [%s] [%s] %s\n", colorizeLevel(level), time.Now().Format(time.RFC3339), l.tag, message)
}

func colorizeLevel(level log.Level) string {
	var color func(interface{}) string

	switch level {
	case log.DEBUG:
		color = log.Magenta
	case log.INFO:
		color = log.Blue
	case log.WARN:
		color = log.Yellow
	case log.ERROR:
		color = log.Red
	default:
		color = log.Red
	}

	return log.Bold(color(level.String()))
}
