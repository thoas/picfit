package application

import (
	"fmt"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/image"
)

const (
	defaultUpscale = true
	defaultWidth   = 0
	defaultHeight  = 0
	defaultDegree  = 90
)

var formats = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"jpg":  imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
}

type Parameters struct {
	Output     *image.ImageFile
	Operations []engine.EngineOperation
}

func NewParameters(e *engine.Engine, input *image.ImageFile, qs map[string]interface{}) (*Parameters, error) {
	format, ok := qs["fmt"].(string)
	filepath := input.Filepath

	if ok {
		if _, ok := engine.ContentTypes[format]; !ok {
			return nil, fmt.Errorf("Unknown format %s", format)
		}

	}

	if format == "" && e.Format != "" {
		format = e.Format
	}

	if format == "" {
		format = input.Format()
	}

	if format == "" {
		format = e.DefaultFormat
	}

	if format != input.Format() {
		index := len(filepath) - len(input.Format())

		filepath = filepath[:index] + format

		if contentType, ok := engine.ContentTypes[format]; ok {
			input.Headers["Content-Type"] = contentType
		}
	}

	output := &image.ImageFile{
		Source:   input.Source,
		Key:      input.Key,
		Headers:  input.Headers,
		Filepath: filepath,
	}

	var operations []engine.EngineOperation

	op, ok := qs["op"].(string)
	if ok {
		operation := engine.Operation(op)
		opts, err := newBackendOptions(e, operation, qs)
		if err != nil {
			return nil, err
		}

		opts.Format = formats[format]
		operations = append(operations, engine.EngineOperation{
			Options:   opts,
			Operation: operation,
		})
	}

	ops, ok := qs["op"].([]string)
	if ok {
		for i := range ops {
			operation := engine.Operation(ops[i])
			opts, err := newBackendOptions(e, operation, qs)
			if err != nil {
				return nil, err
			}

			opts.Format = formats[format]
			operations = append(operations, engine.EngineOperation{
				Options:   opts,
				Operation: operation,
			})
		}
	}

	return &Parameters{
		Output:     output,
		Operations: operations,
	}, nil
}

func newBackendOptions(e *engine.Engine, operation engine.Operation, qs map[string]interface{}) (*backend.Options, error) {
	var (
		err     error
		quality int
		upscale = defaultUpscale
		height  = defaultHeight
		width   = defaultWidth
		degree  = defaultDegree
	)

	q, ok := qs["q"].(string)
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

	position, ok := qs["pos"].(string)
	if !ok && operation == engine.Flip {
		return nil, fmt.Errorf("Parameter \"pos\" not found in query string")
	}

	if deg, ok := qs["deg"].(string); ok {
		degree, err = strconv.Atoi(deg)
		if err != nil {
			return nil, err
		}
	}

	if up, ok := qs["upscale"].(string); ok {
		upscale, err = strconv.ParseBool(up)
		if err != nil {
			return nil, err
		}
	}

	if w, ok := qs["w"].(string); ok {
		width, err = strconv.Atoi(w)
		if err != nil {
			return nil, err
		}
	}

	if h, ok := qs["h"].(string); ok {
		height, err = strconv.Atoi(h)
		if err != nil {
			return nil, err
		}
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
