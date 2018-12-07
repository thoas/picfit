package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/imdario/mergo"

	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/engine/config"
	"github.com/thoas/picfit/image"
)

var defaultParams = map[string]string{
	"upscale": "1",
	"h":       "0",
	"w":       "0",
	"deg":     "90",
}

var formats = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"jpg":  imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
}

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

	return &Engine{
		DefaultFormat:  cfg.DefaultFormat,
		Format:         cfg.Format,
		DefaultQuality: cfg.Quality,
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

func (e Engine) Transform(img *image.ImageFile, operation Operation, qs map[string]string) (*image.ImageFile, error) {
	err := mergo.Merge(&qs, defaultParams)

	if err != nil {
		return nil, err
	}

	format, ok := qs["fmt"]
	filepath := img.Filepath

	if ok {
		if _, ok := ContentTypes[format]; !ok {
			return nil, fmt.Errorf("Unknown format %s", format)
		}

	}

	if format == "" && e.Format != "" {
		format = e.Format
	}

	if format == "" {
		format = img.Format()
	}

	if format == "" {
		format = e.DefaultFormat
	}

	if format != img.Format() {
		index := len(filepath) - len(img.Format())

		filepath = filepath[:index] + format

		if contentType, ok := ContentTypes[format]; ok {
			img.Headers["Content-Type"] = contentType
		}
	}

	file := &image.ImageFile{
		Source:   img.Source,
		Key:      img.Key,
		Headers:  img.Headers,
		Filepath: filepath,
	}

	options, err := newBackendOptions(e, operation, qs)
	if err != nil {
		return nil, err
	}
	options.Format = formats[format]

	for i := range e.backends {
		file.Processed, err = operate(e.backends[i], img, operation, options)
		if err == nil {
			break
		}
		if err != backend.MethodNotImplementedError {
			return nil, err
		}
	}

	return file, err
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
	default:
		return nil, fmt.Errorf("Operation not found for %s", operation)
	}
}

func newBackendOptions(e Engine, operation Operation, qs map[string]string) (*backend.Options, error) {
	var quality int
	q, ok := qs["q"]
	if ok {
		quality, err := strconv.Atoi(q)

		if err != nil {
			return nil, err
		}

		if quality > 100 {
			return nil, fmt.Errorf("Quality should be <= 100")
		}
	} else {
		quality = e.DefaultQuality
	}

	position, ok := qs["pos"]
	if !ok && operation == Flip {
		return nil, fmt.Errorf("Parameter \"pos\" not found in query string")
	}

	degree, err := strconv.Atoi(qs["deg"])
	if err != nil {
		return nil, err
	}

	upscale, err := strconv.ParseBool(qs["upscale"])
	if err != nil {
		return nil, err
	}

	width, err := strconv.Atoi(qs["w"])
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(qs["h"])
	if err != nil {
		return nil, err
	}

	return &backend.Options{
		Width:    width,
		Height:   height,
		Upscale:  upscale,
		Position: position,
		Quality:  quality,
		Degree:   degree,
	}, nil
}
