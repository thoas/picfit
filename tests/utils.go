package tests

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thoas/picfit"
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

func NewDummyProcessor() *picfit.Processor {
	processor, _ := picfit.NewProcessor(config.DefaultConfig())

	return processor
}

type Option func(*Options)

// options are server options.
type Options struct {
	Config string
}

// newOptions initializes server options.
func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// WithConfig overrides config instance.
func WithConfig(cfg string) Option {
	return func(o *Options) {
		o.Config = cfg
	}
}

type Suite struct {
	Processor *picfit.Processor
	Config    *config.Config
}

type FuncTest func(t *testing.T, s *Suite)

func Run(t *testing.T, f FuncTest, opt ...Option) {
	var (
		opts  = newOptions(opt...)
		suite *Suite
	)

	if opts.Config != "" {
		cfg, err := config.LoadFromContent(opts.Config)
		assert.Nil(t, err)

		processor, err := picfit.NewProcessor(cfg)
		assert.Nil(t, err)

		suite = &Suite{
			Config:    cfg,
			Processor: processor,
		}
	} else {
		cfg := config.DefaultConfig()

		processor, _ := picfit.NewProcessor(cfg)

		suite = &Suite{
			Config:    cfg,
			Processor: processor,
		}
	}

	f(t, suite)
}

func NewImageServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			f, err := os.Open(path.Join("tests", "fixtures", r.URL.Path))
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
