package gokvstores

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	kvstore := NewCacheKVStore(100)

	con := kvstore.Connection()
	defer con.Close()

	con.Set("key", "value")

	value, _ := String(con.Get("key"))

	assert.Equal(t, "value", value)

	assert.True(t, con.Exists("key"))

	con.Delete("key")

	assert.Equal(t, nil, con.Get("key"))

	assert.False(t, con.Exists("key"))

	// Append
	con.Set("greetings", "Hello, ")

	con.Append("greetings", "World!")
	value, _ = String(con.Get("greetings"))
	assert.Equal(t, "Hello, World!", value)

	con.Append("greetings", " 123")
	value, _ = String(con.Get("greetings"))
	assert.Equal(t, "Hello, World! 123", value)

	// Sets
	con.SetAdd("myset", "hello")
	con.SetAdd("myset", "world")
	assert.True(t, compareStringSets(con.SetMembers("myset"), []string{"hello", "world"}))
	con.SetAdd("myset", "hello")
	assert.True(t, compareStringSets(con.SetMembers("myset"), []string{"hello", "world"}))
	con.SetAdd("myset", "hi")
	assert.True(t, compareStringSets(con.SetMembers("myset"), []string{"hello", "world", "hi"}))
	con.Delete("myset")
	assert.True(t, con.SetMembers("myset") == nil)
}
