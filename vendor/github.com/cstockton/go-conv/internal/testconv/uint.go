package testconv

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func RunUintTests(t *testing.T, fn func(interface{}) (uint, error)) {
	RunTest(t, reflect.Uint, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunUint8Tests(t *testing.T, fn func(interface{}) (uint8, error)) {
	RunTest(t, reflect.Uint8, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunUint16Tests(t *testing.T, fn func(interface{}) (uint16, error)) {
	RunTest(t, reflect.Uint16, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunUint32Tests(t *testing.T, fn func(interface{}) (uint32, error)) {
	RunTest(t, reflect.Uint32, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

func RunUint64Tests(t *testing.T, fn func(interface{}) (uint64, error)) {
	RunTest(t, reflect.Uint64, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

type testUint64Converter uint64

func (t testUint64Converter) Uint64() (uint64, error) {
	return uint64(t) + 5, nil
}

func init() {
	type ulyUint uint
	type ulyUint8 uint8
	type ulyUint16 uint8
	type ulyUint32 uint8
	type ulyUint64 uint64

	exp := func(e uint, e8 uint8, e16 uint16, e32 uint32, e64 uint64) []Expecter {
		return []Expecter{Exp{e}, Exp{e8}, Exp{e16}, Exp{e32}, Exp{e64}}
	}
	experrs := func(s string) []Expecter {
		return []Expecter{
			experr(uint(0), s), experr(uint8(0), s), experr(uint16(0), s),
			experr(uint32(0), s), experr(uint64(0), s)}
	}

	// basics
	assert(0, exp(0, 0, 0, 0, 0))
	assert(1, exp(1, 1, 1, 1, 1))
	assert(false, exp(0, 0, 0, 0, 0))
	assert(true, exp(1, 1, 1, 1, 1))
	assert("false", exp(0, 0, 0, 0, 0))
	assert("true", exp(1, 1, 1, 1, 1))

	// test length kinds
	assert([]string{"one", "two"}, exp(2, 2, 2, 2, 2))
	assert(map[int]string{1: "one", 2: "two"}, exp(2, 2, 2, 2, 2))

	// test implements Uint64(uint64, error)
	assert(testUint64Converter(5), exp(10, 10, 10, 10, 10))

	// max bounds
	assert(math.MaxUint8, exp(math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8))
	assert(math.MaxUint16, exp(math.MaxUint16, math.MaxUint8, math.MaxUint16,
		math.MaxUint16, math.MaxUint16))
	assert(math.MaxUint32, exp(math.MaxUint32, math.MaxUint8, math.MaxUint16,
		math.MaxUint32, math.MaxUint32))
	assert(uint64(math.MaxUint64), exp(MaxUint, math.MaxUint8,
		math.MaxUint16, math.MaxUint32, uint64(math.MaxUint64)))

	// min bounds
	assert(math.MinInt8, exp(0, 0, 0, 0, 0))
	assert(math.MinInt16, exp(0, 0, 0, 0, 0))
	assert(math.MinInt32, exp(0, 0, 0, 0, 0))
	assert(int64(math.MinInt64), exp(0, 0, 0, 0, 0))

	// perms of various type
	for n := uint(0); n < math.MaxUint8; n += 0xB {
		i := n

		// uints
		if n < 1 {
			i = 0
		} else {
			assert(uintptr(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(i, i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(uint8(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(uint16(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(uint32(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(uint64(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		}

		// underlying
		assert(ulyUint(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(ulyUint8(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(ulyUint16(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(ulyUint32(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(ulyUint64(i), i, uint8(i), uint16(i), uint32(i), uint64(i))

		// implements
		if i < math.MaxUint8-5 {
			assert(testUint64Converter(i),
				i+5, uint8(i+5), uint16(i+5), uint32(i+5), uint64(i+5))
			assert(testUint64Converter(ulyUint(i)),
				i+5, uint8(i+5), uint16(i+5), uint32(i+5), uint64(i+5))
		}

		// ints
		if i < math.MaxInt8 {
			assert(int(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(int8(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(int16(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(int32(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
			assert(int64(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		}

		// floats
		assert(float32(i), i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(float64(i), i, uint8(i), uint16(i), uint32(i), uint64(i))

		// complex
		assert(complex(float32(i), 0),
			i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(complex(float64(i), 0),
			i, uint8(i), uint16(i), uint32(i), uint64(i))

		// from string int
		assert(fmt.Sprintf("%d", i),
			i, uint8(i), uint16(i), uint32(i), uint64(i))
		assert(testStringConverter(fmt.Sprintf("%d", i)),
			i, uint8(i), uint16(i), uint32(i), uint64(i))

		// from string float form
		assert(fmt.Sprintf("%d.0", i),
			i, uint8(i), uint16(i), uint32(i), uint64(i))
	}

	assert(nil, experrs(`cannot convert <nil> (type <nil>) to `))
	assert("foo", experrs(` "foo" (type string) `))
	assert(struct{}{}, experrs(`cannot convert struct {}{} (type struct {}) to `))
}
