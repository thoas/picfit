package binding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestForm(t *testing.T) {

	Convey("Given a proper request with a query string", t, func() {

		Convey("It should deserialize", nil)

		Convey("There should be no errors", nil)

	})

	Convey("Given a proper request with a form body", t, func() {

		Convey("It should deserialize", nil)

		Convey("There should be no errors", nil)

	})

	Convey("Each case in the type switch should be tested", t, func() {

		Convey("uint8", nil)

		Convey("uint16", nil)

		Convey("uint32", nil)

		Convey("uint64", nil)

		Convey("int8", nil)

		Convey("int16", nil)

		Convey("int32", nil)

		Convey("int64", nil)

		Convey("float32", nil)

		Convey("[]float32", nil)

		Convey("float64", nil)

		Convey("[]float64", nil)

		Convey("uint", nil)

		Convey("[]uint", nil)

		Convey("int", nil)

		Convey("[]int", nil)

		Convey("bool", nil)

		Convey("[]bool", nil)

		Convey("string", nil)

		Convey("[]string", nil)

		Convey("time.Time", nil)

		Convey("[]time.Time", nil)

	})

}
