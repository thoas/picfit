package engines

import (
	"bytes"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/thoas/imaging"
	imagefile "github.com/thoas/picfit/image"
	"image"
	"image/jpeg"
	"math"
	"strconv"
)

var Formats = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"jpg":  imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
}

type GoImageEngine struct {
	DefaultFormat string
}

func NewGoImageEngine(DefaultFormat string) Engine {
	return &GoImageEngine{
		DefaultFormat: DefaultFormat,
	}
}

type Transformation func(img image.Image, width int, height int, filter imaging.ResampleFilter) *image.NRGBA

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}

func ImageSize(e image.Image) (int, int) {
	return e.Bounds().Max.X, e.Bounds().Max.Y
}

func (e *GoImageEngine) Scale(img image.Image, dstWidth int, dstHeight int, upscale bool, trans Transformation) *image.NRGBA {
	width, height := ImageSize(img)

	factor := scalingFactor(width, height, dstWidth, dstHeight)

	if factor < 1 || upscale {
		return trans(img, dstWidth, dstHeight, imaging.Lanczos)
	}

	return imaging.Clone(img)
}

func (e *GoImageEngine) Resize(source []byte, width int, height int, options *Options) ([]byte, error) {
	image, err := e.ImageFromSource(source)

	if err != nil {
		return nil, err
	}

	return e.ToBytes(e.Scale(image, width, height, options.Upscale, imaging.Resize), options.Format, options.Quality)
}

func (e *GoImageEngine) ImageFromSource(source []byte) (image.Image, error) {
	return imaging.Decode(bytes.NewReader(source))
}

func (e *GoImageEngine) Thumbnail(source []byte, width int, height int, options *Options) ([]byte, error) {
	image, err := e.ImageFromSource(source)

	if err != nil {
		return nil, err
	}

	return e.ToBytes(e.Scale(image, width, height, options.Upscale, imaging.Thumbnail), options.Format, options.Quality)
}

func (e *GoImageEngine) Transform(img *imagefile.ImageFile, operation *Operation, qs map[string]string) (*imagefile.ImageFile, error) {
	params := map[string]string{
		"upscale": "1",
		"h":       "0",
		"w":       "0",
	}

	err := mergo.Merge(&qs, params)

	if err != nil {
		return nil, err
	}

	upscale, err := strconv.ParseBool(qs["upscale"])

	if err != nil {
		return nil, err
	}

	w, err := strconv.Atoi(qs["w"])

	if err != nil {
		return nil, err
	}

	h, err := strconv.Atoi(qs["h"])

	if err != nil {
		return nil, err
	}

	q, ok := qs["q"]

	var quality int
	var format string

	if ok {
		quality, err := strconv.Atoi(q)

		if err != nil {
			return nil, err
		}

		if quality > 100 {
			return nil, fmt.Errorf("Quality should be <= 100")
		}
	}

	format, ok = qs["fmt"]

	if ok {
		if _, ok := imagefile.ContentTypes[format]; !ok {
			return nil, fmt.Errorf("Unknown format %s", format)
		}
	} else {
		format = e.DefaultFormat
	}

	file := &imagefile.ImageFile{
		Source:   img.Source,
		Key:      img.Key,
		Headers:  img.Headers,
		Filepath: img.Filepath,
	}

	options := &Options{
		Quality: quality,
		Format:  format,
		Upscale: upscale,
	}

	switch operation {
	case Resize:
		content, err := e.Resize(img.Source, w, h, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	case Thumbnail:
		content, err := e.Thumbnail(img.Source, w, h, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (e *GoImageEngine) ToBytes(img image.Image, format string, quality int) ([]byte, error) {
	buf := &bytes.Buffer{}

	var err error

	f := Formats[format]

	fmt.Println(format)

	if f == imaging.JPEG && quality > 0 {
		err = imaging.EncodeWithOptions(buf, img, f, &jpeg.Options{Quality: quality})
	} else {
		err = imaging.Encode(buf, img, f)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
