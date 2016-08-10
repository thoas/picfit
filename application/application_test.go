package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/gokvstores"
	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/signature"
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

func newDummyApplication() context.Context {
	ctx, err := application.LoadFromConfig(config.DefaultConfig())

	assert.Nil(t, err)

	return ctx
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

func TestSignatureApplicationNotAuthorized(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "secret_key": "dummy"
	}`

	app, err := NewFromConfig(content)

	assert.Nil(t, err)

	u, _ := url.Parse(ts.URL + "/avatar.png")

	params := fmt.Sprintf("url=%s&w=100&h=100&op=resize", u.String())

	location := fmt.Sprintf("http://example.com/display?%s", params)

	request, _ := http.NewRequest("GET", location, nil)

	negroni := app.InitRouter()

	res := httptest.NewRecorder()

	negroni.ServeHTTP(res, request)

	assert.Equal(t, 401, res.Code)
}

func TestSignatureApplicationAuthorized(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "secret_key": "dummy"
	}`

	app, err := NewFromConfig(content)

	assert.Nil(t, err)

	u, _ := url.Parse(ts.URL + "/avatar.png")

	params := fmt.Sprintf("h=100&op=resize&url=%s&w=100", u.String())

	values, err := url.ParseQuery(params)

	assert.Nil(t, err)

	sig := signature.Sign("dummy", values.Encode())

	location := fmt.Sprintf("http://example.com/display?%s&sig=%s", params, sig)

	request, _ := http.NewRequest("GET", location, nil)

	negroni := app.InitRouter()

	res := httptest.NewRecorder()

	negroni.ServeHTTP(res, request)

	assert.Equal(t, 200, res.Code)
}

func TestUploadHandler(t *testing.T) {
	tmp := os.TempDir()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "options": {
		  "enable_upload": true
	  },
	  "storage": {
		"src": {
		  "type": "fs",
		  "location": "%s",
		  "base_url": "http://img.example.com"
		}
	  }
	}`

	content = fmt.Sprintf(content, tmp)

	app, err := NewFromConfig(content)
	assert.Nil(t, err)

	f, err := os.Open("testdata/avatar.png")
	assert.Nil(t, err)
	defer f.Close()

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	assert.Nil(t, err)

	stats, err := f.Stat()

	assert.Nil(t, err)

	fileContent, err := ioutil.ReadAll(f)

	assert.Nil(t, err)

	writer, err := w.CreateFormFile("data", "avatar.png")

	assert.Nil(t, err)

	writer.Write(fileContent)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://www.example.com/upload", body)

	assert.Nil(t, err)

	req.Header.Add("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	negroni := app.InitRouter()
	negroni.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200)

	assert.True(t, app.SourceStorage.Exists("avatar.png"))

	file, err := app.SourceStorage.Open("avatar.png")

	assert.Nil(t, err)
	assert.Equal(t, file.Size(), stats.Size())
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}

func TestDeleteHandler(t *testing.T) {
	tmp := os.TempDir()
	tmpSrcStorage := filepath.Join(tmp, "dh_src_storage")
	tmpDstStorage := filepath.Join(tmp, "dh_dst_storage")

	os.MkdirAll(tmpSrcStorage, 0755)
	os.MkdirAll(tmpDstStorage, 0755)

	img, err := ioutil.ReadFile("testdata/schwarzy.jpg")
	assert.Nil(t, err)

	// copy 5 images to src storage
	for i := 0; i < 5; i++ {
		fn := fmt.Sprintf("image%d.jpg", i+1)
		err = ioutil.WriteFile(filepath.Join(tmpSrcStorage, fn), img, 0644)
		assert.Nil(t, err)
	}

	checkDirCount := func(dir string, count int, context string) {
		dircontents, err := ioutil.ReadDir(dir)
		assert.Nil(t, err)
		assert.Equal(t, count, len(dircontents), "%s (%s)", context, dir)
	}
	checkDirCount(tmpSrcStorage, 5, "initial copy")
	checkDirCount(tmpDstStorage, 0, "initial copy")

	cfg := `
{
	"debug": true,
	"port": 3001,
	"options": {
		"enable_delete": true
	},
	"kvstore": {"type": "cache"},
	"storage": {
		"src": {
			"type": "fs",
			"location": "%s"
		},
		"dst": {
			"type": "fs",
			"location": "%s"
		}
	}
}
	`

	cfg = fmt.Sprintf(cfg, tmpSrcStorage, tmpDstStorage)
	app, err := NewFromConfig(cfg)
	assert.Nil(t, err)

	negroni := app.InitRouter()

	// generate 5 resized image1.jpg
	for i := 0; i < 5; i++ {
		// use "get" instead of "display" here to force synchronized behaviour
		url := fmt.Sprintf("http://www.example.com/get/resize/100x%d/image1.jpg", 100+i*10)
		req, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)

		res := httptest.NewRecorder()
		negroni.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 200)
	}

	// generate 2 resized image2.jpg
	for i := 0; i < 2; i++ {
		// use "get" instead of "display" here to force synchronized behaviour
		url := fmt.Sprintf("http://www.example.com/get/resize/100x%d/image2.jpg", 100+i*10)
		req, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)

		res := httptest.NewRecorder()
		negroni.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 200)
	}

	checkDirCount(tmpSrcStorage, 5, "after resize requests")
	checkDirCount(tmpDstStorage, 7, "after resize requests")
	assert.True(t, app.SourceStorage.Exists("image1.jpg"))

	// Delete image1.jpg and all of the derived images
	req, err := http.NewRequest("DELETE", "http://www.example.com/image1.jpg", nil)
	assert.Nil(t, err)

	res := httptest.NewRecorder()
	negroni.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)

	checkDirCount(tmpSrcStorage, 4, "after 1st delete request")
	checkDirCount(tmpDstStorage, 2, "after 1st delete request")
	assert.False(t, app.SourceStorage.Exists("image1.jpg"))

	// Try to delete image1.jpg again
	req, err = http.NewRequest("DELETE", "http://www.example.com/image1.jpg", nil)
	assert.Nil(t, err)

	res = httptest.NewRecorder()
	negroni.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)

	checkDirCount(tmpSrcStorage, 4, "after 2nd delete request")
	checkDirCount(tmpDstStorage, 2, "after 2nd delete request")

	assert.False(t, app.SourceStorage.Exists("image1.jpg"))

	assert.True(t, app.SourceStorage.Exists("image2.jpg"))

	// Delete image2.jpg and all of the derived images
	req, err = http.NewRequest("DELETE", "http://www.example.com/image2.jpg", nil)
	assert.Nil(t, err)

	res = httptest.NewRecorder()
	negroni.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)

	checkDirCount(tmpSrcStorage, 3, "after 3rd delete request")
	checkDirCount(tmpDstStorage, 0, "after 3rd delete request")
	assert.False(t, app.SourceStorage.Exists("image2.jpg"))
}

func TestStorageApplicationWithPath(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	tmp := os.TempDir()

	f, err := os.Open("testdata/avatar.png")
	assert.Nil(t, err)
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	// We store the image at the tmp location to access it
	// with the SourceStorage
	ioutil.WriteFile(path.Join(tmp, "avatar.png"), body, 0755)

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
		  "location": "%s",
		  "base_url": "http://img.example.com"
		}
	  },
	  "shard": {
		"width": 1,
		"depth": 2
	  }
	}`

	content = fmt.Sprintf(content, tmp)

	app, err := NewFromConfig(content)
	assert.Nil(t, err)

	negroni := app.InitRouter()

	connection := app.KVStore.Connection()
	defer connection.Close()

	location := "http://example.com/display/resize/100x100/avatar.png"

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	negroni.ServeHTTP(res, request)

	assert.Equal(t, 200, res.Code)

	// We wait until the goroutine to save the file on disk is finished
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C

	etag := res.Header().Get("ETag")

	key := app.WithPrefix(etag)

	assert.True(t, connection.Exists(key))

	filepath, _ := gokvstores.String(connection.Get(key))

	parts := strings.Split(filepath, "/")

	assert.Equal(t, len(parts), 3)
	assert.Equal(t, len(parts[0]), 1)
	assert.Equal(t, len(parts[1]), 1)

	assert.True(t, app.SourceStorage.Exists(filepath))

	location = "http://example.com/get/resize/100x100/avatar.png"

	request, _ = http.NewRequest("GET", location, nil)

	res = httptest.NewRecorder()

	negroni.ServeHTTP(res, request)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var dat map[string]interface{}

	err = json.Unmarshal(res.Body.Bytes(), &dat)

	assert.Nil(t, err)

	expected := "http://img.example.com/" + filepath

	assert.Equal(t, expected, dat["url"].(string))

	location = "http://example.com/redirect/resize/100x100/avatar.png"

	request, _ = http.NewRequest("GET", location, nil)

	res = httptest.NewRecorder()

	negroni.ServeHTTP(res, request)

	assert.Equal(t, expected, res.Header().Get("Location"))
	assert.Equal(t, 301, res.Code)
}

func TestStorageApplicationWithURL(t *testing.T) {
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
	  }
	}`

	content = fmt.Sprintf(content, tmp)

	app, err := NewFromConfig(content)
	assert.Nil(t, err)

	connection := app.KVStore.Connection()
	defer connection.Close()

	assert.NotNil(t, app.SourceStorage)
	assert.Equal(t, app.SourceStorage, app.DestStorage)

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

	parts := strings.Split(filepath, "/")

	assert.Equal(t, len(parts), 1)

	assert.True(t, app.SourceStorage.Exists(filepath))
}

func TestDummyApplicationErrors(t *testing.T) {
	app := newDummyApplication()

	location := "http://example.com/display/resize/100x100/avatar.png"

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	handler := app.ServeHTTP(ImageHandler)

	handler.ServeHTTP(res, request)
	assert.Equal(t, 400, res.Code)
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
