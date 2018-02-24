// Package convert contains common conversion interfaces.
package convert

import "time"

// Converter supports conversion across Go types.
type Converter interface {

	// Bool returns the bool representation from the given interface value.
	// Returns the default value of false and an error on failure.
	Bool(from interface{}) (to bool, err error)

	// Duration returns the time.Duration representation from the given
	// interface{} value. Returns the default value of 0 and an error on failure.
	Duration(from interface{}) (to time.Duration, err error)

	// Float32 returns the float32 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Float32(from interface{}) (to float32, err error)

	// Float64 returns the float64 representation from the given interface
	// value. Returns the default value of 0 and an error on failure.
	Float64(from interface{}) (to float64, err error)

	// Infer will perform conversion by inferring the conversion operation from
	// a pointer to a supported T of the `into` param.
	Infer(into, from interface{}) error

	// Int returns the int representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Int(from interface{}) (to int, err error)

	// Int8 returns the int8 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Int8(from interface{}) (to int8, err error)

	// Int16 returns the int16 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Int16(from interface{}) (to int16, err error)

	// Int32 returns the int32 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Int32(from interface{}) (to int32, err error)

	// Int64 returns the int64 representation from the given interface
	// value. Returns the default value of 0 and an error on failure.
	Int64(from interface{}) (to int64, err error)

	// String returns the string representation from the given interface
	// value and can not fail. An error is provided only for API cohesion.
	String(from interface{}) (to string, err error)

	// Time returns the time.Time{} representation from the given interface
	// value. Returns an empty time.Time struct and an error on failure.
	Time(from interface{}) (to time.Time, err error)

	// Uint returns the uint representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Uint(from interface{}) (to uint, err error)

	// Uint8 returns the uint8 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Uint8(from interface{}) (to uint8, err error)

	// Uint16 returns the uint16 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Uint16(from interface{}) (to uint16, err error)

	// Uint32 returns the uint32 representation from the given empty interface
	// value. Returns the default value of 0 and an error on failure.
	Uint32(from interface{}) (to uint32, err error)

	// Uint64 returns the uint64 representation from the given interface
	// value. Returns the default value of 0 and an error on failure.
	Uint64(from interface{}) (to uint64, err error)
}
