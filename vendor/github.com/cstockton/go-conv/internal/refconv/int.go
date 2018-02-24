package refconv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/cstockton/go-conv/internal/refutil"
)

var (
	mathMaxInt  int64
	mathMinInt  int64
	mathMaxUint uint64
	mathIntSize = strconv.IntSize
)

func initIntSizes(size int) {
	switch size {
	case 64:
		mathMaxInt = math.MaxInt64
		mathMinInt = math.MinInt64
		mathMaxUint = math.MaxUint64
	case 32:
		mathMaxInt = math.MaxInt32
		mathMinInt = math.MinInt32
		mathMaxUint = math.MaxUint32
	}
}

func init() {
	// this is so it can be unit tested.
	initIntSizes(mathIntSize)
}

func (c Conv) convStrToInt64(v string) (int64, error) {
	if parsed, err := strconv.ParseInt(v, 10, 0); err == nil {
		return parsed, nil
	}
	if parsed, err := strconv.ParseFloat(v, 64); err == nil {
		return int64(parsed), nil
	}
	if parsed, err := c.convStrToBool(v); err == nil {
		if parsed {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %#v (type string) to int64", v)
}

type intConverter interface {
	Int64() (int64, error)
}

// Int64 attempts to convert the given value to int64, returns the zero value
// and an error on failure.
func (c Conv) Int64(from interface{}) (int64, error) {
	if T, ok := from.(string); ok {
		return c.convStrToInt64(T)
	} else if T, ok := from.(int64); ok {
		return T, nil
	}
	if c, ok := from.(intConverter); ok {
		return c.Int64()
	}

	value := refutil.IndirectVal(reflect.ValueOf(from))
	kind := value.Kind()
	switch {
	case reflect.String == kind:
		return c.convStrToInt64(value.String())
	case refutil.IsKindInt(kind):
		return value.Int(), nil
	case refutil.IsKindUint(kind):
		val := value.Uint()
		if val > math.MaxInt64 {
			val = math.MaxInt64
		}
		return int64(val), nil
	case refutil.IsKindFloat(kind):
		return int64(value.Float()), nil
	case refutil.IsKindComplex(kind):
		return int64(real(value.Complex())), nil
	case reflect.Bool == kind:
		if value.Bool() {
			return 1, nil
		}
		return 0, nil
	case refutil.IsKindLength(kind):
		return int64(value.Len()), nil
	}
	return 0, newConvErr(from, "int64")
}

// Int attempts to convert the given value to int, returns the zero value and an
// error on failure.
func (c Conv) Int(from interface{}) (int, error) {
	if T, ok := from.(int); ok {
		return T, nil
	}

	to64, err := c.Int64(from)
	if err != nil {
		return 0, newConvErr(from, "int")
	}
	if to64 > mathMaxInt {
		to64 = mathMaxInt // only possible on 32bit arch
	} else if to64 < mathMinInt {
		to64 = mathMinInt // only possible on 32bit arch
	}
	return int(to64), nil
}

// Int8 attempts to convert the given value to int8, returns the zero value and
// an error on failure.
func (c Conv) Int8(from interface{}) (int8, error) {
	if T, ok := from.(int8); ok {
		return T, nil
	}

	to64, err := c.Int64(from)
	if err != nil {
		return 0, newConvErr(from, "int8")
	}
	if to64 > math.MaxInt8 {
		to64 = math.MaxInt8
	} else if to64 < math.MinInt8 {
		to64 = math.MinInt8
	}
	return int8(to64), nil
}

// Int16 attempts to convert the given value to int16, returns the zero value
// and an error on failure.
func (c Conv) Int16(from interface{}) (int16, error) {
	if T, ok := from.(int16); ok {
		return T, nil
	}

	to64, err := c.Int64(from)
	if err != nil {
		return 0, newConvErr(from, "int16")
	}
	if to64 > math.MaxInt16 {
		to64 = math.MaxInt16
	} else if to64 < math.MinInt16 {
		to64 = math.MinInt16
	}
	return int16(to64), nil
}

// Int32 attempts to convert the given value to int32, returns the zero value
// and an error on failure.
func (c Conv) Int32(from interface{}) (int32, error) {
	if T, ok := from.(int32); ok {
		return T, nil
	}

	to64, err := c.Int64(from)
	if err != nil {
		return 0, newConvErr(from, "int32")
	}
	if to64 > math.MaxInt32 {
		to64 = math.MaxInt32
	} else if to64 < math.MinInt32 {
		to64 = math.MinInt32
	}
	return int32(to64), nil
}
