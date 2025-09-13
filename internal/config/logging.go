package config

type Logging struct {
	Level     LoggingLevel  `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format    LoggingFormat `mapstructure:"format" validate:"required,oneof=json text"`
	File      string        `mapstructure:"file"`
	ErrorFile string        `mapstructure:"error_file"`
}

type LoggingLevel string

const (
	LoggingLevelDebug LoggingLevel = "debug"
	LoggingLevelInfo  LoggingLevel = "info"
	LoggingLevelWarn  LoggingLevel = "warn"
	LoggingLevelError LoggingLevel = "error"
)

func (l LoggingLevel) String() string {
	return string(l)
}

type LoggingFormat string

const (
	LoggingFormatJSON LoggingFormat = "json"
	LoggingFormatText LoggingFormat = "text"
)

func (l LoggingFormat) String() string {
	return string(l)
}
