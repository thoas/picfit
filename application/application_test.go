package application_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime"
	"mime/multipart"
	"net"
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

	"github.com/alicebob/miniredis"
	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/gokvstores"
	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/kvstore"
	"github.com/thoas/picfit/server"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/picfit/storage"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func newDummyApplication() context.Context {
	ctx, _ := application.LoadFromConfig(config.DefaultConfig())

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

	ctx, err := application.LoadFromConfigContent(content)

	assert.Nil(t, err)

	u, _ := url.Parse(ts.URL + "/avatar.png")

	params := fmt.Sprintf("url=%s&w=100&h=100&op=resize", u.String())

	location := fmt.Sprintf("http://example.com/display?%s", params)

	request, _ := http.NewRequest("GET", location, nil)

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, request)

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

	ctx, err := application.LoadFromConfigContent(content)

	assert.Nil(t, err)

	u, _ := url.Parse(ts.URL + "/avatar.png")

	params := fmt.Sprintf("h=100&op=resize&url=%s&w=100", u.String())

	values, err := url.ParseQuery(params)

	assert.Nil(t, err)

	sig := signature.Sign("dummy", values.Encode())

	location := fmt.Sprintf("http://example.com/display?%s&sig=%s", params, sig)

	request, _ := http.NewRequest("GET", location, nil)

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, request)

	assert.Equal(t, 200, res.Code)
}

func TestSizeRestrictedApplicationNotAuthorized(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "options": {
	    "allowed_sizes": [
	      {"width": 100, "height": 100}
	    ]
	  }
	}`

	ctx, err := application.LoadFromConfigContent(content)

	assert.Nil(t, err)

	u, _ := url.Parse(ts.URL + "/avatar.png")

	router, err := server.Router(ctx)

	// unallowed size
	params := fmt.Sprintf("url=%s&w=50&h=50&op=resize", u.String())

	location := fmt.Sprintf("http://example.com/display?%s", params)

	request, _ := http.NewRequest("GET", location, nil)

	assert.Nil(t, err)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, request)

	assert.Equal(t, 403, res.Code)

	// allowed size
	params = fmt.Sprintf("url=%s&w=100&h=100&op=resize", u.String())

	location = fmt.Sprintf("http://example.com/display?%s", params)

	request, _ = http.NewRequest("GET", location, nil)

	assert.Nil(t, err)

	res = httptest.NewRecorder()

	router.ServeHTTP(res, request)

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

	ctx, err := application.LoadFromConfigContent(content)
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

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)

	sourceStorage := storage.SourceFromContext(ctx)

	assert.True(t, sourceStorage.Exists("avatar.png"))

	file, err := sourceStorage.Open("avatar.png")

	assert.Nil(t, err)
	assert.Equal(t, file.Size(), stats.Size())
	assert.Equal(t, "application/json; charset=utf-8", res.Header().Get("Content-Type"))
}

func TestDeleteHandler(t *testing.T) {
	tmp, err := ioutil.TempDir("", RandString(10))

	assert.Nil(t, err)

	tmpSrcStorage := filepath.Join(tmp, "src")
	tmpDstStorage := filepath.Join(tmp, "dst")

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
	ctx, err := application.LoadFromConfigContent(cfg)
	assert.Nil(t, err)

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	// generate 5 resized image1.jpg
	for i := 0; i < 5; i++ {
		// use "get" instead of "display" here to force synchronized behaviour
		url := fmt.Sprintf("http://www.example.com/get/resize/100x%d/image1.jpg", 100+i*10)
		req, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 200)
	}

	checkDirCount(tmpDstStorage, 5, "after resize requests")

	// generate 2 resized image2.jpg
	for i := 0; i < 2; i++ {
		// use "get" instead of "display" here to force synchronized behaviour
		url := fmt.Sprintf("http://www.example.com/get/resize/100x%d/image2.jpg", 100+i*10)
		req, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 200)
	}

	sourceStorage := storage.SourceFromContext(ctx)

	checkDirCount(tmpSrcStorage, 5, "after resize requests")
	checkDirCount(tmpDstStorage, 7, "after resize requests")
	assert.True(t, sourceStorage.Exists("image1.jpg"))

	// Delete image1.jpg and all of the derived images
	req, err := http.NewRequest("DELETE", "http://www.example.com/image1.jpg", nil)
	assert.Nil(t, err)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)

	checkDirCount(tmpSrcStorage, 4, "after 1st delete request")
	checkDirCount(tmpDstStorage, 2, "after 1st delete request")
	assert.False(t, sourceStorage.Exists("image1.jpg"))

	// Try to delete image1.jpg again
	req, err = http.NewRequest("DELETE", "http://www.example.com/image1.jpg", nil)
	assert.Nil(t, err)

	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 404)

	checkDirCount(tmpSrcStorage, 4, "after 2nd delete request")
	checkDirCount(tmpDstStorage, 2, "after 2nd delete request")

	assert.False(t, sourceStorage.Exists("image1.jpg"))

	assert.True(t, sourceStorage.Exists("image2.jpg"))

	// Delete image2.jpg and all of the derived images
	req, err = http.NewRequest("DELETE", "http://www.example.com/image2.jpg", nil)
	assert.Nil(t, err)

	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)

	checkDirCount(tmpSrcStorage, 3, "after 3rd delete request")
	checkDirCount(tmpDstStorage, 0, "after 3rd delete request")
	assert.False(t, sourceStorage.Exists("image2.jpg"))
}

func TestStorageApplicationWithPath(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	rs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start redis server: %v", err)
	}
	defer rs.Close()
	rHost, rPort, err := net.SplitHostPort(rs.Addr())
	if err != nil {
		t.Fatalf("Failed to parse redis addr %q: %v", rs.Addr(), err)
	}

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
			"host": "%s",
			"db": 0,
			"password": "",
			"port": %s
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

	content = fmt.Sprintf(content, rHost, rPort, tmp)

	ctx, err := application.LoadFromConfigContent(content)
	assert.Nil(t, err)

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	store := kvstore.FromContext(ctx)

	connection := store.Connection()
	defer connection.Close()

	location := "http://example.com/display/resize/100x100/avatar.png"

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, request)

	assert.Equal(t, 200, res.Code)

	// We wait until the goroutine to save the file on disk is finished
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C

	etag := res.Header().Get("ETag")

	cfg := config.FromContext(ctx)

	key := cfg.KVStore.Prefix + etag

	assert.True(t, connection.Exists(key))

	filepath, _ := gokvstores.String(connection.Get(key))

	parts := strings.Split(filepath, "/")

	assert.Equal(t, len(parts), 3)
	assert.Equal(t, len(parts[0]), 1)
	assert.Equal(t, len(parts[1]), 1)

	sourceStorage := storage.SourceFromContext(ctx)

	assert.True(t, sourceStorage.Exists(filepath))

	location = "http://example.com/get/resize/100x100/avatar.png"

	request, _ = http.NewRequest("GET", location, nil)

	res = httptest.NewRecorder()

	router.ServeHTTP(res, request)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "application/json; charset=utf-8", res.Header().Get("Content-Type"))

	var dat map[string]interface{}

	err = json.Unmarshal(res.Body.Bytes(), &dat)

	assert.Nil(t, err)

	expected := "http://img.example.com/" + filepath

	assert.Equal(t, expected, dat["url"].(string))

	location = "http://example.com/redirect/resize/100x100/avatar.png"

	request, _ = http.NewRequest("GET", location, nil)

	res = httptest.NewRecorder()

	router.ServeHTTP(res, request)

	assert.Equal(t, expected, res.Header().Get("Location"))
	assert.Equal(t, 301, res.Code)
}

func TestStorageApplicationWithURL(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	rs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start redis server: %v", err)
	}
	defer rs.Close()
	rHost, rPort, err := net.SplitHostPort(rs.Addr())
	if err != nil {
		t.Fatalf("Failed to parse redis addr %q: %v", rs.Addr(), err)
	}

	tmp := os.TempDir()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "kvstore": {
			"prefix": "picfit:",
			"type": "redis",
			"host": "%s",
			"db": 0,
			"password": "",
			"port": %s
	  },
	  "storage": {
			"src": {
				"type": "fs",
				"location": "%s"
			}
	  }
	}`

	content = fmt.Sprintf(content, rHost, rPort, tmp)

	ctx, err := application.LoadFromConfigContent(content)
	assert.Nil(t, err)

	connection := kvstore.FromContext(ctx).Connection()
	defer connection.Close()

	sourceStorage := storage.SourceFromContext(ctx)

	assert.NotNil(t, sourceStorage)
	assert.Equal(t, sourceStorage, storage.DestinationFromContext(ctx))

	filename := "avatar.png"

	u, _ := url.Parse(ts.URL + "/" + filename)

	location := fmt.Sprintf("http://example.com/display?url=%s&w=100&h=100&op=resize", u.String())

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	router.ServeHTTP(res, request)

	assert.Equal(t, res.Code, 200)

	// We wait until the goroutine to save the file on disk is finished
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C

	etag := res.Header().Get("ETag")

	key := config.FromContext(ctx).KVStore.Prefix + etag

	assert.True(t, connection.Exists(key))

	filepath, _ := gokvstores.String(connection.Get(key))

	parts := strings.Split(filepath, "/")

	assert.Equal(t, len(parts), 1)

	assert.True(t, sourceStorage.Exists(filepath))
}

func TestDummyApplicationErrors(t *testing.T) {
	ctx := newDummyApplication()

	location := "http://example.com/display/resize/100x100/avatar.png"

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	router, err := server.Router(ctx)

	assert.Nil(t, err)

	router.ServeHTTP(res, request)
	assert.Equal(t, 404, res.Code)
}

func TestDummyApplication(t *testing.T) {
	ts := newHTTPServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	ctx := newDummyApplication()

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

		router, err := server.Router(ctx)

		assert.Nil(t, err)

		for _, test := range tests {
			request, _ := http.NewRequest("GET", test.URL, nil)

			res := httptest.NewRecorder()

			router.ServeHTTP(res, request)

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
