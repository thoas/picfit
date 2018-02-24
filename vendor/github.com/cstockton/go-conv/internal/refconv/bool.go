package refconv

import (
	"fmt"
	"reflect"
	"time"

	"github.com/cstockton/go-conv/internal/refutil"
)

type boolConverter interface {
	Bool() (bool, error)
}

// Bool attempts to convert the given value to bool, returns the zero value
// and an error on failure.
func (c Conv) Bool(from interface{}) (bool, error) {
	if T, ok := from.(string); ok {
		return c.convStrToBool(T)
	} else if T, ok := from.(bool); ok {
		return T, nil
	} else if c, ok := from.(boolConverter); ok {
		return c.Bool()
	}

	value := refutil.IndirectVal(reflect.ValueOf(from))
	kind := value.Kind()
	switch {
	case reflect.String == kind:
		return c.convStrToBool(value.String())
	case refutil.IsKindNumeric(kind):
		if parsed, ok := c.convNumToBool(kind, value); ok {
			return parsed, nil
		}
	case reflect.Bool == kind:
		return value.Bool(), nil
	case refutil.IsKindLength(kind):
		return value.Len() > 0, nil
	case reflect.Struct == kind && value.CanInterface():
		v := value.Interface()
		if t, ok := v.(time.Time); ok {
			return emptyTime != t, nil
		}
	}
	return false, newConvErr(from, "bool")
}

func (c Conv) convStrToBool(v string) (bool, error) {
	// @TODO Need to find a clean way to expose the truth list to be modified by
	// API to allow INTL.
	if 1 > len(v) || len(v) > 5 {
		return false, fmt.Errorf("cannot parse string with len %d as bool", len(v))
	}

	// @TODO lut
	switch v {
	case "1", "t", "T", "true", "True", "TRUE", "y", "Y", "yes", "Yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "False", "FALSE", "n", "N", "no", "No", "NO":
		return false, nil
	}
	return false, fmt.Errorf("cannot parse %#v (type string) as bool", v)
}
