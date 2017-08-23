package lapi

const (
	LOG_ERROR = 0
	LOG_INFO  = 1
	LOG_DEBUG = 2
	LOG_WARN  = 3
)

// Logger controls log
type Logger interface {
	LogWriter
	LogLeveler
}

type LogWriter interface {
	// Write logs a message
	Write(level uint8, message string, args ...interface{}) error
}

type LogLeveler interface {
	// Error writes error message to log
	Error(message string, args ...interface{}) error

	// Debug writes debug message to log
	Debug(message string, args ...interface{}) error

	// Warn writes warn message to log
	Warn(message string, args ...interface{}) error

	// Info writes info message to log
	Info(message string, args ...interface{}) error
}
