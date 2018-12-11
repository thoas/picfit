package middleware

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetParamsFromURLValues(t *testing.T) {
	v := url.Values{}
	v.Set("w", "100")
	v.Set("h", "123")
	v.Set("op", "resize")

	params := setParamsFromURLValues(make(map[string]interface{}), v)

	assert.Equal(t, params["w"].(string), "100")
	assert.Equal(t, params["h"].(string), "123")
	assert.Equal(t, params["op"].(string), "resize")

	v.Add("op", "rotate")
	v.Add("w", "99")
	v.Add("h", "321")

	params = setParamsFromURLValues(make(map[string]interface{}), v)

	assert.Equal(t, params["w"].(string), "100")
	assert.Equal(t, params["h"].(string), "123")
	assert.Equal(t, len(params["op"].([]string)), 2)
	assert.Equal(t, params["op"].([]string)[0], "resize")
	assert.Equal(t, params["op"].([]string)[1], "rotate")
}
