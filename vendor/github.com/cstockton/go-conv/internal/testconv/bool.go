package testconv

import (
	"math"
	"math/cmplx"
	"reflect"
	"testing"
	"time"
)

func RunBoolTests(t *testing.T, fn func(interface{}) (bool, error)) {
	RunTest(t, reflect.Bool, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

type testBoolConverter bool

func (t testBoolConverter) Bool() (bool, error) {
	return !bool(t), nil
}

func init() {

	// strings: truthy
	trueStrings := []string{
		"1", "t", "T", "true", "True", "TRUE", "y", "Y", "yes", "Yes", "YES"}
	for _, truthy := range trueStrings {
		assert(truthy, true)
		assert(testStringConverter(truthy), true)
	}

	// strings: falsy
	falseStrings := []string{
		"0", "f", "F", "false", "False", "FALSE", "n", "N", "no", "No", "NO"}
	for _, falsy := range falseStrings {
		assert(falsy, false)
		assert(testStringConverter(falsy), false)
	}

	// numerics: true
	for _, i := range []int{-1, 1} {
		assert(i, true)
		assert(int8(i), true)
		assert(int16(i), true)
		assert(int32(i), true)
		assert(int64(i), true)
		assert(uint(i), true)
		assert(uint8(i), true)
		assert(uint16(i), true)
		assert(uint32(i), true)
		assert(uint64(i), true)
		assert(float32(i), true)
		assert(float64(i), true)
		assert(complex(float32(i), 0), true)
		assert(complex(float64(i), 0), true)
	}

	// int/uint: false
	assert(int(0), false)
	assert(int8(0), false)
	assert(int16(0), false)
	assert(int32(0), false)
	assert(int64(0), false)
	assert(uint(0), false)
	assert(uint8(0), false)
	assert(uint16(0), false)
	assert(uint32(0), false)
	assert(uint64(0), false)

	// float: NaN and 0 are false.
	assert(float32(math.NaN()), false)
	assert(math.NaN(), false)
	assert(float32(0), false)
	assert(float64(0), false)

	// complex: NaN and 0 are false.
	assert(complex64(cmplx.NaN()), false)
	assert(cmplx.NaN(), false)
	assert(complex(float32(0), 0), false)
	assert(complex(float64(0), 0), false)

	// time
	assert(time.Time{}, false)
	assert(time.Now(), true)

	// bool
	assert(false, false)
	assert(true, true)

	// underlying bool
	type ulyBool bool
	assert(ulyBool(false), false)
	assert(ulyBool(true), true)

	// implements bool converter
	assert(testBoolConverter(false), true)
	assert(testBoolConverter(true), false)

	// test length kinds
	assert([]string{"one", "two"}, true)
	assert(map[int]string{1: "one", 2: "two"}, true)
	assert([]string{}, false)
	assert([]string(nil), false)

	// errors
	assert(nil, experr(false, `cannot convert <nil> (type <nil>) to bool`))
	assert("foo", experr(false, `cannot parse "foo" (type string) as bool`))
	assert("tooLong", experr(
		false, `cannot parse string with len 7 as bool`))
	assert(struct{}{}, experr(
		false, `cannot convert struct {}{} (type struct {}) to `))

}
