package log

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
	// NONE is used to disable logs
	NONE
)

// parseLevel returns a log.Level based on input string
func parseLevel(level string) Level {
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
