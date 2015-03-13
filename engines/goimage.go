package engines

import (
	"bytes"
	"fmt"
	"github.com/thoas/imaging"
	"image"
	"image/jpeg"
	"math"
)

var Formats = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"jpg":  imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
}

type GoImageEngine struct{}

func NewGoImageEngine() Engine {
	return &GoImageEngine{}
}

type Transformation func(img image.Image, width int, height int, filter imaging.ResampleFilter) *image.NRGBA

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}

func ImageSize(i image.Image) (int, int) {
	return i.Bounds().Max.X, i.Bounds().Max.Y
}

func (i *GoImageEngine) Scale(img image.Image, dstWidth int, dstHeight int, upscale bool, trans Transformation) *image.NRGBA {
	width, height := ImageSize(img)

	factor := scalingFactor(width, height, dstWidth, dstHeight)

	if factor < 1 || upscale {
		return trans(img, dstWidth, dstHeight, imaging.Lanczos)
	}

	return imaging.Clone(img)
}

func (i *GoImageEngine) Resize(source []byte, width int, height int, options *Options) ([]byte, error) {
	image, err := i.ImageFromSource(source)

	if err != nil {
		return nil, err
	}

	return i.ToBytes(i.Scale(image, width, height, options.Upscale, imaging.Resize), options.Format, options.Quality)
}

func (i *GoImageEngine) ImageFromSource(source []byte) (image.Image, error) {
	return imaging.Decode(bytes.NewReader(source))
}

func (i *GoImageEngine) Thumbnail(source []byte, width int, height int, options *Options) ([]byte, error) {
	image, err := i.ImageFromSource(source)

	if err != nil {
		return nil, err
	}

	return i.ToBytes(i.Scale(image, width, height, options.Upscale, imaging.Thumbnail), options.Format, options.Quality)
}

func (i *GoImageEngine) ToBytes(img image.Image, format string, quality int) ([]byte, error) {
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
