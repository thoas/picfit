package picfit

// Option is a functional option.
type Option func(*Options)

// Options are server options.
type Options struct {
	Load bool
}

// NewOptions initializes server options.
func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// WithLoad overrides load value.
func WithLoad(load bool) Option {
	return func(o *Options) {
		o.Load = load
	}
}
