package convert

import (
	"testing"
)

func TestSmoke(t *testing.T) {
	var c interface{} = FnConv(func(into, from interface{}) error { return nil })
	_ = c.(Converter)
}
