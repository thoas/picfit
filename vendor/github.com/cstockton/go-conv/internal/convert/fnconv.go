package convert

import (
	"sync/atomic"
	"time"
)

func WithError(err error, fn FnConv) FnConv {
	return WithErrorAfter(err, 0, fn)
}

func WithErrorAfter(err error, after int, fn FnConv) FnConv {
	var n int64
	return func(into, from interface{}) error {
		if atomic.AddInt64(&n, 1) > int64(after) {
			return err
		}
		n++
		return nil
	}
}

type FnConv func(into, from interface{}) error

func (fn FnConv) Bool(from interface{}) (out bool, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Duration(from interface{}) (out time.Duration, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Float32(from interface{}) (out float32, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Float64(from interface{}) (out float64, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Infer(into, from interface{}) (err error) {
	err = fn(into, from)
	return
}

func (fn FnConv) Int(from interface{}) (out int, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Int8(from interface{}) (out int8, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Int16(from interface{}) (out int16, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Int32(from interface{}) (out int32, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Int64(from interface{}) (out int64, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Map(into, from interface{}) (err error) {
	err = fn(into, from)
	return
}

func (fn FnConv) Slice(into, from interface{}) (err error) {
	err = fn(into, from)
	return
}

func (fn FnConv) String(from interface{}) (out string, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Struct(into, from interface{}) (err error) {
	err = fn(into, from)
	return
}

func (fn FnConv) Time(from interface{}) (out time.Time, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Uint(from interface{}) (out uint, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Uint8(from interface{}) (out uint8, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Uint16(from interface{}) (out uint16, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Uint32(from interface{}) (out uint32, err error) {
	err = fn(&out, from)
	return
}

func (fn FnConv) Uint64(from interface{}) (out uint64, err error) {
	err = fn(&out, from)
	return
}
