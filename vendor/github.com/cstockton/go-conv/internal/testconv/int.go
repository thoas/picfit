package testconv

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func RunIntTests(t *testing.T, fn func(interface{}) (int, error)) {
	RunTest(t, reflect.Int, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunInt8Tests(t *testing.T, fn func(interface{}) (int8, error)) {
	RunTest(t, reflect.Int8, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunInt16Tests(t *testing.T, fn func(interface{}) (int16, error)) {
	RunTest(t, reflect.Int16, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunInt32Tests(t *testing.T, fn func(interface{}) (int32, error)) {
	RunTest(t, reflect.Int32, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunInt64Tests(t *testing.T, fn func(interface{}) (int64, error)) {
	RunTest(t, reflect.Int64, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

// Numeric constants for testing.
const (
	IntSize int  = 32 << (^uint(0) >> 63)
	MinUint uint = 0
	MaxUint uint = ^MinUint
	MaxInt  int  = int(MaxUint >> 1)
	MinInt  int  = ^MaxInt
)

type testInt64Converter int64

func (t testInt64Converter) Int64() (int64, error) {
	return int64(t) + 5, nil
}

func init() {
	type ulyInt int
	type ulyInt8 int8
	type ulyInt16 int8
	type ulyInt32 int8
	type ulyInt64 int64

	exp := func(e int, e8 int8, e16 int16, e32 int32, e64 int64) []Expecter {
		return []Expecter{Exp{e}, Exp{e8}, Exp{e16}, Exp{e32}, Exp{e64}}
	}
	experrs := func(s string) []Expecter {
		return []Expecter{
			experr(int(0), s), experr(int8(0), s), experr(int16(0), s),
			experr(int32(0), s), experr(int64(0), s)}
	}

	// basics
	assert(-1, exp(-1, -1, -1, -1, -1))
	assert(0, exp(0, 0, 0, 0, 0))
	assert(1, exp(1, 1, 1, 1, 1))
	assert(false, exp(0, 0, 0, 0, 0))
	assert(true, exp(1, 1, 1, 1, 1))
	assert("false", exp(0, 0, 0, 0, 0))
	assert("true", exp(1, 1, 1, 1, 1))

	// test length kinds
	assert([]string{"one", "two"}, 2, 2, 2, 2, 2)
	assert(map[int]string{1: "one", 2: "two"}, 2, 2, 2, 2, 2)

	// test implements Int64(int64, error)
	assert(testInt64Converter(5), 10, 10, 10, 10, 10)

	// overflow
	assert(uint64(math.MaxUint64), exp(MaxInt, math.MaxInt8,
		math.MaxInt16, math.MaxInt32, math.MaxInt64))

	// underflow
	assert(int64(math.MinInt64), exp(MinInt, math.MinInt8, math.MinInt16,
		math.MinInt32, math.MinInt64))

	// max bounds
	assert(math.MaxInt8, exp(math.MaxInt8, math.MaxInt8, math.MaxInt8,
		math.MaxInt8, math.MaxInt8))
	assert(math.MaxInt16, exp(math.MaxInt16, math.MaxInt8, math.MaxInt16,
		math.MaxInt16, math.MaxInt16))
	assert(math.MaxInt32, exp(math.MaxInt32, math.MaxInt8, math.MaxInt16,
		math.MaxInt32, math.MaxInt32))
	assert(math.MaxInt64, exp(MaxInt, math.MaxInt8, math.MaxInt16,
		math.MaxInt32, math.MaxInt64))

	// min bounds
	assert(math.MinInt8, exp(math.MinInt8, math.MinInt8, math.MinInt8,
		math.MinInt8, math.MinInt8))
	assert(math.MinInt16, exp(math.MinInt16, math.MinInt8, math.MinInt16,
		math.MinInt16, math.MinInt16))
	assert(math.MinInt32, exp(math.MinInt32, math.MinInt8, math.MinInt16,
		math.MinInt32, math.MinInt32))
	assert(int64(math.MinInt64), exp(MinInt, math.MinInt8, math.MinInt16,
		math.MinInt32, math.MinInt64))

	// perms of various type
	for i := int(math.MinInt8); i < math.MaxInt8; i += 0xB {

		// uints
		if i > 0 {
			assert(uint(i), i, int8(i), int16(i), int32(i), int64(i))
			assert(uint8(i), i, int8(i), int16(i), int32(i), int64(i))
			assert(uint16(i), i, int8(i), int16(i), int32(i), int64(i))
			assert(uint32(i), i, int8(i), int16(i), int32(i), int64(i))
			assert(uint64(i), i, int8(i), int16(i), int32(i), int64(i))
		}

		// underlying
		assert(ulyInt(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(ulyInt8(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(ulyInt16(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(ulyInt32(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(ulyInt64(i), i, int8(i), int16(i), int32(i), int64(i))

		// implements
		if i < math.MaxInt8-5 {
			assert(testInt64Converter(i),
				i+5, int8(i+5), int16(i+5), int32(i+5), int64(i+5))
			assert(testInt64Converter(ulyInt(i)),
				i+5, int8(i+5), int16(i+5), int32(i+5), int64(i+5))
		}

		// ints
		assert(i, i, int8(i), int16(i), int32(i), int64(i))
		assert(int8(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(int16(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(int32(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(int64(i), i, int8(i), int16(i), int32(i), int64(i))

		// floats
		assert(float32(i), i, int8(i), int16(i), int32(i), int64(i))
		assert(float64(i), i, int8(i), int16(i), int32(i), int64(i))

		// complex
		assert(complex(float32(i), 0),
			i, int8(i), int16(i), int32(i), int64(i))
		assert(complex(float64(i), 0),
			i, int8(i), int16(i), int32(i), int64(i))

		// from string int
		assert(fmt.Sprintf("%d", i),
			i, int8(i), int16(i), int32(i), int64(i))
		assert(testStringConverter(fmt.Sprintf("%d", i)),
			i, int8(i), int16(i), int32(i), int64(i))

		// from string float form
		assert(fmt.Sprintf("%d.0", i),
			i, int8(i), int16(i), int32(i), int64(i))
	}

	assert("foo", experrs(`"foo" (type string) `))
	assert(struct{}{}, experrs(`cannot convert struct {}{} (type struct {}) to `))
	assert(nil, experrs(`cannot convert <nil> (type <nil>) to `))
}
