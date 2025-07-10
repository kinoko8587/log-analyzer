package analyzer

import "time"

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelDebug LogLevel = "DEBUG"
)

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
}

func NewLog(level LogLevel, message string) *Log {
	return &Log{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
}

func (l *Log) IsError() bool {
	return l.Level == LogLevelError
}