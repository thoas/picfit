package application_test

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/image"
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
	ctx, _ := application.Load(config.DefaultConfig())

	return ctx
}

type option func(*options)

// options are server options.
type options struct {
	Config string
}

// newOptions initializes server options.
func newOptions(opts ...option) options {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// WithConfig overrides config instance.
func WithConfig(cfg string) option {
	return func(o *options) {
		o.Config = cfg
	}
}

func Run(t *testing.T, f func(t *testing.T, ctx context.Context), opt ...option) {
	opts := newOptions(opt...)

	var ctx context.Context

	if opts.Config != "" {
		cfg, err := config.LoadFromContent(opts.Config)
		assert.Nil(t, err)

		ctx, err = application.Load(cfg)
		assert.Nil(t, err)
	} else {
		ctx = context.Background()
	}

	f(t, ctx)
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

				contentType, _ := image.MimetypeDetectorExtension(r.URL)

				w.WriteHeader(200)

				w.Header().Set("Content-Type", contentType)
				w.Write(bytes)
			}
		}
	}))
}
