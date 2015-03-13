package engines

type Options struct {
	Upscale bool
	Format  string
	Quality int
}

type Engine interface {
	Resize(source []byte, width int, height int, options *Options) ([]byte, error)
	Thumbnail(source []byte, width int, height int, options *Options) ([]byte, error)
}
