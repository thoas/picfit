package picfit

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/constants"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/failure"
	"github.com/thoas/picfit/image"
)

const (
	defaultDegree  = 90
	defaultHeight  = 0
	defaultUpscale = true
	defaultWidth   = 0
)

var formats = map[string]image.Format{
	"bmp":  image.BMP,
	"gif":  image.GIF,
	"jpeg": image.JPEG,
	"jpg":  image.JPEG,
	"png":  image.PNG,
	"tiff": image.TIFF,
	"webp": image.WEBP,
}

type Parameters struct {
	output     *image.ImageFile
	operations []engine.EngineOperation
}

// newParameters returns Parameters for engine.
func (p *Processor) NewParameters(ctx context.Context, input *image.ImageFile, qs map[string]interface{}) (*Parameters, error) {
	format, ok := qs["fmt"].(string)
	filepath := input.Filepath

	if ok {
		if _, ok := engine.ContentTypes[format]; !ok {
			return nil, fmt.Errorf("Unknown format %s", format)
		}

	}

	if format == "" && p.engine.Format != "" {
		format = p.engine.Format
	}

	if format == "" {
		format = input.Format()
	}

	if format == "" {
		format = p.engine.DefaultFormat
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
		opts, err := p.newBackendOptionsFromParameters(operation, qs)
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
				engineOperation.Options, err = p.newBackendOptionsFromParameters(operation, qs)
				if err != nil {
					return nil, err
				}
			} else {
				engineOperation, err = p.NewEngineOperationFromQuery(ctx, ops[i])
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
		output:     output,
		operations: operations,
	}, nil
}

func (p Processor) NewEngineOperationFromQuery(ctx context.Context, op string) (*engine.EngineOperation, error) {
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
	opts, err := p.newBackendOptionsFromParameters(operation, params)
	if err != nil {
		return nil, err
	}

	for i := range imagePaths {
		if _, err := p.sourceStorage.Stat(ctx, imagePaths[i]); errors.Is(err, gostorages.ErrNotExist) {
			return nil, errors.Wrapf(failure.ErrFileNotExists, "file does not exist: %s", imagePaths[i])
		}

		file, err := image.FromStorage(ctx, p.sourceStorage, imagePaths[i])
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load file from storage: %s", imagePaths[i])
		}
		opts.Images = append(opts.Images, *file)
	}

	return &engine.EngineOperation{
		Options:   opts,
		Operation: operation,
	}, nil
}

func (p Processor) newBackendOptionsFromParameters(operation engine.Operation, qs map[string]interface{}) (*backend.Options, error) {
	var (
		err     error
		quality = p.engine.DefaultQuality
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
			return nil, failure.ErrQuality
		}
	}

	position, ok := qs["pos"].(string)
	if !ok && operation == engine.Flip {
		return nil, fmt.Errorf("Parameter \"pos\" not found in query string")
	}

	stick, _ := qs["stick"].(string)
	if stick != "" {
		var exists bool
		for i := range constants.StickPositions {
			if stick == constants.StickPositions[i] {
				exists = true
				break
			}
		}
		if !exists {
			return nil, fmt.Errorf("Parameter \"stick\" has wrong value. Available values are: %v", constants.StickPositions)
		}
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
		Stick:    stick,
		Quality:  quality,
		Degree:   degree,
		Color:    color,
	}, nil
}
