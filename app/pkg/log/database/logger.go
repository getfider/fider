package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log/console"

	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

// Logger writes logs to a SQL database
type Logger struct {
	db      *dbx.Database
	console log.Logger
	level   log.Level
	tag     string
	props   log.Props
	enabled bool
}

// NewLogger creates a new Logger
func NewLogger(tag string, db *dbx.Database) *Logger {
	level := strings.ToUpper(env.GetEnvOrDefault("LOG_LEVEL", ""))
	return &Logger{
		tag:     tag,
		db:      db,
		console: console.NewLogger(tag),
		level:   log.ParseLevel(level),
		props:   make(log.Props, 0),
		enabled: true,
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
	l.console.SetLevel(level)
}

// IsEnabled returns true if given level is enabled
func (l *Logger) IsEnabled(level log.Level) bool {
	return l.enabled && level >= l.level
}

// Debug logs a DEBUG message
func (l *Logger) Debug(message string) {
	l.console.Debug(message)
	l.log(log.DEBUG, message, nil)
}

// Debugf logs a DEBUG message
func (l *Logger) Debugf(message string, props log.Props) {
	l.console.Debugf(message, props)
	l.log(log.DEBUG, message, props)
}

// Info logs a INFO message
func (l *Logger) Info(message string) {
	l.console.Info(message)
	l.log(log.INFO, message, nil)
}

// Infof logs a INFO message
func (l *Logger) Infof(message string, props log.Props) {
	l.console.Infof(message, props)
	l.log(log.INFO, message, props)
}

// Warn logs a WARN message
func (l *Logger) Warn(message string) {
	l.console.Warn(message)
	l.log(log.WARN, message, nil)
}

// Warnf logs a WARN message
func (l *Logger) Warnf(message string, props log.Props) {
	l.console.Warnf(message, props)
	l.log(log.WARN, message, props)
}

// Errorf logs a ERROR message
func (l *Logger) Errorf(message string, props log.Props) {
	l.console.Errorf(message, props)
	l.log(log.ERROR, message, props)
}

// Error logs a ERROR message
func (l *Logger) Error(err error) {
	if err != nil {
		l.log(log.ERROR, err.Error(), nil)
	} else {
		l.log(log.ERROR, "nil", nil)
	}
	l.console.Error(err)
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *Logger) Write(p []byte) (int, error) {
	l.console.Write(p)
	l.Debug(fmt.Sprintf("%s", p))
	return len(p), nil
}

// New returns a copy of current logger with empty context
func (l *Logger) New() log.Logger {
	return NewLogger(l.tag, l.db)
}

// SetProperty with given key:value into current logger context
func (l *Logger) SetProperty(key string, value interface{}) {
	l.props[key] = value
	l.console.SetProperty(key, value)
}

func (l *Logger) log(level log.Level, message string, props log.Props) {
	if !l.IsEnabled(level) {
		return
	}

	props = l.props.Merge(props)
	trx, err := l.db.Begin()
	if err != nil {
		l.console.Error(errors.Wrap(err, "failed to open transaction"))
		return
	}

	message = log.Parse(message, props, false)

	trx.NoLogs()
	_, err = trx.Execute(
		"INSERT INTO logs (tag, level, text, created_at, properties) VALUES ($1, $2, $3, $4, $5)",
		l.tag, level.String(), message, time.Now(), props,
	)
	trx.ResumeLogs()

	if err != nil {
		l.console.Error(errors.Wrap(err, "failed to insert log"))
		err = trx.Rollback()
		if err != nil {
			l.console.Error(errors.Wrap(err, "failed to rollback transaction"))
		}
	} else {
		err = trx.Commit()
		if err != nil {
			l.console.Error(errors.Wrap(err, "failed to commit transaction"))
		}
	}
}
