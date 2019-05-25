package picfit_test

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

	conv "github.com/cstockton/go-conv"

	"github.com/disintegration/imaging"

	"github.com/stretchr/testify/assert"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/server"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/picfit/tests"
)

func TestSignatureApplicationNotAuthorized(t *testing.T) {
	ts := tests.NewImageServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "secret_key": "dummy"
	}`

	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		u, _ := url.Parse(ts.URL + "/avatar.png")

		params := fmt.Sprintf("url=%s&w=100&h=100&op=resize", u.String())

		location := fmt.Sprintf("http://example.com/display?%s", params)

		request, _ := http.NewRequest("GET", location, nil)

		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, request)

		assert.Equal(t, 401, res.Code)
	}, tests.WithConfig(content))
}

func TestSignatureApplicationAuthorized(t *testing.T) {
	ts := tests.NewImageServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "secret_key": "dummy"
	}`

	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		u, _ := url.Parse(ts.URL + "/avatar.png")

		params := fmt.Sprintf("h=100&op=resize&url=%s&w=100", u.String())

		values, err := url.ParseQuery(params)
		assert.Nil(t, err)

		sig := signature.Sign("dummy", values.Encode())

		location := fmt.Sprintf("http://example.com/display?%s&sig=%s", params, sig)

		request, _ := http.NewRequest("GET", location, nil)

		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, request)

		assert.Equal(t, 200, res.Code)
	}, tests.WithConfig(content))

}

func TestSizeRestrictedApplicationNotAuthorized(t *testing.T) {
	ts := tests.NewImageServer()
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

	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		u, _ := url.Parse(ts.URL + "/avatar.png")

		// unallowed size
		params := fmt.Sprintf("url=%s&w=50&h=50&op=resize", u.String())

		location := fmt.Sprintf("http://example.com/display?%s", params)

		request, _ := http.NewRequest("GET", location, nil)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, request)

		assert.Equal(t, 403, res.Code)

		// allowed size
		params = fmt.Sprintf("url=%s&w=100&h=100&op=resize", u.String())

		location = fmt.Sprintf("http://example.com/display?%s", params)

		request, _ = http.NewRequest("GET", location, nil)

		res = httptest.NewRecorder()

		server.ServeHTTP(res, request)

		assert.Equal(t, 200, res.Code)
	}, tests.WithConfig(content))
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

	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		f, err := os.Open("tests/fixtures/avatar.png")
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

		server.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)

		assert.True(t, suite.Processor.FileExists("avatar.png"))

		file, err := suite.Processor.OpenFile("avatar.png")

		assert.Nil(t, err)
		assert.Equal(t, file.Size(), stats.Size())
		assert.Equal(t, "application/json; charset=utf-8", res.Header().Get("Content-Type"))
	}, tests.WithConfig(content))
}

func TestDeleteHandler(t *testing.T) {
	tmp, err := ioutil.TempDir("", tests.RandString(10))

	assert.Nil(t, err)

	tmpSrcStorage := filepath.Join(tmp, "src")
	tmpDstStorage := filepath.Join(tmp, "dst")

	os.MkdirAll(tmpSrcStorage, 0755)
	os.MkdirAll(tmpDstStorage, 0755)

	img, err := ioutil.ReadFile("tests/fixtures/schwarzy.jpg")
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
		"enable_delete": true,
		"enable_cascade_delete": true
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
	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		// generate 5 resized image1.jpg
		for i := 0; i < 5; i++ {
			// use "get" instead of "display" here to force synchronized behaviour
			url := fmt.Sprintf("http://www.example.com/get/resize/100x%d/image1.jpg", 100+i*10)
			req, err := http.NewRequest("GET", url, nil)
			assert.Nil(t, err)

			res := httptest.NewRecorder()
			server.ServeHTTP(res, req)
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
			server.ServeHTTP(res, req)
			assert.Equal(t, res.Code, 200)
		}

		checkDirCount(tmpSrcStorage, 5, "after resize requests")
		checkDirCount(tmpDstStorage, 7, "after resize requests")
		assert.True(t, suite.Processor.FileExists("image1.jpg"))

		// Delete image1.jpg and all of the derived images
		req, err := http.NewRequest("DELETE", "http://www.example.com/image1.jpg", nil)
		assert.Nil(t, err)

		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 200)

		checkDirCount(tmpSrcStorage, 4, "after 1st delete request")
		checkDirCount(tmpDstStorage, 2, "after 1st delete request")
		assert.False(t, suite.Processor.FileExists("image1.jpg"))

		// Try to delete image1.jpg again
		req, err = http.NewRequest("DELETE", "http://www.example.com/image1.jpg", nil)
		assert.Nil(t, err)

		res = httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 404)

		checkDirCount(tmpSrcStorage, 4, "after 2nd delete request")
		checkDirCount(tmpDstStorage, 2, "after 2nd delete request")

		assert.False(t, suite.Processor.FileExists("image1.jpg"))

		assert.True(t, suite.Processor.FileExists("image2.jpg"))

		// Delete image2.jpg and all of the derived images
		req, err = http.NewRequest("DELETE", "http://www.example.com/image2.jpg", nil)
		assert.Nil(t, err)

		res = httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assert.Equal(t, res.Code, 200)

		checkDirCount(tmpSrcStorage, 3, "after 3rd delete request")
		checkDirCount(tmpDstStorage, 0, "after 3rd delete request")
		assert.False(t, suite.Processor.FileExists("image2.jpg"))
	}, tests.WithConfig(cfg))
}

func TestStorageApplicationWithPath(t *testing.T) {
	ts := tests.NewImageServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	tmp := os.TempDir()

	f, err := os.Open("tests/fixtures/avatar.png")
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
		"redis": {
			"host": "127.0.0.1",
			"db": 0,
			"password": "",
			"port": 6379
		}
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

	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		location := "http://example.com/display/resize/100x100/avatar.png"

		request, _ := http.NewRequest("GET", location, nil)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, request)

		assert.Equal(t, 200, res.Code)

		// We wait until the goroutine to save the file on disk is finished
		timer1 := time.NewTimer(time.Second * 2)
		<-timer1.C

		etag := res.Header().Get("ETag")

		exists, err := suite.Processor.KeyExists(etag)
		assert.Nil(t, err)
		assert.True(t, exists)

		raw, err := suite.Processor.GetKey(etag)
		assert.Nil(t, err)

		filepath, err := conv.String(raw)
		assert.Nil(t, err)

		parts := strings.Split(filepath, "/")

		assert.Equal(t, len(parts), 3)
		assert.Equal(t, len(parts[0]), 1)
		assert.Equal(t, len(parts[1]), 1)

		assert.True(t, suite.Processor.FileExists(filepath))

		location = "http://example.com/get/resize/100x100/avatar.png"

		request, _ = http.NewRequest("GET", location, nil)

		res = httptest.NewRecorder()

		server.ServeHTTP(res, request)

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

		server.ServeHTTP(res, request)

		assert.Equal(t, expected, res.Header().Get("Location"))
		assert.Equal(t, 301, res.Code)

	}, tests.WithConfig(content))
}

func TestStorageApplicationWithURL(t *testing.T) {
	ts := tests.NewImageServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	tmp := os.TempDir()

	content := `{
	  "debug": true,
	  "port": 3001,
	  "kvstore": {
		"prefix": "picfit:",
		"type": "redis",
		"redis": {
			"host": "127.0.0.1",
			"db": 0,
			"password": "",
			"port": 6379
		}
	  },
	  "storage": {
		"src": {
		  "type": "fs",
		  "location": "%s"
		}
	  }
	}`

	content = fmt.Sprintf(content, tmp)

	tests.Run(t, func(t *testing.T, suite *tests.Suite) {
		server, err := server.New(suite.Config)
		assert.Nil(t, err)

		filename := "avatar.png"

		u, _ := url.Parse(ts.URL + "/" + filename)

		location := fmt.Sprintf("http://example.com/display?url=%s&w=100&h=100&op=resize", u.String())

		request, _ := http.NewRequest("GET", location, nil)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, request)

		assert.Equal(t, res.Code, 200)

		// We wait until the goroutine to save the file on disk is finished
		timer1 := time.NewTimer(time.Second * 2)
		<-timer1.C

		etag := res.Header().Get("ETag")

		exists, err := suite.Processor.KeyExists(etag)
		assert.Nil(t, err)
		assert.True(t, exists)

		raw, err := suite.Processor.GetKey(etag)
		assert.Nil(t, err)

		filepath, err := conv.String(raw)
		assert.Nil(t, err)

		parts := strings.Split(filepath, "/")

		assert.Equal(t, len(parts), 1)

		assert.True(t, suite.Processor.FileExists(filepath))
	}, tests.WithConfig(content))
}

func TestDummyApplicationErrors(t *testing.T) {
	location := "http://example.com/display/resize/100x100/avatar.png"

	request, _ := http.NewRequest("GET", location, nil)

	res := httptest.NewRecorder()

	server, err := server.New(config.DefaultConfig())
	assert.Nil(t, err)

	server.ServeHTTP(res, request)
	assert.Equal(t, 404, res.Code)
}

func TestDummyApplication(t *testing.T) {
	ts := tests.NewImageServer()
	defer ts.Close()
	defer ts.CloseClientConnections()

	server, err := server.New(config.DefaultConfig())
	assert.Nil(t, err)

	for _, filename := range []string{"avatar.png", "schwarzy.jpg", "giphy.gif"} {
		u, _ := url.Parse(ts.URL + "/" + filename)

		contentType := mime.TypeByExtension(path.Ext(filename))

		tests := []*tests.TestRequest{
			{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=resize", u.String()),
				Dimensions: &tests.Dimension{
					Width:  50,
					Height: 50,
				},
			},
			{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=100&h=30&op=resize", u.String()),
				Dimensions: &tests.Dimension{
					Width:  100,
					Height: 30,
				},
			},
			{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=thumbnail", u.String()),
				Dimensions: &tests.Dimension{
					Width:  50,
					Height: 50,
				},
			},
			{
				URL: fmt.Sprintf("http://example.com/display?url=%s&w=50&h=50&op=thumbnail&fmt=jpg", u.String()),
				Dimensions: &tests.Dimension{
					Width:  50,
					Height: 50,
				},
				ContentType: "image/jpeg",
			},
			{
				URL: fmt.Sprintf("http://example.com/display?url=%s&op=op:resize+w:100+h:50&op=op:rotate+deg:90", u.String()),
				Dimensions: &tests.Dimension{
					Width:  50,
					Height: 100,
				},
			},
			{
				URL: fmt.Sprintf("http://example.com/display?url=%s&op=resize&w=100&h=50&op=op:rotate+deg:90", u.String()),
				Dimensions: &tests.Dimension{
					Width:  50,
					Height: 100,
				},
			},
		}

		for _, test := range tests {
			request, _ := http.NewRequest("GET", test.URL, nil)

			res := httptest.NewRecorder()

			server.ServeHTTP(res, request)

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
