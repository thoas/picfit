package refutil

import (
	"fmt"
	"reflect"
)

// IsKindComplex returns true if the given Kind is a complex value.
func IsKindComplex(k reflect.Kind) bool {
	return reflect.Complex64 == k || k == reflect.Complex128
}

// IsKindFloat returns true if the given Kind is a float value.
func IsKindFloat(k reflect.Kind) bool {
	return reflect.Float32 == k || k == reflect.Float64
}

// IsKindInt returns true if the given Kind is a int value.
func IsKindInt(k reflect.Kind) bool {
	return reflect.Int <= k && k <= reflect.Int64
}

// IsKindUint returns true if the given Kind is a uint value.
func IsKindUint(k reflect.Kind) bool {
	return reflect.Uint <= k && k <= reflect.Uintptr
}

// IsKindNumeric returns true if the given Kind is a numeric value.
func IsKindNumeric(k reflect.Kind) bool {
	return (reflect.Int <= k && k <= reflect.Uint64) ||
		(reflect.Float32 <= k && k <= reflect.Complex128)
}

// IsKindNillable will return true if the Kind is a chan, func, interface, map,
// pointer, or slice value, false otherwise.
func IsKindNillable(k reflect.Kind) bool {
	return (reflect.Chan <= k && k <= reflect.Slice) || k == reflect.UnsafePointer
}

// IsKindLength will return true if the Kind has a length.
func IsKindLength(k reflect.Kind) bool {
	return reflect.Array == k || reflect.Chan == k || reflect.Map == k ||
		reflect.Slice == k || reflect.String == k
}

// IndirectVal is like Indirect but faster when the caller is using working with
// a reflect.Value.
func IndirectVal(val reflect.Value) reflect.Value {
	var last uintptr
	for {
		if val.Kind() != reflect.Ptr {
			return val
		}

		ptr := val.Pointer()
		if ptr == last {
			return val
		}
		last, val = ptr, val.Elem()
	}
}

// Indirect will perform recursive indirection on the given value. It should
// never panic and will return a value unless indirection is impossible due to
// infinite recursion in cases like `type Element *Element`.
func Indirect(value interface{}) interface{} {
	for {

		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Ptr {
			// Value is not a pointer.
			return value
		}

		res := reflect.Indirect(val)
		if !res.IsValid() || !res.CanInterface() {
			// Invalid value or can't be returned as interface{}.
			return value
		}

		// Test for a circular type.
		if res.Kind() == reflect.Ptr && val.Pointer() == res.Pointer() {
			return value
		}

		// Next round.
		value = res.Interface()
	}
}

// Recover will attempt to execute f, if f return a non-nil error it will be
// returned. If f panics this function will attempt to recover() and return a
// error instead.
func Recover(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch T := r.(type) {
			case error:
				err = T
			default:
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()
	err = f()
	return
}
