package image

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"math"
	"mime"
	"strconv"
)

type ImageFile struct {
	Source      image.Image
	ContentType string
	Key         string
	Header      map[string]string
	Filename    string
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

func (i *ImageFile) Transform(method *Method, qs map[string]string) (*ImageFile, error) {
	_, ok := qs["upscale"]

	if !ok {
		qs["upscale"] = "1"
	}

	switch method {
	case Resize, Thumbnail:
		w, err := strconv.Atoi(qs["w"])

		if err != nil {
			return nil, err
		}

		h, err := strconv.Atoi(qs["h"])

		if err != nil {
			return nil, err
		}

		upscale, err := strconv.ParseBool(qs["upscale"])

		if err != nil {
			return nil, err
		}

		dest := i.Scale([]int{w, h}, upscale, method.Transformation)

		file := &ImageFile{
			Source:      dest,
			ContentType: i.ContentType,
			Key:         i.Key,
			Header:      i.Header,
			Filename:    i.Filename,
		}

		return file, err
	}

	return nil, fmt.Errorf("Method not found for %s", method)
}

func (i *ImageFile) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := imaging.Encode(buf, i.Source, Formats[i.ContentType])

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (i *ImageFile) Format() string {
	if i.ContentType != "" {
		return Extensions[i.ContentType]
	}

	if i.Filename != "" {
		return Extensions[mime.TypeByExtension(i.Filename)]
	}

	return ""
}

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}
