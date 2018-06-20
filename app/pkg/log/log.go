package log

import (
	"fmt"
	"regexp"
	"strings"
)

// Level defines all possible log levels
type Level uint8

// Logger defines the logging interface.
type Logger interface {
	SetLevel(level Level)
	SetProperty(key string, value interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Error(err error)
	IsEnabled(level Level) bool
	Write(p []byte) (int, error)
	New() Logger
}

const (
	// PropertyKeyContextID is the context id of current logger
	PropertyKeyContextID = "context_id"
	// PropertyKeyUserID is the user id of current logger
	PropertyKeyUserID = "user_id"
	// PropertyKeyTenantID is the tenant id of current logger
	PropertyKeyTenantID = "tenant_id"
)

const (
	// DEBUG for verbose logs
	DEBUG Level = iota + 1
	// INFO for WARN+ERROR+INFO logs
	INFO
	// WARN for WARN+ERROR logs
	WARN
	// ERROR for ERROR only logs
	ERROR
	// NONE is used to disable logs
	NONE
)

// ParseLevel returns a log.Level based on input string
func ParseLevel(level string) Level {
	switch level {
	case "DEBUG":
		return DEBUG
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case NONE:
		return "NONE"
	default:
		return "???"
	}
}

// Props is a map of key:value
type Props map[string]interface{}

var placeholderFinder = regexp.MustCompile("@{.*?}")

// Parse is used to merge props into format and return a text message
func Parse(format string, props Props) string {
	if props == nil || len(props) == 0 {
		return format
	}

	for {
		indexes := placeholderFinder.FindSubmatchIndex([]byte(format))
		if len(indexes) == 0 {
			return format
		}

		ph := format[indexes[0]:indexes[1]]
		phContent := ph[2 : len(ph)-1]
		phSeparatorIdx := strings.Index(phContent, ":")
		value := props[phContent]
		if phSeparatorIdx >= 0 {
			phName := phContent[:phSeparatorIdx]
			phColor := phContent[phSeparatorIdx+1:]
			value = Color(phColor, props[phName])
		}
		format = fmt.Sprintf("%s%v%s", format[:indexes[0]], value, format[indexes[1]:])
	}
}
