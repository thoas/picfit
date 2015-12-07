package binding

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidate(t *testing.T) {
	Convey("Given a struct populated properly and as expected", t, func() {

		Convey("No errors should be produced", FailureContinues, func() {
			req, err := http.NewRequest("POST", "http://www.example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			model := NewCompleteModel()
			errs := Validate(req, &model)

			expectedErrs := make(map[string]bool)
			for _, v := range model.FieldMap() {
				f, ok := v.(Field)
				if !ok {
					t.Fatal("unexpected value in FieldMap")
				}

				expectedErrs[f.Form] = false
			}

			for _, bindErr := range errs {
				for _, name := range bindErr.FieldNames {
					if bindErr.Classification == RequiredError {
						expectedErrs[name] = true
					}
				}
			}

			for k, v := range expectedErrs {
				Convey(fmt.Sprintf("A Required error for %s should not be produced", k), func() {
					if v {
						Println(k, "has an unexpected Required error")
					}
					So(v, ShouldBeFalse)
				})
			}
		})

	})

	Convey("Given a populated struct missing one required field", t, func() {

		Convey("A Required error should be produced", nil)

	})

	Convey("Given a populated struct missing multiple required fields", t, func() {

		Convey("As many Required errors should be produced", FailureContinues, func() {
			req, err := http.NewRequest("POST", "http://www.example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			model := new(AllTypes)
			errs := Validate(req, model)

			expectedErrs := make(map[string]bool)
			for _, v := range model.FieldMap() {
				f, ok := v.(Field)
				if !ok {
					t.Fatal("unexpected value in FieldMap")
				}

				expectedErrs[f.Form] = false
			}

			for _, bindErr := range errs {
				for _, name := range bindErr.FieldNames {
					if bindErr.Classification == RequiredError {
						expectedErrs[name] = true
					}
				}
			}

			for k, v := range expectedErrs {
				Convey(fmt.Sprintf("A Required error for %s should be produced", k), func() {
					if !v {
						Println(k, "is missing the expected Required error")
					}
					So(v, ShouldBeTrue)
				})
			}
		})
	})

	Convey("Given a struct that is a Validator", t, func() {

		Convey("The user's Validate method should be invoked and its errors appended", nil)

	})

	Convey("Each case in the type switch should be tested", t, func() {

		Convey("uint8", nil)
		Convey("*uint8", nil)
		Convey("[]uint8", nil)

		Convey("uint16", nil)
		Convey("*uint16", nil)
		Convey("[]uint16", nil)

		Convey("uint32", nil)
		Convey("*uint32", nil)
		Convey("[]uint32", nil)

		Convey("uint64", nil)
		Convey("*uint64", nil)
		Convey("[]uint64", nil)

		Convey("int8", nil)
		Convey("*int8", nil)
		Convey("[]int8", nil)

		Convey("int16", nil)
		Convey("*int16", nil)
		Convey("[]int16", nil)

		Convey("int32", nil)
		Convey("*int32", nil)
		Convey("[]int32", nil)

		Convey("int64", nil)
		Convey("*int64", nil)
		Convey("[]int64", nil)

		Convey("float32", nil)
		Convey("*float32", nil)
		Convey("[]float32", nil)

		Convey("float64", nil)
		Convey("*float64", nil)
		Convey("[]float64", nil)

		Convey("uint", nil)
		Convey("*uint", nil)
		Convey("[]uint", nil)

		Convey("int", nil)
		Convey("*int", nil)
		Convey("[]int", nil)

		Convey("bool", nil)
		Convey("*bool", nil)
		Convey("[]bool", nil)

		Convey("string", nil)
		Convey("*string", nil)
		Convey("[]string", nil)

		Convey("time.Time", nil)
		Convey("*time.Time", nil)
		Convey("[]time.Time", nil)

	})
}
