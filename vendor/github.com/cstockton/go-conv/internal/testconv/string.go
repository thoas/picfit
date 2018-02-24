package testconv

import (
	"reflect"
	"testing"
)

func RunStringTests(t *testing.T, fn func(interface{}) (string, error)) {
	RunTest(t, reflect.String, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

type testStringConverter string

// @TODO Check Stringer before Converter interface, or after?
func (t testStringConverter) String() (string, error) {
	return string(t) + "Tested", nil
}

func init() {

	// basic
	assert(`hello`, `hello`)
	assert(``, ``)
	assert([]byte(`hello`), `hello`)
	assert([]byte(``), ``)

	// ptr indirection
	assert(new(string), ``)
	assert(new([]byte), ``)

	// underlying string
	type ulyString string
	assert(ulyString(`hello`), `hello`)
	assert(ulyString(``), ``)

	// implements string converter
	assert(testStringConverter(`hello`), `helloTested`)
	assert(testStringConverter(`hello`), `helloTested`)
}
