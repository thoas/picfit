// Package refconv implements the Converter interface by using the standard
// libraries reflection package.
package refconv

import (
	"fmt"
)

// Conv implements the Converter interface by using the reflection package. It
// will never panic, does not require initialization or share state so is safe
// for use by multiple Goroutines.
type Conv struct{}

func newConvErr(from interface{}, to string) error {
	return fmt.Errorf("cannot convert %#v (type %[1]T) to %v", from, to)
}
