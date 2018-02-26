package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/imdario/mergo"

	imagefile "github.com/thoas/picfit/image"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type GoImageEngine struct {
	DefaultFormat  string
	Format         string
	DefaultQuality int
}

type ImageTransformation func(img image.Image) *image.NRGBA

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

var flipTransformations = map[string]ImageTransformation{
	"h": imaging.FlipH,
	"v": imaging.FlipV,
}

var rotateTransformations = map[int]ImageTransformation{
	90:  imaging.Rotate90,
	270: imaging.Rotate270,
	180: imaging.Rotate180,
}

type Transformation func(img image.Image, width int, height int, filter imaging.ResampleFilter) *image.NRGBA

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}

func scalingFactorImage(img image.Image, dstWidth int, dstHeight int) float64 {
	width, height := imageSize(img)

	return scalingFactor(width, height, dstWidth, dstHeight)
}

func imageSize(e image.Image) (int, int) {
	return e.Bounds().Max.X, e.Bounds().Max.Y
}

func scale(img image.Image, options *Options, trans Transformation) image.Image {
	factor := scalingFactorImage(img, options.Width, options.Height)

	if factor < 1 || options.Upscale {
		return trans(img, options.Width, options.Height, imaging.Lanczos)
	}

	return img
}

func imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}

func (e *GoImageEngine) TransformGIF(img *imagefile.ImageFile, options *Options, trans Transformation) ([]byte, error) {
	first, err := gif.Decode(bytes.NewReader(img.Source))

	if err != nil {
		return nil, err
	}

	factor := scalingFactorImage(first, options.Width, options.Height)

	if factor > 1 && !options.Upscale {
		return img.Source, nil
	}

	g, err := gif.DecodeAll(bytes.NewReader(img.Source))

	if err != nil {
		return nil, err
	}

	firstFrame := g.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
	im := image.NewRGBA(b)

	for i, frame := range g.Image {
		bounds := frame.Bounds()
		draw.Draw(im, bounds, frame, bounds.Min, draw.Over)
		g.Image[i] = imageToPaletted(scale(im, options, trans))
	}

	srcW, srcH := imageSize(first)

	if options.Width == 0 {
		tmpW := float64(options.Height) * float64(srcW) / float64(srcH)
		options.Width = int(math.Max(1.0, math.Floor(tmpW+0.5)))
	}
	if options.Height == 0 {
		tmpH := float64(options.Width) * float64(srcH) / float64(srcW)
		options.Height = int(math.Max(1.0, math.Floor(tmpH+0.5)))
	}

	g.Config.Height = options.Height
	g.Config.Width = options.Width

	buf := &bytes.Buffer{}

	err = gif.EncodeAll(buf, g)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *GoImageEngine) Resize(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	if options.Format == imaging.GIF {
		content, err := e.TransformGIF(img, options, imaging.Resize)

		if err != nil {
			return nil, err
		}

		return content, nil
	}

	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	return e.transform(image, options, imaging.Resize)
}

func (e *GoImageEngine) transform(img image.Image, options *Options, trans Transformation) ([]byte, error) {
	return e.ToBytes(scale(img, options, trans), options.Format, options.Quality)
}

func (e *GoImageEngine) Source(img *imagefile.ImageFile) (image.Image, error) {
	return imaging.Decode(bytes.NewReader(img.Source))
}

func (e *GoImageEngine) Rotate(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	deg := options.Degree

	transform, ok := rotateTransformations[deg]

	if !ok {
		return nil, fmt.Errorf("Invalid rotate transformation degree=%d is not supported", deg)
	}

	return e.ToBytes(transform(image), options.Format, options.Quality)
}

func (e *GoImageEngine) Flip(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	pos := options.Position

	transform, ok := flipTransformations[pos]

	if !ok {
		return nil, fmt.Errorf("Invalid flip transformation, %s is not supported", pos)
	}

	return e.ToBytes(transform(image), options.Format, options.Quality)
}

func (e *GoImageEngine) Thumbnail(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	if options.Format == imaging.GIF {
		content, err := e.TransformGIF(img, options, imaging.Thumbnail)

		if err != nil {
			return nil, err
		}

		return content, nil
	}

	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	return e.transform(image, options, imaging.Thumbnail)
}

func (e *GoImageEngine) Fit(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	if options.Format == imaging.GIF {
		content, err := e.TransformGIF(img, options, imaging.Thumbnail)

		if err != nil {
			return nil, err
		}

		return content, nil
	}

	image, err := e.Source(img)

	if err != nil {
		return nil, err
	}

	return e.transform(image, options, imaging.Fit)
}

func (e *GoImageEngine) Transform(img *imagefile.ImageFile, operation *Operation, qs map[string]string) (*imagefile.ImageFile, error) {
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

	file := &imagefile.ImageFile{
		Source:   img.Source,
		Key:      img.Key,
		Headers:  img.Headers,
		Filepath: filepath,
	}

	options := &Options{
		Quality: quality,
		Format:  formats[format],
	}

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

		content, err := e.Flip(img, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	case Rotate:
		deg, err := strconv.Atoi(qs["deg"])

		if err != nil {
			return nil, err
		}

		options.Degree = deg

		content, err := e.Rotate(img, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
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
			content, err := e.Resize(img, options)

			if err != nil {
				return nil, err
			}

			file.Processed = content

			return file, err
		case Thumbnail:
			content, err := e.Thumbnail(img, options)

			if err != nil {
				return nil, err
			}

			file.Processed = content

			return file, err
		case Fit:
			content, err := e.Fit(img, options)

			if err != nil {
				return nil, err
			}

			file.Processed = content

			return file, err
		}
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (e *GoImageEngine) ToBytes(img image.Image, format imaging.Format, quality int) ([]byte, error) {
	buf := &bytes.Buffer{}

	var err error

	err = encode(buf, img, format, quality)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func encode(w io.Writer, img image.Image, format imaging.Format, quality int) error {
	var err error
	switch format {
	case imaging.JPEG:
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
		}

	case imaging.PNG:
		err = png.Encode(w, img)
	case imaging.GIF:
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	case imaging.TIFF:
		err = tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case imaging.BMP:
		err = bmp.Encode(w, img)
	default:
		err = imaging.ErrUnsupportedFormat
	}
	return err
}
