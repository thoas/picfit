package engine

import (
	"context"
)

const key = "engine"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the engine associated with this context.
func FromContext(c context.Context) Engine {
	return c.Value(key).(Engine)
}

// ToContext adds the Managers to this context if it supports
// the Setter interface.
func ToContext(c Setter, e Engine) {
	c.Set(key, e)
}

// NewContext instantiate a new context with an engine
func NewContext(ctx context.Context, e Engine) context.Context {
	return context.WithValue(ctx, key, e)
}
