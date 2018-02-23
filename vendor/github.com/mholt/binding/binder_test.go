package binding_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/mholt/binding"
)

type MyBinder map[string]string

func (t MyBinder) Bind(fieldName string, strVals []string) error {
	t["formData"] = strVals[0]
	return nil
}

type MyBinderContainer struct {
	Important MyBinder
}

func (c *MyBinderContainer) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.Important: "important",
	}
}

func ExampleBinder() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		v := new(MyBinderContainer)
		v.Important = make(MyBinder)
		if err := binding.Bind(req, v); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, v.Important["formData"])
	}))
	defer ts.Close()

	resp, err := http.DefaultClient.PostForm(ts.URL, url.Values{"important": []string{"1008"}})
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	// Output: 1008
}
