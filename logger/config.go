package logger

const (
	DevelopmentLevel = "development"
	ProductionLevel  = "production"

	// defaultLevel is the default logger level
	defaultLevel = DevelopmentLevel
)

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
