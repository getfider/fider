package log

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Level defines all possible log levels
type Level uint8

// Logger defines the logging interface.
type Logger interface {
	Disable()
	Enable()
	SetLevel(level Level)
	SetProperty(key string, value interface{})
	Debug(message string)
	Debugf(message string, props Props)
	Info(message string)
	Infof(message string, props Props)
	Warn(format string)
	Warnf(message string, props Props)
	Error(err error)
	Errorf(message string, props Props)
	IsEnabled(level Level) bool
	Write(p []byte) (int, error)
	New() Logger
}

const (
	// PropertyKeySessionID is the session id of current logger
	PropertyKeySessionID = "SessionID"
	// PropertyKeyContextID is the context id of current logger
	PropertyKeyContextID = "ContextID"
	// PropertyKeyUserID is the user id of current logger
	PropertyKeyUserID = "UserID"
	// PropertyKeyTenantID is the tenant id of current logger
	PropertyKeyTenantID = "TenantID"
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

// Value converts props into a database value
func (p Props) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Merge current props with given props
func (p Props) Merge(props Props) Props {
	new := Props{}
	if p != nil {
		for k, v := range p {
			new[k] = v
		}
	}
	if props != nil {
		for k, v := range props {
			new[k] = v
		}
	}
	return new
}

var placeholderFinder = regexp.MustCompile("@{.*?}")

// Parse is used to merge props into format and return a text message
func Parse(format string, props Props, colorize bool) string {
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
			value = props[phName]
			if colorize {
				value = Color(phColor, value)
			}
		}
		format = fmt.Sprintf("%s%v%s", format[:indexes[0]], value, format[indexes[1]:])
	}
}
