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
	logger  *stdLog.Logger
	level   log.Level
	tag     string
	props   log.Props
	enabled bool
}

// NewLogger creates a new Logger
func NewLogger(tag string) *Logger {
	level := strings.ToUpper(env.GetEnvOrDefault("LOG_LEVEL", ""))
	return &Logger{
		tag:     tag,
		logger:  stdLog.New(os.Stdout, "", 0),
		level:   log.ParseLevel(level),
		props:   make(log.Props, 0),
		enabled: !env.IsTest(),
	}
}

// Disable logs for this logger
func (l *Logger) Disable() {
	l.enabled = false
}

// Enable logs for this logger
func (l *Logger) Enable() {
	l.enabled = true
}

// SetLevel increases/decreases current log level
func (l *Logger) SetLevel(level log.Level) {
	l.level = level
}

// IsEnabled returns true if given level is enabled
func (l *Logger) IsEnabled(level log.Level) bool {
	return l.enabled && level >= l.level
}

// Debug logs a DEBUG message
func (l *Logger) Debug(message string) {
	l.log(log.DEBUG, message, nil)
}

// Debugf logs a DEBUG message
func (l *Logger) Debugf(message string, props log.Props) {
	l.log(log.DEBUG, message, props)
}

// Info logs a INFO message
func (l *Logger) Info(message string) {
	l.log(log.INFO, message, nil)
}

// Infof logs a INFO message
func (l *Logger) Infof(message string, props log.Props) {
	l.log(log.INFO, message, props)
}

// Warn logs a WARN message
func (l *Logger) Warn(message string) {
	l.log(log.WARN, message, nil)
}

// Warnf logs a WARN message
func (l *Logger) Warnf(message string, props log.Props) {
	l.log(log.WARN, message, props)
}

// Errorf logs a ERROR message
func (l *Logger) Errorf(message string, props log.Props) {
	l.log(log.ERROR, message, props)
}

// Error logs a ERROR message
func (l *Logger) Error(err error) {
	if err != nil {
		l.log(log.ERROR, err.Error(), nil)
	} else {
		l.log(log.ERROR, "nil", nil)
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *Logger) Write(p []byte) (int, error) {
	l.Debug(fmt.Sprintf("%s", p))
	return len(p), nil
}

// New returns a copy of current logger with empty context
func (l *Logger) New() log.Logger {
	return NewLogger(l.tag)
}

// SetProperty with given key:value into current logger context
func (l *Logger) SetProperty(key string, value interface{}) {
	l.props[key] = value
}

func (l *Logger) log(level log.Level, message string, props log.Props) {
	props = l.props.Merge(props)
	if !l.IsEnabled(level) {
		return
	}
	message = log.Parse(message, props, true)
	contextID := l.props[log.PropertyKeyContextID]
	l.logger.Printf("%s [%s] [%s] [%s] %s\n", colorizeLevel(level), time.Now().Format(time.RFC3339), l.tag, contextID, message)
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
