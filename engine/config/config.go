package config

type Backends struct {
	Gifsicle *CommandBackend `mapstructure:"gifsicle"`
	GoImage  *Backend        `mapstructure:"goimage"`
	Vips     *VipsBackend    `mapstructure:"vips"`
}

type VipsBackend struct {
	Weight    int
	Mimetypes []string
	// Concurrency is the number of libvips worker threads per image, 0 uses the libvips default.
	Concurrency int
}

type Backend struct {
	Weight    int
	Mimetypes []string
}

type CommandBackend struct {
	Path      string
	Mimetypes []string
	Weight    int
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
