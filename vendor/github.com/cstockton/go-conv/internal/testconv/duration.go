package testconv

import (
	"math"
	"math/cmplx"
	"testing"
	"time"
)

func RunDurationTests(t *testing.T, fn func(interface{}) (time.Duration, error)) {
	RunTest(t, DurationKind, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

type testDurationConverter time.Duration

func (d testDurationConverter) Duration() (time.Duration, error) {
	return time.Duration(d) + time.Minute, nil
}

func init() {
	var (
		dZero      time.Duration
		d42ns              = time.Nanosecond * 42
		d2m                = time.Minute * 2
		d34s               = time.Second * 34
		d567ms             = time.Millisecond * 567
		d234567            = d2m + d34s + d567ms
		d154567f32 float32 = 154.567
		d154567f64         = 154.567
		dMAX               = time.Duration(math.MaxInt64)
	)

	// strings
	assert("2m34.567s", d234567)
	assert("-2m34.567s", -d234567)
	assert("154.567", d234567)
	assert("-154.567", -d234567)
	assert("42", d42ns)
	assert(testStringConverter("42"), d42ns)

	// durations
	assert(d234567, d234567)
	assert(dZero, dZero)
	assert(new(time.Duration), dZero)

	// underlying
	type ulyDuration time.Duration
	assert(ulyDuration(time.Second), time.Second)
	assert(ulyDuration(-time.Second), -time.Second)

	// implements converter
	assert(testDurationConverter(time.Second), time.Second+time.Minute)
	assert(testDurationConverter(-time.Second), time.Minute-time.Second)

	// numerics
	assert(int(42), d42ns)
	assert(int8(42), d42ns)
	assert(int16(42), d42ns)
	assert(int32(42), d42ns)
	assert(int64(42), d42ns)
	assert(uint(42), d42ns)
	assert(uint8(42), d42ns)
	assert(uint16(42), d42ns)
	assert(uint32(42), d42ns)
	assert(uint64(42), d42ns)

	// floats
	assert(d154567f32, DurationExp{d234567, time.Millisecond})
	assert(d154567f64, d234567)
	assert(math.NaN(), dZero)
	assert(math.Inf(1), dZero)
	assert(math.Inf(-1), dZero)

	// complex
	assert(complex(d154567f32, 0), DurationExp{d234567, time.Millisecond})
	assert(complex(d154567f64, 0), d234567)
	assert(cmplx.NaN(), dZero)
	assert(cmplx.Inf(), dZero)

	// overflow
	assert(uint64(math.MaxUint64), dMAX)

	// errors
	assert(nil, experr(dZero, `cannot convert <nil> (type <nil>) to time.Duration`))
	assert("foo", experr(dZero, `cannot parse "foo" (type string) as time.Duration`))
	assert("tooLong", experr(
		dZero, `cannot parse "tooLong" (type string) as time.Duration`))
	assert(struct{}{}, experr(
		dZero, `cannot convert struct {}{} (type struct {}) to `))
	assert([]string{"1s"}, experr(
		dZero, `cannot convert []string{"1s"} (type []string) to `))
	assert([]string{}, experr(
		dZero, `cannot convert []string{} (type []string) to `))
}
