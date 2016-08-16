package storage

import (
	"context"

	"github.com/thoas/gostorages"
)

const sourceKey = "srcStorage"
const destinationKey = "dstStorage"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// SourceFromContext returns the Storage associated with this context.
func SourceFromContext(c context.Context) gostorages.Storage {
	return c.Value(sourceKey).(gostorages.Storage)
}

// SourceToContext adds the Storage to this context if it supports
// the Setter interface.
func SourceToContext(c Setter, s gostorages.Storage) {
	c.Set(sourceKey, s)
}

// NewSourceContext instantiate a new context with a storage
func NewSourceContext(ctx context.Context, s gostorages.Storage) context.Context {
	return context.WithValue(ctx, sourceKey, s)
}

// DestinationFromContext returns the Storage associated with this context.
func DestinationFromContext(c context.Context) gostorages.Storage {
	return c.Value(destinationKey).(gostorages.Storage)
}

// DestinationToContext adds the Storage to this context if it supports
// the Setter interface.
func DestinationToContext(c Setter, s gostorages.Storage) {
	c.Set(destinationKey, s)
}

// NewDestinationContext instantiate a new context with a storage
func NewDestinationContext(ctx context.Context, s gostorages.Storage) context.Context {
	return context.WithValue(ctx, destinationKey, s)
}
