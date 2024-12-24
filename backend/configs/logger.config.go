package configs

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelError LogLevel = "ERROR"
)

type LoggerConfig struct {
	Level LogLevel `json:"log_level"`
}
