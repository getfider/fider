package database

import (
	"strings"

	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

// Logger writes logs to a SQL database
type Logger struct {
	db    *dbx.Database
	level log.Level
	tag   string
}

// NewLogger creates a new Logger
func NewLogger(tag string, db *dbx.Database) *Logger {
	level := strings.ToUpper(env.GetEnvOrDefault("LOG_LEVEL", ""))
	logger := &Logger{tag: tag, db: db}

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

func (l *Logger) log(level log.Level, format string, args ...interface{}) {
	trx, err := l.db.Begin()
	if err != nil {
		//TODO: log somewhere
		return
	}
	defer trx.Commit()
	trx.NoLogs()

	// var count int
	// trx.Scalar(&count, "SELECT COUNT(*) FROM tenants")
	// println("Message: " + string(count))
}
