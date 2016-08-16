package kvstore

import (
	"context"

	"github.com/thoas/gokvstores"
)

const key = "kvstore"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the KVStore associated with this context.
func FromContext(c context.Context) gokvstores.KVStore {
	return c.Value(key).(gokvstores.KVStore)
}

// ToContext adds the KVStore to this context if it supports
// the Setter interface.
func ToContext(c Setter, s gokvstores.KVStore) {
	c.Set(key, s)
}

// NewContext instantiate a new context with a kvstore
func NewContext(ctx context.Context, s gokvstores.KVStore) context.Context {
	return context.WithValue(ctx, key, s)
}
