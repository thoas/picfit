package binding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorAdd(t *testing.T) {

	Convey("When using Add to add an error to the slice", t, func() {

		Convey("The slice should be the correct size", nil)

		Convey("The slice should have that error", nil)

	})

}

func TestErrorsLen(t *testing.T) {

	Convey("When using Len to get the error count", t, func() {

		Convey("It should return the correct value", nil)

	})

}

func TestErrorsHas(t *testing.T) {

	Convey("When checking for an error classification that exists", t, func() {

		Convey("Has should return true", nil)

	})

	Convey("When checking for an error classification that doesn't exist", t, func() {

		Convey("Has should return false", nil)

	})

}

func TestErrorGetters(t *testing.T) {

	Convey("Given a simple Errors instance", t, func() {

		Convey("Fields should return the fields", nil)

		Convey("Kind should return the classification", nil)

		Convey("Error should return the message", nil)

	})

}
