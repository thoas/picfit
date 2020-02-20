package logger

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarningLevel = "warn"
	ErrorLevel   = "error"
	FatalLevel   = "fatal"

	// defaultLevel is the default logger level
	defaultLevel = DebugLevel
)

const (
	ConsoleType    = "console"
	JsonType       = "json"
	HowdooJsonType = "howdoo_json"

	// defaultType is the default logger type
	defaultType = ConsoleType
)

// Config is a struct to configure logger
type Config struct {
	App     string
	Level   string
	Type    string
	Channel string
}

// GetLevel returns the level of the logger
func (l *Config) GetLevel() string {
	if l.Level == "" {
		return defaultLevel
	}

	return l.Level
}

// GetType returns the type of the logger
func (l *Config) GetType() string {
	if l.Type == "" {
		return defaultType
	}

	return l.Type
}
