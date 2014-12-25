package image

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/thoas/storages"
	"image"
	"math"
	"mime"
	"path"
	"strconv"
	"strings"
)

type ImageFile struct {
	Source   image.Image
	Key      string
	Headers  map[string]string
	Filepath string
	Storage  storages.Storage
}

func (i *ImageFile) ImageSize() (int, int) {
	return i.Source.Bounds().Max.X, i.Source.Bounds().Max.Y
}

func (i *ImageFile) Scale(geometry []int, upscale bool, trans Transformation) *image.NRGBA {
	width, height := i.ImageSize()

	factor := scalingFactor(width, height, geometry[0], geometry[1])

	if factor < 1 || upscale {
		width = int(float64(width) * factor)
		height = int(float64(height) * factor)

		return trans(i.Source, width, height, imaging.Lanczos)
	}

	return imaging.Clone(i.Source)
}

func (i *ImageFile) Transform(operation *Operation, qs map[string]string) (*ImageFile, error) {
	_, ok := qs["upscale"]

	if !ok {
		qs["upscale"] = "1"
	}

	switch operation {
	case Resize, Thumbnail:
		var w int
		var h int
		var err error

		if _, ok := qs["w"]; !ok {
			w = 0
		} else {
			w, err = strconv.Atoi(qs["w"])

			if err != nil {
				return nil, err
			}
		}

		if _, ok := qs["h"]; !ok {
			h = 0
		} else {
			h, err = strconv.Atoi(qs["h"])

			if err != nil {
				return nil, err
			}
		}

		upscale, err := strconv.ParseBool(qs["upscale"])

		if err != nil {
			return nil, err
		}

		dest := i.Scale([]int{w, h}, upscale, operation.Transformation)

		file := &ImageFile{
			Source:   dest,
			Key:      i.Key,
			Headers:  i.Headers,
			Filepath: i.Filepath,
		}

		return file, err
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (i *ImageFile) ToBytes() ([]byte, error) {
	format, ok := Formats[i.ContentType()]

	if !ok {
		format = DefaultFormat
	}

	return i.ToBytesWithFormat(format)
}

func (i *ImageFile) URL() string {
	return i.Storage.URL(i.Filepath)
}

func (i *ImageFile) Path() string {
	return i.Storage.Path(i.Filepath)
}

func (i *ImageFile) ToBytesWithFormat(format imaging.Format) ([]byte, error) {
	buf := &bytes.Buffer{}

	err := imaging.Encode(buf, i.Source, format)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (i *ImageFile) Format() string {
	return Extensions[i.ContentType()]
}

func (i *ImageFile) ContentType() string {
	return mime.TypeByExtension(i.FilenameExt())
}

func (i *ImageFile) Filename() string {
	return i.Filepath[strings.LastIndex(i.Filepath, "/")+1:]
}

func (i *ImageFile) FilenameExt() string {
	return path.Ext(i.Filename())
}

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}
