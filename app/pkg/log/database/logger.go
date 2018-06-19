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
}

// NewLogger creates a new Logger
func NewLogger(tag string, db *dbx.Database) *Logger {
	level := strings.ToUpper(env.GetEnvOrDefault("LOG_LEVEL", ""))
	logger := &Logger{
		tag:     tag,
		db:      db,
		console: console.NewLogger(tag),
	}

	switch level {
	case "DEBUG":
		logger.SetLevel(log.DEBUG)
	case "WARN":
		logger.SetLevel(log.WARN)
	case "ERROR":
		logger.SetLevel(log.ERROR)
	default:
		logger.SetLevel(log.INFO)
	}
	return logger
}

// SetLevel increases/decreases current log level
func (l *Logger) SetLevel(level log.Level) {
	l.level = level
	l.console.SetLevel(level)
}

// IsEnabled returns true if given level is enabled
func (l *Logger) IsEnabled(level log.Level) bool {
	return level >= l.level
}

// Debugf logs a DEBUG message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(log.DEBUG, format, args...)
	l.console.Debugf(format, args...)
}

// Infof logs a INFO message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(log.INFO, format, args...)
	l.console.Infof(format, args...)
}

// Warnf logs a WARN message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(log.WARN, format, args...)
	l.console.Warnf(format, args...)
}

// Errorf logs a ERROR message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(log.ERROR, format, args...)
	l.console.Errorf(format, args...)
}

// Error logs a ERROR message
func (l *Logger) Error(err error) {
	if err != nil {
		l.log(log.ERROR, err.Error())
	} else {
		l.log(log.ERROR, "nil")
	}
	l.console.Error(err)
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *Logger) Write(p []byte) (int, error) {
	l.console.Write(p)
	l.Debugf("%s", p)
	return len(p), nil
}

func (l *Logger) log(level log.Level, format string, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	go func() {
		trx, err := l.db.Begin()
		if err != nil {
			l.console.Error(errors.Wrap(err, "failed to open transaction"))
			return
		}
		trx.NoLogs()

		message := ""
		if format == "" {
			message = fmt.Sprint(args...)
		} else {
			message = fmt.Sprintf(format, args...)
		}

		_, err = trx.Execute(
			"INSERT INTO logs (tag, level, text, created_on) VALUES ($1, $2, $3, $4)",
			l.tag, level.String(), message, time.Now(),
		)

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
	}()
}
