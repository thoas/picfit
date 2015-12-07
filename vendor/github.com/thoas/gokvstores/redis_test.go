package gokvstores

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func compareStringSets(a []interface{}, b []string) bool {
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	ma := make(map[string]bool)
	for _, aa := range a {
		s, err := String(aa)
		if err != nil {
			return false
		}
		ma[s] = true
	}

	for _, s := range b {
		if _, ok := ma[s]; !ok {
			return false
		}
	}
	return true
}

func TestRedis(t *testing.T) {
	kvstore := NewRedisKVStore("127.0.0.1", 6379, "", 0)

	con := kvstore.Connection()
	defer con.Close()

	con.Flush()

	con.Set("key", "value")

	value, _ := String(con.Get("key"))

	assert.Equal(t, "value", value)

	assert.True(t, con.Exists("key"))

	con.Delete("key")

	assert.Equal(t, nil, con.Get("key"))

	assert.False(t, con.Exists("key"))

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

	// Append
	con.Set("greetings", "Hello, ")

	con.Append("greetings", "World!")
	value, _ = String(con.Get("greetings"))
	assert.Equal(t, "Hello, World!", value)

	con.Append("greetings", " 123")
	value, _ = String(con.Get("greetings"))
	assert.Equal(t, "Hello, World! 123", value)
}
