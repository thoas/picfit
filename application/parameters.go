package application

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/errs"
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

// NewParameters returns Parameters for engine.
func NewParameters(e *engine.Engine, s gostorages.Storage, input *image.ImageFile, qs map[string]interface{}) (*Parameters, error) {
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
		opts, err := newBackendOptionsFromParameters(e, operation, qs)
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
			var err error
			engineOperation := &engine.EngineOperation{}
			operation, k := engine.Operations[ops[i]]
			if k {
				engineOperation.Operation = operation
				engineOperation.Options, err = newBackendOptionsFromParameters(e, operation, qs)
				if err != nil {
					return nil, err
				}
			} else {
				engineOperation, err = newEngineOperationFromQuery(e, s, ops[i])
				if err != nil {
					return nil, err
				}
			}

			if engineOperation != nil {
				engineOperation.Options.Format = formats[format]
				operations = append(operations, *engineOperation)
			}
		}
	}

	return &Parameters{
		Output:     output,
		Operations: operations,
	}, nil
}

func newEngineOperationFromQuery(e *engine.Engine, s gostorages.Storage, op string) (*engine.EngineOperation, error) {
	params := make(map[string]interface{})
	var imagePaths []string
	for _, p := range strings.Split(op, " ") {
		l := strings.Split(p, ":")
		if len(l) > 1 {
			if l[0] == "path" {
				imagePaths = append(imagePaths, l[1])
			} else {
				params[l[0]] = l[1]
			}
		}
	}

	op, ok := params["op"].(string)
	if !ok {
		return nil, nil
	}

	operation := engine.Operation(op)
	opts, err := newBackendOptionsFromParameters(e, operation, params)
	if err != nil {
		return nil, err
	}

	for i := range imagePaths {
		if !s.Exists(imagePaths[i]) {
			return nil, errs.ErrFileNotExists
		}

		file, err := image.FromStorage(s, imagePaths[i])
		if err != nil {
			return nil, err
		}
		opts.Images = append(opts.Images, *file)
	}

	return &engine.EngineOperation{
		Options:   opts,
		Operation: operation,
	}, nil
}

func newBackendOptionsFromParameters(e *engine.Engine, operation engine.Operation, qs map[string]interface{}) (*backend.Options, error) {
	var (
		err     error
		quality = e.DefaultQuality
		upscale = defaultUpscale
		height  = defaultHeight
		width   = defaultWidth
		degree  = defaultDegree
	)

	q, ok := qs["q"].(string)
	if ok {
		quality, err = strconv.Atoi(q)
		if err != nil {
			return nil, err
		}

		if quality > 100 {
			return nil, fmt.Errorf("Quality should be <= 100")
		}
	}

	position, ok := qs["pos"].(string)
	if !ok && operation == engine.Flip {
		return nil, fmt.Errorf("Parameter \"pos\" not found in query string")
	}

	color, _ := qs["color"].(string)

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
		Color:    color,
	}, nil
}
