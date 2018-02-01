package log

// NoopLogger doesn't log anything
type NoopLogger struct {
}

// SetLevel increases/decreases current log level
func (l *NoopLogger) SetLevel(level Level) {
}

// IsEnabled returns true if given level is enabled
func (l *NoopLogger) IsEnabled(level Level) bool {
	return false
}

// Debugf logs a DEBUG message
func (l *NoopLogger) Debugf(format string, args ...interface{}) {
}

// Infof logs a INFO message
func (l *NoopLogger) Infof(format string, args ...interface{}) {
}

// Warnf logs a WARN message
func (l *NoopLogger) Warnf(format string, args ...interface{}) {
}

// Errorf logs a ERROR message
func (l *NoopLogger) Errorf(format string, args ...interface{}) {
}

// Error logs a ERROR message
func (l *NoopLogger) Error(err error) {
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *NoopLogger) Write(p []byte) (int, error) {
	return 0, nil
}

func (l *NoopLogger) log(level Level, format string, args ...interface{}) {
}
