package logger

import "context"

const key = "logger"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the KVStore associated with this context.
func FromContext(c context.Context) Logger {
	return c.Value(key).(Logger)
}

// ToContext adds a Logger to this context if it supports
// the Setter interface.
func ToContext(c Setter, l Logger) {
	c.Set(key, l)
}

// NewContext instantiate a new context with a kvstore
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, key, l)
}
