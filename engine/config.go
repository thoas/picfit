package engine

// Config is the engine config
type Config struct {
	Backends      []string `mapstructure:"backends"`
	DefaultFormat string   `mapstructure:"default_format"`
	Format        string   `mapstructure:"format"`
	Quality       int      `mapstructure:"quality"`
	MaxBufferSize int      `mapstructure:"max_buffer_size"`
}
