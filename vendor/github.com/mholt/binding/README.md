[<img src="http://mholt.github.io/binding/resources/images/binding-sm.png" height="250" alt="binding is reflectionless data binding for Go"></a>](http://mholt.github.io/binding)

[![GoDoc](https://godoc.org/github.com/mholt/binding?status.svg)](https://godoc.org/github.com/mholt/binding)

binding
=======

Reflectionless data binding for Go's net/http



Features
---------

- HTTP request data binding
- Data validation (custom and built-in)
- Error handling



Benefits
---------

- Moves data binding, validation, and error handling out of your application's handler
- Reads Content-Type to deserialize form, multipart form, and JSON data from requests
- No middleware: just a function call
- Usable in any setting where `net/http` is present (Negroni, gocraft/web, std lib, etc.)
- No reflection



Usage example
--------------

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/mholt/binding"
)

// First define a type to hold the data
// (If the data comes from JSON, see: http://mholt.github.io/json-to-go)
type ContactForm struct {
	User struct {
		ID int
	}
	Email   string
	Message string
}

// Then provide a field mapping (pointer receiver is vital)
func (cf *ContactForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.User.ID: "user_id",
		&cf.Email:   "email",
		&cf.Message: binding.Field{
			Form:     "message",
			Required: true,
		},
	}
}

// Now your handlers can stay clean and simple
func handler(resp http.ResponseWriter, req *http.Request) {
	contactForm := new(ContactForm)
	if errs := binding.Bind(req, contactForm); errs != nil {
		http.Error(resp, errs.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(resp, "From:    %d\n", contactForm.User.ID)
	fmt.Fprintf(resp, "Message: %s\n", contactForm.Message)
}

func main() {
	http.HandleFunc("/contact", handler)
	http.ListenAndServe(":3000", nil)
}
```

Multipart/form-data usage example
---------------------------------

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/mholt/binding"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

// We expect a multipart/form-data upload with
// a file field named 'data'
type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

// Handlers are still clean and simple
func handler(resp http.ResponseWriter, req *http.Request) {
	multipartForm := new(MultipartForm)
	if errs := binding.Bind(req, multipartForm); errs != nil {
		http.Error(resp, errs.Error(), http.StatusBadRequest)
		return
	}

	// To access the file data you need to Open the file
	// handler and read the bytes out.
	var fh io.ReadCloser
	var err error
	if fh, err = multipartForm.Data.Open(); err != nil {
		http.Error(resp,
			fmt.Sprint("Error opening Mime::Data %v", err),
			http.StatusInternalServerError)
		return
	}
	defer fh.Close()
	dataBytes := bytes.Buffer{}
	var size int64
	if size, err = dataBytes.ReadFrom(fh); err != nil {
		http.Error(resp,
			fmt.Sprint("Error reading Mime::Data %v", err),
			http.StatusInternalServerError)
		return
	}
	// Now you have the attachment in databytes.
	// Maximum size is default is 10MB.
	log.Printf("Read %v bytes with filename %s",
		size, multipartForm.Data.Filename)
}

func main() {
	http.HandleFunc("/upload", handler)
	http.ListenAndServe(":3000", nil)
}

```

You can test from CLI using the excellent httpie client

   `http -f POST localhost:3000/upload data@myupload`

Custom data validation
-----------------------

You may optionally have your type implement the `binding.Validator` interface to perform your own data validation. The `.Validate()` method is called after the struct is populated.

```go
func (cf ContactForm) Validate(req *http.Request) error {
	if cf.Message == "Go needs generics" {
		return binding.Errors{
			binding.NewError([]string{"message"}, "ComplaintError", "Go has generics. They're called interfaces.")
		}
	}
	return nil
}
```



Binding custom types
---------------------

For types you've defined, you can bind form data to it by implementing the `Binder` interface. Here's a contrived example:


```go
type MyBinder map[string]string

func (t MyBinder) Bind(fieldName string, strVals []string) error {
	t["formData"] = strVals[0]
	return nil
}
```

If you can't add a method to the type, you can still specify a `Binder` func in the field spec. Here's a contrived example that binds an integer (not necessary, but you get the idea):

```go
func (t *MyType) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		"a-key": binding.Field{
			Form: "number",
			Binder: func(fieldName string, formVals []string) error {
				val, err := strconv.Atoi(formVals[0])
				if err != nil {
					return binding.Errors{binding.NewError([]string{fieldName}, binding.DeserializationError, err.Error())}
				}
				t.SomeNumber = val
				return nil
			},
		},
	}
}
```

The `Errors` type has a convenience method, `Add`, which you can use to append to the slice if you prefer.

Supported types (forms)
------------------------

The following types are supported in form deserialization by default. (JSON requests are delegated to `encoding/json`.)

- uint, \*uint, []uint, uint8, \*uint8, []uint8, uint16, \*uint16, []uint16, uint32, \*uint32, []uint32, uint64, \*uint64, []uint64
- int, \*int, []int, int8, \*int8, []int8, int16, \*int16, []int16, int32, \*int32, []int32, int64, \*int64, []int64
- float32, \*float32, []float32, float64, \*float64, []float64
- bool, \*bool, []bool
- string, \*string, []string
- time.Time, \*time.Time, []time.Time
- \*multipart.FileHeader, []\*multipart.FileHeader
