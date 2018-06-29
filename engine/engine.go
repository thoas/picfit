package engine

import (
	"fmt"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/imdario/mergo"

	"github.com/thoas/picfit/engine/backend"
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

	b backend.Backend
}

const (
	goEngineType       = "go"
	lilliputEngineType = "lilliput"
)

// Options is the engine options
type Options struct {
	Upscale  bool
	Format   imaging.Format
	Quality  int
	Width    int
	Height   int
	Position string
	Degree   int
}

// New initializes an Engine
func New(cfg Config) *Engine {
	var back backend.Backend
	if cfg.Type == lilliputEngineType {
		back = backend.NewLilliputEngine(cfg.MaxBufferSize)
	} else {
		back = &backend.GoImageEngine{}
	}
	return &Engine{
		DefaultFormat:  cfg.DefaultFormat,
		Format:         cfg.Format,
		DefaultQuality: cfg.Quality,
		b:              back,
	}
}

func (e *Engine) Transform(img *image.ImageFile, operation Operation, qs map[string]string) (*image.ImageFile, error) {
	err := mergo.Merge(&qs, defaultParams)

	if err != nil {
		return nil, err
	}

	var quality int
	var format string

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

	format, ok = qs["fmt"]
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

	options := &backend.Options{
		Quality: quality,
		Format:  formats[format],
	}

	var content []byte
	switch operation {
	case Noop:
		file.Processed = file.Source

		return file, err
	case Flip:
		pos, ok := qs["pos"]
		if !ok {
			return nil, fmt.Errorf("Parameter \"pos\" not found in query string")
		}

		options.Position = pos
		content, err = e.b.Flip(img, options)

	case Rotate:
		deg, err := strconv.Atoi(qs["deg"])
		if err != nil {
			return nil, err
		}

		options.Degree = deg
		content, err = e.b.Rotate(img, options)

	case Thumbnail, Resize, Fit:
		var upscale bool
		var w int
		var h int

		if upscale, err = strconv.ParseBool(qs["upscale"]); err != nil {
			return nil, err
		}

		if w, err = strconv.Atoi(qs["w"]); err != nil {
			return nil, err
		}

		if h, err = strconv.Atoi(qs["h"]); err != nil {
			return nil, err
		}

		options.Width = w
		options.Height = h
		options.Upscale = upscale

		switch operation {
		case Resize:
			content, err = e.b.Resize(img, options)
		case Thumbnail:
			content, err = e.b.Thumbnail(img, options)
		case Fit:
			content, err = e.b.Fit(img, options)
		}

	default:
		return nil, fmt.Errorf("Operation not found for %s", operation)
	}

	if err != nil {
		return nil, err
	}

	file.Processed = content

	return file, err
}
