package engine

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/engine/config"
	"github.com/thoas/picfit/image"
)

type Engine struct {
	DefaultFormat  string
	DefaultQuality int
	Format         string
	backends       []*BackendWrapper
}

type BackendWrapper struct {
	backend   backend.Backend
	mimetypes []string
	weight    int
}

// New initializes an Engine
func New(cfg config.Config) *Engine {
	var b []*BackendWrapper

	if cfg.Backends == nil {
		b = append(b, &BackendWrapper{
			backend:   &backend.GoImage{},
			mimetypes: MimeTypes,
		})
	} else {
		if cfg.Backends.Gifsicle != nil {
			path := cfg.Backends.Gifsicle.Path
			if path == "" {
				path = "gifsicle"
			}

			if _, err := exec.LookPath(path); err == nil {
				b = append(b, &BackendWrapper{
					backend:   &backend.Gifsicle{Path: path},
					mimetypes: cfg.Backends.Gifsicle.Mimetypes,
					weight:    cfg.Backends.Gifsicle.Weight,
				})
			}
		}
		if cfg.Backends.GoImage != nil {
			b = append(b, &BackendWrapper{
				backend:   &backend.GoImage{},
				mimetypes: cfg.Backends.GoImage.Mimetypes,
				weight:    cfg.Backends.GoImage.Weight,
			})
		}
	}

	sort.Slice(b, func(i, j int) bool {
		return b[i].weight < b[j].weight
	})

	quality := config.DefaultQuality
	if cfg.Quality != 0 {
		quality = cfg.Quality
	}

	return &Engine{
		DefaultFormat:  cfg.DefaultFormat,
		DefaultQuality: quality,
		Format:         cfg.Format,
		backends:       b,
	}
}

func (e Engine) String() string {
	backendNames := []string{}
	for _, backend := range e.backends {
		backendNames = append(backendNames, backend.backend.String())
	}

	return strings.Join(backendNames, " ")
}

func (e Engine) Transform(output *image.ImageFile, operations []EngineOperation) (*image.ImageFile, error) {
	var (
		err       error
		processed []byte
		source    = output.Source
	)

	ct := output.ContentType()
	for i := range operations {
		for j := range e.backends {
			var processing bool
			for k := range e.backends[j].mimetypes {
				if ct == e.backends[j].mimetypes[k] {
					processing = true
					break
				}
			}

			if !processing {
				continue
			}

			processed, err = operate(e.backends[j].backend, output, operations[i].Operation, operations[i].Options)
			if err == nil {
				output.Source = processed
				break
			}
			if err != backend.MethodNotImplementedError {
				return nil, err
			}
		}
	}

	output.Source = source
	output.Processed = processed

	return output, err
}

func operate(b backend.Backend, img *image.ImageFile, operation Operation, options *backend.Options) ([]byte, error) {
	switch operation {
	case Noop:
		return img.Source, nil
	case Flip:
		return b.Flip(img, options)
	case Rotate:
		return b.Rotate(img, options)
	case Resize:
		return b.Resize(img, options)
	case Thumbnail:
		return b.Thumbnail(img, options)
	case Fit:
		return b.Fit(img, options)
	case Flat:
		return b.Flat(img, options)
	default:
		return nil, fmt.Errorf("Operation not found for %s", operation)
	}
}
