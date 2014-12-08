package image

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/thoas/storages"
	"image"
	"math"
	"mime"
	"net/url"
	"strconv"
	"strings"
)

type ImageFile struct {
	Source   image.Image
	Key      string
	Header   map[string]string
	Filepath string
	Storage  storages.Storage
}

type Transformation func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA

func (i *ImageFile) GetImageSize() (int, int) {
	return i.Source.Bounds().Max.X, i.Source.Bounds().Max.Y
}

func (i *ImageFile) scale(width int, height int, trans Transformation) *image.NRGBA {
	return trans(i.Source, width, height, imaging.Lanczos)
}

func (i *ImageFile) Scale(geometry []int, upscale bool, trans Transformation) *image.NRGBA {
	width, height := i.GetImageSize()

	factor := scalingFactor(width, height, geometry[0], geometry[1])

	if factor < 1 || upscale {
		width = int(float64(width) * factor)
		height = int(float64(height) * factor)

		return i.scale(width, height, trans)
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
			Header:   i.Header,
			Filepath: i.Filepath,
		}

		return file, err
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (i *ImageFile) ToBytes() ([]byte, error) {
	format, ok := Formats[i.GetContentType()]

	if !ok {
		format = DefaultFormat
	}

	return i.ToBytesWithFormat(format)
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
	return Extensions[i.GetContentType()]
}

func (i *ImageFile) GetContentType() string {
	return mime.TypeByExtension(i.GetFilename())
}

func (i *ImageFile) GetFilename() string {
	return i.Filepath[strings.LastIndex(i.Filepath, "/")+1:]
}

func (i *ImageFile) LoadFromStorage(storage storages.Storage) (*ImageFile, error) {
	var file *ImageFile
	var err error

	// URL provided we use http protocol to retrieve it
	if storage.HasBaseURL() {
		u, err := url.Parse(storage.URL(i.Filepath))

		if err != nil {
			return nil, err
		}

		file, err = ImageFileFromURL(u)

		if err != nil {
			return nil, err
		}
	} else {
		body, err := storage.Open(i.Filepath)

		if err != nil {
			return nil, err
		}

		modifiedTime, err := storage.ModifiedTime(i.Filepath)

		if err != nil {
			return nil, err
		}

		contentType := i.GetContentType()

		headers := map[string]string{
			"Last-Modified": modifiedTime.Format(storages.LastModifiedFormat),
			"Content-Type":  contentType,
		}

		reader := bytes.NewReader(body)

		dest, err := imaging.Decode(reader)

		if err != nil {
			return nil, err
		}

		return &ImageFile{Source: dest, Storage: storage, Header: headers, Filepath: i.Filepath}, nil
	}

	return file, err
}

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}
