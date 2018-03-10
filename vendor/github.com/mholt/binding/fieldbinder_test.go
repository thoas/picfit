package binding_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"

	"github.com/mholt/binding"
)

type MyType struct {
	SomeNumber int
}

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

func Example_fieldBinder() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b := new(MyType)
		if err := binding.Bind(req, b); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, "%d", b.SomeNumber)
	}))
	defer ts.Close()

	resp, err := http.DefaultClient.PostForm(ts.URL, url.Values{"number": []string{"1008"}})
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	// Output: 1008
}
