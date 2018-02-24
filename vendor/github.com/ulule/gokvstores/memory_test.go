package gokvstores

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	store, err := NewMemoryStore(time.Second*10, time.Second*10)
	assert.Nil(t, err)

	testStore(t, store)
}
