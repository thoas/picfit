package gokvstores

import (
	"sort"
	"testing"
	"time"

	conv "github.com/cstockton/go-conv"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore(t *testing.T) {
	store, err := NewRedisClientStore(&RedisClientOptions{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}, time.Second*30)

	assert.Nil(t, err)

	testStore(t, store)

	is := assert.New(t)

	mapResults := map[string]map[string]interface{}{
		"order1": {"language": "go"},
		"order2": {"integer": "1"},
		"order3": {"float": "20.2"},
	}
	expectedStrings := []string{"order1", "order2", "order3"}

	for key, expected := range mapResults {
		err = store.SetMap(key, expected)
		is.NoError(err)
	}

	values, err := store.Keys("order*")
	is.NoError(err)

	sort.Strings(expectedStrings)
	result := make([]string, len(values))
	for k, v := range values {
		result[k], _ = conv.String(v)
	}
	sort.Strings(result)

	is.Equal(expectedStrings, result)

	assert.Nil(t, store.Close())
}
