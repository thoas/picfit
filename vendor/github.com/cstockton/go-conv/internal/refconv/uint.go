package refconv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/cstockton/go-conv/internal/refutil"
)

func (c Conv) convStrToUint64(v string) (uint64, error) {
	if parsed, err := strconv.ParseUint(v, 10, 0); err == nil {
		return parsed, nil
	}
	if parsed, err := strconv.ParseFloat(v, 64); err == nil {
		return uint64(math.Max(0, parsed)), nil
	}
	if parsed, err := c.convStrToBool(v); err == nil {
		if parsed {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %#v (type string) to uint64", v)
}

type uintConverter interface {
	Uint64() (uint64, error)
}

// Uint64 attempts to convert the given value to uint64, returns the zero value
// and an error on failure.
func (c Conv) Uint64(from interface{}) (uint64, error) {
	if T, ok := from.(string); ok {
		return c.convStrToUint64(T)
	} else if T, ok := from.(uint64); ok {
		return T, nil
	}
	if c, ok := from.(uintConverter); ok {
		return c.Uint64()
	}

	value := refutil.IndirectVal(reflect.ValueOf(from))
	kind := value.Kind()
	switch {
	case reflect.String == kind:
		return c.convStrToUint64(value.String())
	case refutil.IsKindUint(kind):
		return value.Uint(), nil
	case refutil.IsKindInt(kind):
		val := value.Int()
		if val < 0 {
			val = 0
		}
		return uint64(val), nil
	case refutil.IsKindFloat(kind):
		return uint64(math.Max(0, value.Float())), nil
	case refutil.IsKindComplex(kind):
		return uint64(math.Max(0, real(value.Complex()))), nil
	case reflect.Bool == kind:
		if value.Bool() {
			return 1, nil
		}
		return 0, nil
	case refutil.IsKindLength(kind):
		return uint64(value.Len()), nil
	}

	return 0, newConvErr(from, "uint64")
}

// Uint attempts to convert the given value to uint, returns the zero value and
// an error on failure.
func (c Conv) Uint(from interface{}) (uint, error) {
	if T, ok := from.(uint); ok {
		return T, nil
	}

	to64, err := c.Uint64(from)
	if err != nil {
		return 0, newConvErr(from, "uint")
	}
	if to64 > mathMaxUint {
		to64 = mathMaxUint // only possible on 32bit arch
	}
	return uint(to64), nil
}

// Uint8 attempts to convert the given value to uint8, returns the zero value
// and an error on failure.
func (c Conv) Uint8(from interface{}) (uint8, error) {
	if T, ok := from.(uint8); ok {
		return T, nil
	}

	to64, err := c.Uint64(from)
	if err != nil {
		return 0, newConvErr(from, "uint8")
	}
	if to64 > math.MaxUint8 {
		to64 = math.MaxUint8
	}
	return uint8(to64), nil
}

// Uint16 attempts to convert the given value to uint16, returns the zero value
// and an error on failure.
func (c Conv) Uint16(from interface{}) (uint16, error) {
	if T, ok := from.(uint16); ok {
		return T, nil
	}

	to64, err := c.Uint64(from)
	if err != nil {
		return 0, newConvErr(from, "uint16")
	}
	if to64 > math.MaxUint16 {
		to64 = math.MaxUint16
	}
	return uint16(to64), nil
}

// Uint32 attempts to convert the given value to uint32, returns the zero value
// and an error on failure.
func (c Conv) Uint32(from interface{}) (uint32, error) {
	if T, ok := from.(uint32); ok {
		return T, nil
	}

	to64, err := c.Uint64(from)
	if err != nil {
		return 0, newConvErr(from, "uint32")
	}
	if to64 > math.MaxUint32 {
		to64 = math.MaxUint32
	}
	return uint32(to64), nil
}
