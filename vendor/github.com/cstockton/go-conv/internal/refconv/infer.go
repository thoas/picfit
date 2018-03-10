package refconv

import (
	"fmt"
	"reflect"
)

// Infer will perform conversion by inferring the conversion operation from
// the T of `into`.
func (c Conv) Infer(into, from interface{}) error {
	var value reflect.Value
	switch into := into.(type) {
	case reflect.Value:
		value = into
	default:
		value = reflect.ValueOf(into)
	}

	if !value.IsValid() {
		return fmt.Errorf("%T is not a valid value", into)
	}

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if !value.CanSet() {
		return fmt.Errorf(`cannot infer conversion for unchangeable %v (type %[1]T)`, into)
	}

	v, err := c.infer(value, from)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(v))
	return nil
}

func (c Conv) infer(val reflect.Value, from interface{}) (interface{}, error) {
	switch val.Kind() {
	case reflect.Bool:
		return c.Bool(from)
	case reflect.Float32:
		return c.Float32(from)
	case reflect.Float64:
		return c.Float64(from)
	case reflect.Int:
		return c.Int(from)
	case reflect.Int8:
		return c.Int8(from)
	case reflect.Int16:
		return c.Int16(from)
	case reflect.Int32:
		return c.Int32(from)
	case reflect.String:
		return c.String(from)
	case reflect.Uint:
		return c.Uint(from)
	case reflect.Uint8:
		return c.Uint8(from)
	case reflect.Uint16:
		return c.Uint16(from)
	case reflect.Uint32:
		return c.Uint32(from)
	case reflect.Uint64:
		return c.Uint64(from)
	case reflect.Int64:
		if val.Type() == typeOfDuration {
			return c.Duration(from)
		}
		return c.Int64(from)
	case reflect.Struct:
		if val.Type() == typeOfTime {
			return c.Time(from)
		}
		fallthrough
	default:
		return nil, fmt.Errorf(`cannot infer conversion for %v (type %[1]v)`, val)
	}
}
