package logger

// defaultLevel is the default logger level
const defaultLevel = "debug"

// Config is a struct to configure logger
type Config struct {
	Level string
}

// GetLevel returns the level of the logger
func (l *Config) GetLevel() string {
	if l.Level == "" {
		return defaultLevel
	}

	return l.Level
}
