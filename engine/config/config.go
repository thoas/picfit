package config

type Backends struct {
	Lilliput *Backend `mapstructure:"lilliput"`
	GoImage  *Backend `mapstructure:"goimage"`
	Vips     *Backend `mapstructure:"vips"`
}

type Backend struct {
	Mimetypes []string
}

// Config is the engine config
type Config struct {
	Backends        *Backends `mapstructure:"backends"`
	DefaultFormat   string    `mapstructure:"default_format"`
	Format          string    `mapstructure:"format"`
	Quality         int       `mapstructure:"quality"`
	MaxBufferSize   int       `mapstructure:"max_buffer_size"`
	ImageBufferSize int       `mapstructure:"image_buffer_size"`
	JpegQuality     int       `mapstructure:"jpeg_quality"`
	PngCompression  int       `mapstructure:"png_compression"`
	WebpQuality     int       `mapstructure:"webp_quality"`
}
