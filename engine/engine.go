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

	backends  []backend.Backend
	mimetypes map[string]backend.Backend
}

// New initializes an Engine
func New(cfg config.Config) *Engine {
	var (
		b         []backend.Backend
		mimetypes = map[string]backend.Backend{}
	)

	if cfg.Backends == nil {
		b = append(b, &backend.GoImage{})
	} else {
		if cfg.Backends.Vips != nil {
			back := backend.NewVips(cfg)

			b = append(b, back)

			for _, mimetype := range cfg.Backends.Vips.Mimetypes {
				mimetypes[mimetype] = back
			}
		}

		if cfg.Backends.Lilliput != nil {
			back := backend.NewLilliput(cfg)

			b = append(b, back)

			for _, mimetype := range cfg.Backends.Lilliput.Mimetypes {
				mimetypes[mimetype] = back
			}
		}

		if cfg.Backends.GoImage != nil {
			back := &backend.GoImage{}

			b = append(b, back)

			for _, mimetype := range cfg.Backends.GoImage.Mimetypes {
				mimetypes[mimetype] = back
			}
		}
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
		mimetypes:      mimetypes,
	}
}

func (e Engine) String() string {
	backendNames := []string{}
	for _, backend := range e.backends {
		backendNames = append(backendNames, backend.String())
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
		backends := e.backends

		back, ok := e.mimetypes[output.ContentType()]
		if ok {
			backends = []backend.Backend{back}
		}

		for j := range backends {
			processed, err = operate(backends[j], output, operations[i].Operation, operations[i].Options)
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
