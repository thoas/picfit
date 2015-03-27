package application

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/gokvstores"
	"github.com/thoas/picfit/dummy"
	"github.com/thoas/picfit/engines"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"
)

type Dimension struct {
	Width  int
	Height int
}

type TestRequest struct {
	URL         string
	Dimensions  *Dimension
	ContentType string
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

				contentType := mime.TypeByExtension(path.Ext(r.URL.Path))

				w.WriteHeader(200)

				w.Header().Set("Content-Type", contentType)
				w.Write(bytes)
			}
		}
	}))
}

func TestStorageApplication(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	tmp := os.TempDir()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "kvstore": {
		"prefix": "picfit:",
		"type": "redis",
		"host": "127.0.0.1",
		"db": 0,
		"password": "",
		"port": 6379
	  },
	  "storage": {
		"src": {
		  "type": "fs",
		  "location": "%s"
		}
	  },
	  "shard": {
		"width": 1,
		"depth": 2
	  }
	}`

	content = fmt.Sprintf(content, tmp)

	app, err := NewFromConfig(content)

	connection := app.KVStore.Connection()
	defer connection.Close()

	assert.Nil(t, err)
	assert.NotNil(t, app.SourceStorage)

	filename := "avatar.png"

	u, _ := url.Parse(ts.URL + "/" + filename)

	location := fmt.Sprintf("http://example.com/display?url=%s&w=100&h=100&op=resize", u.String())

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	handler := app.ServeHTTP(ImageHandler)

	handler.ServeHTTP(res, request)

	assert.Equal(t, res.Code, 200)

	// We wait until the goroutine to save the file on disk is finished
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C

	etag := res.Header().Get("ETag")

	key := app.WithPrefix(etag)

	assert.True(t, connection.Exists(key))

	filepath, _ := gokvstores.String(connection.Get(key))

	assert.True(t, app.SourceStorage.Exists(filepath))
}

func TestDummyApplication(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	app := newDummyApplication()

	for _, filename := range []string{"avatar.png", "schwarzy.jpg", "giphy.gif"} {
		u, _ := url.Parse(ts.URL + "/" + filename)

		contentType := mime.TypeByExtension(path.Ext(filename))

		tests := []*TestRequest{
			&TestRequest{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=resize", u.String()),
				Dimensions: &Dimension{
					Width:  50,
					Height: 50,
				},
			},
			&TestRequest{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=100&h=30&op=resize", u.String()),
				Dimensions: &Dimension{
					Width:  100,
					Height: 30,
				},
			},
			&TestRequest{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=thumbnail", u.String()),
				Dimensions: &Dimension{
					Width:  50,
					Height: 50,
				},
			},
			&TestRequest{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=thumbnail&fmt=jpg", u.String()),
				Dimensions: &Dimension{
					Width:  50,
					Height: 50,
				},
				ContentType: "image/jpeg",
			},
		}

		for _, test := range tests {
			request, _ := http.NewRequest("GET", test.URL, nil)

			res := httptest.NewRecorder()

			handler := app.ServeHTTP(ImageHandler)

			handler.ServeHTTP(res, request)

			img, err := imaging.Decode(res.Body)

			assert.Nil(t, err)

			if test.ContentType != "" {
				assert.Equal(t, res.Header().Get("Content-Type"), test.ContentType)
			} else {
				assert.Equal(t, res.Header().Get("Content-Type"), contentType)
			}

			assert.Equal(t, res.Code, 200)

			if img.Bounds().Max.X != test.Dimensions.Width {
				t.Fatalf("Invalid width for %s: %d != %d", filename, img.Bounds().Max.X, test.Dimensions.Width)
			}

			if img.Bounds().Max.Y != test.Dimensions.Height {
				t.Fatalf("Invalid width for %s: %d != %d", filename, img.Bounds().Max.Y, test.Dimensions.Height)
			}
		}
	}
}
