package testconv

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func RunFloat32Tests(t *testing.T, fn func(interface{}) (float32, error)) {
	RunTest(t, reflect.Float32, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunFloat64Tests(t *testing.T, fn func(interface{}) (float64, error)) {
	RunTest(t, reflect.Float64, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

type testFloat64Converter float64

func (t testFloat64Converter) Float64() (float64, error) {
	return float64(t) + 5, nil
}

func init() {
	type ulyFloat32 float32
	type ulyFloat64 float64

	exp := func(e32 float32, e64 float64) []Expecter {
		return []Expecter{Float32Exp{e32}, Float64Exp{e64}}
	}
	experrs := func(s string) []Expecter {
		return []Expecter{experr(float32(0), s), experr(float64(0), s)}
	}

	// basics
	assert(0, exp(0, 0))
	assert(1, exp(1, 1))
	assert(false, exp(0, 0))
	assert(true, exp(1, 1))
	assert("false", exp(0, 0))
	assert("true", exp(1, 1))

	// test length kinds
	assert([]string{"one", "two"}, exp(2, 2))
	assert(map[int]string{1: "one", 2: "two"}, exp(2, 2))

	// test implements Float64(float64, error)
	assert(testFloat64Converter(5), exp(10, 10))

	// max bounds
	assert(math.MaxFloat32, exp(math.MaxFloat32, math.MaxFloat32))
	assert(math.MaxFloat64, exp(math.MaxFloat32, math.MaxFloat64))

	// min bounds
	assert(-math.MaxFloat32, exp(-math.MaxFloat32, -math.MaxFloat32))
	assert(-math.MaxFloat64, exp(-math.MaxFloat32, -math.MaxFloat64))

	// ints
	assert(int(10), exp(10, float64(10)))
	assert(int8(10), exp(10, float64(10)))
	assert(int16(10), exp(10, float64(10)))
	assert(int32(10), exp(10, float64(10)))
	assert(int64(10), exp(10, float64(10)))

	// uints
	assert(uint(10), exp(10, float64(10)))
	assert(uint8(10), exp(10, float64(10)))
	assert(uint16(10), exp(10, float64(10)))
	assert(uint32(10), exp(10, float64(10)))
	assert(uint64(10), exp(10, float64(10)))

	// perms of various type
	for i := float32(-3.0); i < 3.0; i += .5 {

		// underlying
		assert(ulyFloat32(i), exp(i, float64(i)))
		assert(ulyFloat64(i), exp(i, float64(i)))

		// implements
		assert(testFloat64Converter(i), exp(i+5, float64(i+5)))
		assert(testFloat64Converter(ulyFloat64(i)), exp(i+5, float64(i+5)))

		// floats
		assert(i, exp(i, float64(i)))
		assert(float64(i), exp(i, float64(i)))

		// complex
		assert(complex(i, 0), exp(i, float64(i)))
		assert(complex(float64(i), 0), exp(i, float64(i)))

		// from string int
		assert(fmt.Sprintf("%#v", i), exp(i, float64(i)))
		assert(testStringConverter(fmt.Sprintf("%#v", i)), exp(i, float64(i)))

		// from string float form
		assert(fmt.Sprintf("%#v", i), exp(i, float64(i)))
	}

	assert("foo", experrs(`cannot convert "foo" (type string) to `))
	assert(struct{}{}, experrs(`cannot convert struct {}{} (type struct {}) to `))
	assert(nil, experrs(`cannot convert <nil> (type <nil>) to `))

}
