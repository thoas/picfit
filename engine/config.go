package engine

// Config is the engine config
type Config struct {
	Type          string `mapstructure:"type"`
	DefaultFormat string `mapstructure:"default_format"`
	Format        string `mapstructure:"format"`
	Quality       int    `mapstructure:"quality"`
	MaxBufferSize int    `mapstructure:"max_buffer_size"`
}
