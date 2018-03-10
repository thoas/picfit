package testconv

import (
	"reflect"
	"testing"
)

func RunTest(t *testing.T, k reflect.Kind, fn func(interface{}) (interface{}, error)) {
	t.Run(k.String(), func(t *testing.T) {
		if n := assertions.EachOf(k, func(a *Assertion, e Expecter) {
			if err := e.Expect(fn(a.From)); err != nil {
				t.Errorf("(FAIL) %v %v\n%v", a, e, err)
			} else {
				t.Logf("(PASS) %v %v", a, e)
			}
		}); n < 1 {
			t.Fatalf("no test coverage ran for %v conversions", k)
		}
	})
}
