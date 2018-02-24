// Package refconv implements the Converter interface by using the standard
// libraries reflection package.
package refconv

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"

	"github.com/cstockton/go-conv/internal/refutil"
)

type stringConverter interface {
	String() (string, error)
}

// String returns the string representation from the given interface{} value
// and can not currently fail. Although an error is currently provided only for
// API cohesion you should still check it to be future proof.
func (c Conv) String(from interface{}) (string, error) {
	switch T := from.(type) {
	case string:
		return T, nil
	case stringConverter:
		return T.String()
	case []byte:
		return string(T), nil
	case *[]byte:
		// @TODO Maybe validate the bytes are valid runes
		return string(*T), nil
	case *string:
		return *T, nil
	default:
		return fmt.Sprintf("%v", from), nil
	}
}

func (c Conv) convNumToBool(k reflect.Kind, value reflect.Value) (bool, bool) {
	switch {
	case refutil.IsKindInt(k):
		return 0 != value.Int(), true
	case refutil.IsKindUint(k):
		return 0 != value.Uint(), true
	case refutil.IsKindFloat(k):
		T := value.Float()
		if math.IsNaN(T) {
			return false, true
		}
		return 0 != T, true
	case refutil.IsKindComplex(k):
		T := value.Complex()
		if cmplx.IsNaN(T) {
			return false, true
		}
		return 0 != real(T), true
	}
	return false, false
}
