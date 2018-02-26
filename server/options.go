package server

import "context"

// Option is a functional option.
type Option func(*Options)

// Options are server options.
type Options struct {
	Context context.Context
}

// NewOptions initializes server options.
func NewOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// WithContext overrides context instance.
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
