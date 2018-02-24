// Package conv provides fast and intuitive conversions across Go types.
package conv

import (
	"time"

	"github.com/cstockton/go-conv/internal/refconv"
)

var converter = refconv.Conv{}

// Infer will perform conversion by inferring the conversion operation from
// the base type of a pointer to a supported T.
//
// Example:
//
//   var into int64
//   err := conv.Infer(&into, `12`)
//   // into -> 12
//
// See examples for more usages.
func Infer(into, from interface{}) error {
	return converter.Infer(into, from)
}

// Bool will convert the given value to a bool, returns the default value of
// false if a conversion can not be made.
func Bool(from interface{}) (bool, error) {
	return converter.Bool(from)
}

// Duration will convert the given value to a time.Duration, returns the default
// value of 0ns if a conversion can not be made.
func Duration(from interface{}) (time.Duration, error) {
	return converter.Duration(from)
}

// String will convert the given value to a string, returns the default value
// of "" if a conversion can not be made.
func String(from interface{}) (string, error) {
	return converter.String(from)
}

// Time will convert the given value to a time.Time, returns the empty struct
// time.Time{} if a conversion can not be made.
func Time(from interface{}) (time.Time, error) {
	return converter.Time(from)
}

// Float32 will convert the given value to a float32, returns the default value
// of 0.0 if a conversion can not be made.
func Float32(from interface{}) (float32, error) {
	return converter.Float32(from)
}

// Float64 will convert the given value to a float64, returns the default value
// of 0.0 if a conversion can not be made.
func Float64(from interface{}) (float64, error) {
	return converter.Float64(from)
}

// Int will convert the given value to a int, returns the default value of 0 if
// a conversion can not be made.
func Int(from interface{}) (int, error) {
	return converter.Int(from)
}

// Int8 will convert the given value to a int8, returns the default value of 0
// if a conversion can not be made.
func Int8(from interface{}) (int8, error) {
	return converter.Int8(from)
}

// Int16 will convert the given value to a int16, returns the default value of 0
// if a conversion can not be made.
func Int16(from interface{}) (int16, error) {
	return converter.Int16(from)
}

// Int32 will convert the given value to a int32, returns the default value of 0
// if a conversion can not be made.
func Int32(from interface{}) (int32, error) {
	return converter.Int32(from)
}

// Int64 will convert the given value to a int64, returns the default value of 0
// if a conversion can not be made.
func Int64(from interface{}) (int64, error) {
	return converter.Int64(from)
}

// Uint will convert the given value to a uint, returns the default value of 0
// if a conversion can not be made.
func Uint(from interface{}) (uint, error) {
	return converter.Uint(from)
}

// Uint8 will convert the given value to a uint8, returns the default value of 0
// if a conversion can not be made.
func Uint8(from interface{}) (uint8, error) {
	return converter.Uint8(from)
}

// Uint16 will convert the given value to a uint16, returns the default value of
// 0 if a conversion can not be made.
func Uint16(from interface{}) (uint16, error) {
	return converter.Uint16(from)
}

// Uint32 will convert the given value to a uint32, returns the default value of
// 0 if a conversion can not be made.
func Uint32(from interface{}) (uint32, error) {
	return converter.Uint32(from)
}

// Uint64 will convert the given value to a uint64, returns the default value of
// 0 if a conversion can not be made.
func Uint64(from interface{}) (uint64, error) {
	return converter.Uint64(from)
}
