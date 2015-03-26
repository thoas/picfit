package engines

import (
	"bytes"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/thoas/imaging"
	imagefile "github.com/thoas/picfit/image"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"math"
	"strconv"
	"time"
)

type GoImageEngine struct {
	DefaultFormat string
}

type Result struct {
	Paletted *image.Paletted
	Image    *image.NRGBA
	Position int
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

func (e *GoImageEngine) TransformGIF(img *imagefile.ImageFile, width int, height int, options *Options, trans Transformation) ([]byte, error) {
	g, err := gif.DecodeAll(bytes.NewReader(img.Source))

	if err != nil {
		return nil, err
	}

	length := len(g.Image)
	done := make(chan *Result)
	images := make([]*image.Paletted, length)
	processed := 0

	for i := range g.Image {
		go func(paletted *image.Paletted, width int, height int, position int, options *Options) {
			done <- &Result{
				Image:    e.Scale(paletted, width, height, options.Upscale, imaging.Resize),
				Position: position,
				Paletted: image.NewPaletted(image.Rect(0, 0, width, height), paletted.Palette),
			}
		}(g.Image[i], width, height, i, options)
	}

	for {
		select {
		case result := <-done:
			draw.Draw(result.Paletted, image.Rect(0, 0, width, height), result.Image, image.Pt(0, 0), draw.Src)

			images[result.Position] = result.Paletted

			processed++
		case <-time.After(time.Second * 5):
			break
		}

		if processed == length {
			break
		}
	}

	close(done)

	g.Image = images

	buf := &bytes.Buffer{}

	err = gif.EncodeAll(buf, g)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *GoImageEngine) Resize(img *imagefile.ImageFile, width int, height int, options *Options) ([]byte, error) {
	if options.Format == imaging.GIF {
		content, err := e.TransformGIF(img, width, height, options, imaging.Resize)

		if err != nil {
			return nil, err
		}

		return content, nil
	}

	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	return e.resize(image, width, height, options)
}

func (e *GoImageEngine) resize(img image.Image, width int, height int, options *Options) ([]byte, error) {
	return e.ToBytes(e.Scale(img, width, height, options.Upscale, imaging.Resize), options.Format, options.Quality)
}

func (e *GoImageEngine) Source(img *imagefile.ImageFile) (image.Image, error) {
	return imaging.Decode(bytes.NewReader(img.Source))
}

func (e *GoImageEngine) Thumbnail(img *imagefile.ImageFile, width int, height int, options *Options) ([]byte, error) {
	if options.Format == imaging.GIF {
		content, err := e.TransformGIF(img, width, height, options, imaging.Thumbnail)

		if err != nil {
			return nil, err
		}

		return content, nil
	}

	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	return e.thumbnail(image, width, height, options)
}

func (e *GoImageEngine) thumbnail(img image.Image, width int, height int, options *Options) ([]byte, error) {
	return e.ToBytes(e.Scale(img, width, height, options.Upscale, imaging.Thumbnail), options.Format, options.Quality)
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
	filepath := img.Filepath

	if ok {
		if _, ok := ContentTypes[format]; !ok {
			return nil, fmt.Errorf("Unknown format %s", format)
		}

	} else {
		format = img.Format()
	}

	if format == "" {
		format = e.DefaultFormat
	}

	if format != img.Format() {
		index := len(filepath) - len(img.Format())

		filepath = filepath[:index] + format
	}

	file := &imagefile.ImageFile{
		Source:   img.Source,
		Key:      img.Key,
		Headers:  img.Headers,
		Filepath: filepath,
	}

	options := &Options{
		Quality: quality,
		Format:  Formats[format],
		Upscale: upscale,
	}

	switch operation {
	case Resize:
		content, err := e.Resize(img, w, h, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	case Thumbnail:
		content, err := e.Thumbnail(img, w, h, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (e *GoImageEngine) ToBytes(img image.Image, format imaging.Format, quality int) ([]byte, error) {
	buf := &bytes.Buffer{}

	var err error

	if format == imaging.JPEG && quality > 0 {
		err = imaging.EncodeWithOptions(buf, img, format, &jpeg.Options{Quality: quality})
	} else {
		err = imaging.Encode(buf, img, format)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
