package application

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/imaging"
	"github.com/thoas/picfit/dummy"
	"github.com/thoas/picfit/engines"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
)

type Dimension struct {
	Width  int
	Height int
}

type TestRequest struct {
	URL        string
	Dimensions *Dimension
}

func newDummyApplication() *Application {
	app := NewApplication()
	app.SourceStorage = &dummy.DummyStorage{}
	app.DestStorage = &dummy.DummyStorage{}
	app.KVStore = &dummy.DummyKVStore{}
	app.Engine = &engines.GoImageEngine{}

	return app
}

func newHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			f, err := os.Open(path.Join("testdata", r.URL.Path))
			defer f.Close()

			if err != nil {
				w.WriteHeader(500)
			} else {
				bytes, _ := ioutil.ReadAll(f)

				w.Header().Set("Content-Length", fmt.Sprintf("%d\n\n%v", len(bytes), bytes))
				w.Write(bytes)
				w.WriteHeader(200)
			}
		}
	}))
}

func TestDummyApplication(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()

	app := newDummyApplication()

	u, _ := url.Parse(ts.URL + "/avatar.png")

	tests := []*TestRequest{
		&TestRequest{
			URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=resize", u.String()),
			Dimensions: &Dimension{
				Width:  50,
				Height: 50,
			},
		},
		&TestRequest{
			URL: fmt.Sprintf("http://example.com/display?url=%s&w=100&h=0&op=resize", u.String()),
			Dimensions: &Dimension{
				Width:  100,
				Height: 100,
			},
		},
		&TestRequest{
			URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=thumbnail", u.String()),
			Dimensions: &Dimension{
				Width:  50,
				Height: 50,
			},
		},
	}

	for _, test := range tests {
		request, _ := http.NewRequest("GET", test.URL, nil)

		res := httptest.NewRecorder()

		handler := app.ServeHTTP(ImageHandler)

		handler.ServeHTTP(res, request)

		img, err := imaging.Decode(res.Body)

		assert.Nil(t, err)
		assert.Equal(t, img.Bounds().Max.X, test.Dimensions.Width)
		assert.Equal(t, img.Bounds().Max.Y, test.Dimensions.Height)
	}
}
