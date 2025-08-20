package events

import "encoding/json"

type LogLevel string

const (
	LogLevelInfo  LogLevel = "LOG_INFO"
	LogLevelWarn  LogLevel = "LOG_WARN"
	LogLevelError LogLevel = "LOG_ERROR"
	LogLevelDebug LogLevel = "LOG_DEBUG"
)

type LogEvent struct {
	Name   string   `json:"name"`
	Data   string   `json:"data"`
	Level  LogLevel `json:"level"`
	Source string   `json:"source"`
	// TraceID   string    `json:"trace_id,omitempty"`
}

// ToJSON converts the event to JSON bytes.
func (le *LogEvent) ToJSON() ([]byte, error) {
	return json.Marshal(le)
}
