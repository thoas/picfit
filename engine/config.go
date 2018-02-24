package engine

// Config is the engine config
type Config struct {
	DefaultFormat string `mapstructure:"default_format"`
	Format        string `mapstructure:"format"`
	Quality       int    `mapstructure:"quality"`
}
