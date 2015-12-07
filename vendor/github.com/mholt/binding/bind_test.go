package binding

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Model struct {
	Foo   string     `json:"foo"`
	Bar   *string    `json:"bar"`
	Baz   []int      `json:"baz"`
	Child ChildModel `json:"child"`
}

func (m *Model) FieldMap() FieldMap {
	return FieldMap{
		&m.Foo: Field{
			Form:     "foo",
			Required: true,
		},
		&m.Bar: Field{
			Form:     "bar",
			Required: true,
		},
		&m.Baz:   "baz",
		&m.Child: "child",
	}
}

type ChildModel struct {
	Wibble string `json:"wibble"`
}

func TestBind(t *testing.T) {
	Convey("A request", t, func() {

		Convey("Without a Content-Type", func() {

			Convey("But with a query string", func() {

				Convey("Should invoke the Form deserializer", nil)

			})

			Convey("And without a query string", func() {

				Convey("Should yield an error", nil)

			})

		})

		Convey("With a form-urlencoded Content-Type", func() {
			data := url.Values{}
			data.Add("foo", "foo-value")
			data.Add("child.wibble", "wobble")
			data.Add("baz", "1")
			data.Add("baz", "2")
			data.Add("baz", "3")
			req, err := http.NewRequest("POST", "http://www.example.com", strings.NewReader(data.Encode()))
			So(err, ShouldBeNil)
			req.Header.Add("Content-type", "application/x-www-form-urlencoded")

			Convey("Should invoke the Form deserializer", func() {
				model := new(Model)
				invoked := false
				formBinder = func(req *http.Request, v FieldMapper) Errors {
					invoked = true
					return defaultFormBinder(req, v)
				}
				Bind(req, model)
				So(invoked, ShouldBeTrue)
				formBinder = defaultFormBinder
			})
		})

		Convey("With a multipart/form-data Content-Type", func() {
			body := new(bytes.Buffer)
			w := multipart.NewWriter(body)
			_ = w.WriteField("foo", "foo-value")
			_ = w.WriteField("child.wibble", "wobble")
			_ = w.WriteField("baz", "1")
			_ = w.WriteField("baz", "2")
			_ = w.WriteField("baz", "3")
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "http://www.example.com", body)
			So(err, ShouldBeNil)
			req.Header.Add("Content-Type", "multipart/form-data")

			Convey("Should invoke the MultipartForm deserializer", func() {
				model := new(Model)
				invoked := false
				multipartFormBinder = func(req *http.Request, v FieldMapper) Errors {
					invoked = true
					return defaultMultipartFormBinder(req, v)
				}
				Bind(req, model)
				So(invoked, ShouldBeTrue)
				multipartFormBinder = defaultMultipartFormBinder
			})
		})

		Convey("With a json Content-Type", func() {
			data := `{ "foo": "foo-value", "child": { "wibble": "wobble" }, "baz": [1,2,3]}`
			req, err := http.NewRequest("POST", "http://www.example.com", strings.NewReader(data))
			So(err, ShouldBeNil)
			req.Header.Add("Content-type", "application/json; charset=utf-8")

			Convey("Should invoke Json deserializer", func() {
				model := new(Model)
				invoked := false
				jsonBinder = func(req *http.Request, v FieldMapper) Errors {
					invoked = true
					return defaultJsonBinder(req, v)
				}
				Bind(req, model)
				So(invoked, ShouldBeTrue)
				jsonBinder = defaultJsonBinder
			})
		})

		Convey("With an unsupported Content-Type", func() {

			Convey("Should yield an error", nil)

		})

	})
}

func TestBindForm(t *testing.T) {
	Convey("Given a struct reference and complete form data", t, func() {
		expected := NewCompleteModel()
		formData := expected.FormValues()

		Convey("Given that all of the struct's fields are required", func() {
			actual := AllTypes{}
			Convey("When bindForm is called", func() {
				req, err := http.NewRequest("POST", "http://www.example.com", nil)
				So(err, ShouldBeNil)
				var errs Errors
				errs = bindForm(req, &actual, formData, nil, errs)
				Convey("Then all of the struct's fields should be populated", func() {
					Convey("Then the Uint8 field should have the expected value", func() {
						So(actual.Uint8, ShouldEqual, expected.Uint8)
					})
					Convey("Then the PointerToUint8 field should have the expected value", func() {
						So(*actual.PointerToUint8, ShouldEqual, *expected.PointerToUint8)
					})
					Convey("Then the Uint8Slice field should have the expected values", func() {
						So(len(actual.Uint8Slice), ShouldEqual, len(expected.Uint8Slice))
						for i := range actual.Uint8Slice {
							So(actual.Uint8Slice[i], ShouldEqual, expected.Uint8Slice[i])
						}
					})
					Convey("Then the Uint16 field should have the expected value", func() {
						So(actual.Uint16, ShouldEqual, expected.Uint16)
					})
					Convey("Then the PointerToUint16 field should have the expected value", func() {
						So(*actual.PointerToUint16, ShouldEqual, *expected.PointerToUint16)
					})
					Convey("Then the Uint16Slice field should have the expected values", func() {
						So(len(actual.Uint16Slice), ShouldEqual, len(expected.Uint16Slice))
						for i := range actual.Uint16Slice {
							So(actual.Uint16Slice[i], ShouldEqual, expected.Uint16Slice[i])
						}
					})
					Convey("Then the Uint32 field should have the expected value", func() {
						So(actual.Uint32, ShouldEqual, expected.Uint32)
					})
					Convey("Then the PointerToUint32 field should have the expected value", func() {
						So(*actual.PointerToUint32, ShouldEqual, *expected.PointerToUint32)
					})
					Convey("Then the Uint32Slice field should have the expected values", func() {
						So(len(actual.Uint32Slice), ShouldEqual, len(expected.Uint32Slice))
						for i := range actual.Uint32Slice {
							So(actual.Uint32Slice[i], ShouldEqual, expected.Uint32Slice[i])
						}
					})
					Convey("Then the Uint64 field should have the expected value", func() {
						So(actual.Uint64, ShouldEqual, expected.Uint64)
					})
					Convey("Then the PointerToUint64 field should have the expected value", func() {
						So(*actual.PointerToUint64, ShouldEqual, *expected.PointerToUint64)
					})
					Convey("Then the Uint64Slice field should have the expected values", func() {
						So(len(actual.Uint64Slice), ShouldEqual, len(expected.Uint64Slice))
						for i := range actual.Uint64Slice {
							So(actual.Uint64Slice[i], ShouldEqual, expected.Uint64Slice[i])
						}
					})
					Convey("Then the Int8 field should have the expected value", func() {
						So(actual.Int8, ShouldEqual, expected.Int8)
					})
					Convey("Then the PointerToInt8 field should have the expected value", func() {
						So(*actual.PointerToInt8, ShouldEqual, *expected.PointerToInt8)
					})
					Convey("Then the Int8Slice field should have the expected values", func() {
						So(len(actual.Int8Slice), ShouldEqual, len(expected.Int8Slice))
						for i := range actual.Int8Slice {
							So(actual.Int8Slice[i], ShouldEqual, expected.Int8Slice[i])
						}
					})
					Convey("Then the Int16 field should have the expected value", func() {
						So(actual.Int16, ShouldEqual, expected.Int16)
					})
					Convey("Then the PointerToInt16 field should have the expected value", func() {
						So(*actual.PointerToInt16, ShouldEqual, *expected.PointerToInt16)
					})
					Convey("Then the Int16Slice field should have the expected values", func() {
						So(len(actual.Int16Slice), ShouldEqual, len(expected.Int16Slice))
						for i := range actual.Int16Slice {
							So(actual.Int16Slice[i], ShouldEqual, expected.Int16Slice[i])
						}
					})
					Convey("Then the Int32 field should have the expected value", func() {
						So(actual.Int32, ShouldEqual, expected.Int32)
					})
					Convey("Then the PointerToInt32 field should have the expected value", func() {
						So(*actual.PointerToInt32, ShouldEqual, *expected.PointerToInt32)
					})
					Convey("Then the Int32Slice field should have the expected values", func() {
						So(len(actual.Int32Slice), ShouldEqual, len(expected.Int32Slice))
						for i := range actual.Int32Slice {
							So(actual.Int32Slice[i], ShouldEqual, expected.Int32Slice[i])
						}
					})
					Convey("Then the Int64 field should have the expected value", func() {
						So(actual.Int64, ShouldEqual, expected.Int64)
					})
					Convey("Then the PointerToInt64 field should have the expected value", func() {
						So(*actual.PointerToInt64, ShouldEqual, *expected.PointerToInt64)
					})
					Convey("Then the Int64Slice field should have the expected values", func() {
						So(len(actual.Int64Slice), ShouldEqual, len(expected.Int64Slice))
						for i := range actual.Int64Slice {
							So(actual.Int64Slice[i], ShouldEqual, expected.Int64Slice[i])
						}
					})
					Convey("Then the Float32 field should have the expected value", func() {
						So(actual.Float32, ShouldEqual, expected.Float32)
					})
					Convey("Then the PointerToFloat32 field should have the expected value", func() {
						So(*actual.PointerToFloat32, ShouldEqual, *expected.PointerToFloat32)
					})
					Convey("Then the Float32Slice field should have the expected values", func() {
						So(len(actual.Float32Slice), ShouldEqual, len(expected.Float32Slice))
						for i := range actual.Float32Slice {
							So(actual.Float32Slice[i], ShouldEqual, expected.Float32Slice[i])
						}
					})
					Convey("Then the Float64 field should have the expected value", func() {
						So(actual.Float64, ShouldEqual, expected.Float64)
					})
					Convey("Then the PointerToFloat64 field should have the expected value", func() {
						So(*actual.PointerToFloat64, ShouldEqual, *expected.PointerToFloat64)
					})
					Convey("Then the Float64Slice field should have the expected values", func() {
						So(len(actual.Float64Slice), ShouldEqual, len(expected.Float64Slice))
						for i := range actual.Float64Slice {
							So(actual.Float64Slice[i], ShouldEqual, expected.Float64Slice[i])
						}
					})
					Convey("Then the Uint field should have the expected value", func() {
						So(actual.Uint, ShouldEqual, expected.Uint)
					})
					Convey("Then the PointerToUint field should have the expected value", func() {
						So(*actual.PointerToUint, ShouldEqual, *expected.PointerToUint)
					})
					Convey("Then the UintSlice field should have the expected values", func() {
						So(len(actual.UintSlice), ShouldEqual, len(expected.UintSlice))
						for i := range actual.UintSlice {
							So(actual.UintSlice[i], ShouldEqual, expected.UintSlice[i])
						}
					})
					Convey("Then the Int field should have the expected value", func() {
						So(actual.Int, ShouldEqual, expected.Int)
					})
					Convey("Then the PointerToInt field should have the expected value", func() {
						So(*actual.PointerToInt, ShouldEqual, *expected.PointerToInt)
					})
					Convey("Then the IntSlice field should have the expected values", func() {
						So(len(actual.IntSlice), ShouldEqual, len(expected.IntSlice))
						for i := range actual.IntSlice {
							So(actual.IntSlice[i], ShouldEqual, expected.IntSlice[i])
						}
					})
					Convey("Then the Bool field should have the expected value", func() {
						So(actual.Bool, ShouldEqual, expected.Bool)
					})
					Convey("Then the PointerToBool field should have the expected value", func() {
						So(*actual.PointerToBool, ShouldEqual, *expected.PointerToBool)
					})
					Convey("Then the BoolSlice field should have the expected values", func() {
						So(len(actual.BoolSlice), ShouldEqual, len(expected.BoolSlice))
						for i := range actual.BoolSlice {
							So(actual.BoolSlice[i], ShouldEqual, expected.BoolSlice[i])
						}
					})
					Convey("Then the String field should have the expected value", func() {
						So(actual.String, ShouldEqual, expected.String)
					})
					Convey("Then the PointerToString field should have the expected value", func() {
						So(*actual.PointerToString, ShouldEqual, *expected.PointerToString)
					})
					Convey("Then the StringSlice field should have the expected values", func() {
						So(len(actual.StringSlice), ShouldEqual, len(expected.StringSlice))
						for i := range actual.StringSlice {
							So(actual.StringSlice[i], ShouldEqual, expected.StringSlice[i])
						}
					})
					Convey("Then the Time field should have the expected value", func() {
						So(actual.Time.Equal(expected.Time), ShouldBeTrue)
					})
					Convey("Then the PointerToTime field should have the expected value", func() {
						So((*actual.PointerToTime).Equal(*expected.PointerToTime), ShouldBeTrue)
					})
					Convey("Then the TimeSlice field should have the expected values", func() {
						So(len(actual.TimeSlice), ShouldEqual, len(expected.TimeSlice))
						for i := range actual.TimeSlice {
							So(actual.TimeSlice[i].Equal(expected.TimeSlice[i]), ShouldBeTrue)
						}
					})
				})

				Convey("Then no errors should be produced", FailureContinues, func() {
					So(errs.Len(), ShouldEqual, 0)
					if errs.Len() > 0 {
						for _, e := range errs {
							Println(fmt.Sprintf("%v. %s", e.FieldNames, e.Message))
						}
					}
				})
			})
		})
	})
}
