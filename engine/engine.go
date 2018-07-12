package engine

import (
	"fmt"
	"strings"

	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/engine/config"
	"github.com/thoas/picfit/image"
)

type Engine struct {
	DefaultFormat  string
	Format         string
	DefaultQuality int

	backends []backend.Backend
}

const (
	goEngineType       = "go"
	lilliputEngineType = "lilliput"
)

// New initializes an Engine
func New(cfg config.Config) *Engine {
	var b []backend.Backend
	for i := range cfg.Backends {
		if cfg.Backends[i] == lilliputEngineType {
			b = append(b, backend.NewLilliput(cfg))
		} else if cfg.Backends[i] == goEngineType {
			b = append(b, &backend.GoImage{})
		}
	}

	if len(b) == 0 {
		b = append(b, &backend.GoImage{})
	}

	quality := config.DefaultQuality
	if cfg.Quality != 0 {
		quality = cfg.Quality
	}

	return &Engine{
		DefaultFormat:  cfg.DefaultFormat,
		Format:         cfg.Format,
		DefaultQuality: quality,
		backends:       b,
	}
}

func (e Engine) String() string {
	backendNames := make([]string, len(e.backends))
	for i := range e.backends {
		backendNames[i] = e.backends[i].String()
	}

	return strings.Join(backendNames, " ")
}

func (e Engine) Transform(output *image.ImageFile, operations []EngineOperation) (*image.ImageFile, error) {
	var (
		err       error
		processed []byte
		source    = output.Source
	)
	for i := range operations {
		for j := range e.backends {
			processed, err = operate(e.backends[j], output, operations[i].Operation, operations[i].Options)
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
