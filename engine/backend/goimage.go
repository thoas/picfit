package backend

import (
	"context"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"

	"github.com/thoas/picfit/constants"

	"github.com/go-spectest/imaging"

	imagefile "github.com/thoas/picfit/image"

	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type (
	imageTransformation func(img image.Image) *image.NRGBA
	transformation      func(img image.Image, width int, height int, filter imaging.ResampleFilter) *image.NRGBA
)

var (
	flipTransformations = map[string]imageTransformation{
		"h": imaging.FlipH,
		"v": imaging.FlipV,
	}

	rotateTransformations = map[int]imageTransformation{
		90:  imaging.Rotate90,
		270: imaging.Rotate270,
		180: imaging.Rotate180,
	}
)

type GoImage struct{}

func (e *GoImage) String() string {
	return "goimage"
}
func (e *GoImage) Resize(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	return e.resize(dst, img, options, imaging.Resize)
}

func (e *GoImage) Thumbnail(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	return e.resize(dst, img, options, imaging.Thumbnail)
}

func (e *GoImage) Rotate(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	image, err := e.source(img)
	if err != nil {
		return err
	}

	deg := options.Degree

	transform, ok := rotateTransformations[deg]
	if !ok {
		return fmt.Errorf("Invalid rotate transformation degree=%d is not supported", deg)
	}

	return encode(dst, transform(image), options.Format, options.Quality)
}

func (e *GoImage) Flip(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	image, err := e.source(img)
	if err != nil {
		return err
	}

	pos := options.Position

	transform, ok := flipTransformations[pos]
	if !ok {
		return fmt.Errorf("Invalid flip transformation, %s is not supported", pos)
	}

	return encode(dst, transform(image), options.Format, options.Quality)
}

func (e *GoImage) Fit(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	if options.Format == imagefile.GIF {
		err := e.transformGIF(dst, img, options, imaging.Thumbnail)
		if err != nil {
			return err
		}

		return nil
	}

	image, err := e.source(img)
	if err != nil {
		return err
	}

	return e.transform(dst, image, options, imaging.Fit)
}

func (e *GoImage) Effect(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	image, err := e.source(img)
	if err != nil {
		return err
	}
	width, height := imageSize(image)
	size := max(width, height)
	sigma := size / 20

	const maxSigma = 50
	if sigma > maxSigma {
		sigma = maxSigma
	}
	switch options.Filter {
	case constants.FilterBlur:
		return encode(dst, imaging.Blur(image, float64(sigma)), options.Format, options.Quality)
	}

	return MethodNotImplementedError
}

func (e *GoImage) transformGIF(dst io.Writer, img *imagefile.ImageFile, options *Options, trans transformation) error {
	g, err := gif.DecodeAll(img.Stream)
	if err != nil {
		return err
	}
	if len(g.Image) == 0 {
		return fmt.Errorf("GIF has no frames")
	}
	if len(g.Image) == 0 {
		return fmt.Errorf("GIF has no frames")
	}

	first := g.Image[0]
	factor := scalingFactorImage(first, options.Width, options.Height)
	if factor > 1 && !options.Upscale {
		return gif.EncodeAll(dst, g)
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
		options.Width = int(max(1.0, math.Floor(tmpW+0.5)))
	}
	if options.Height == 0 {
		tmpH := float64(options.Width) * float64(srcH) / float64(srcW)
		options.Height = int(max(1.0, math.Floor(tmpH+0.5)))
	}

	g.Config.Height = options.Height
	g.Config.Width = options.Width

	if err := gif.EncodeAll(dst, g); err != nil {
		return err
	}

	return nil
}

func (e *GoImage) resize(dst io.Writer, img *imagefile.ImageFile, options *Options, trans transformation) error {
	if options.Format == imagefile.GIF {
		err := e.transformGIF(dst, img, options, trans)
		if err != nil {
			return err
		}

		return nil
	}

	image, err := e.source(img)
	if err != nil {
		return err
	}

	return e.transform(dst, image, options, trans)
}

func (e *GoImage) transform(dst io.Writer, img image.Image, options *Options, trans transformation) error {
	if options.Height == 0 && options.Width == 0 {
		return encode(dst, img, options.Format, options.Quality)
	}

	return encode(dst, scale(img, options, trans), options.Format, options.Quality)
}

func (e *GoImage) source(img *imagefile.ImageFile) (image.Image, error) {
	return decode(img.Stream)
}

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}

func scalingFactorImage(img image.Image, dstWidth int, dstHeight int) float64 {
	width, height := imageSize(img)

	return scalingFactor(width, height, dstWidth, dstHeight)
}

func imageSize(e image.Image) (int, int) {
	return e.Bounds().Max.X, e.Bounds().Max.Y
}

func scale(img image.Image, options *Options, trans transformation) image.Image {
	factor := scalingFactorImage(img, options.Width, options.Height)

	if factor < 1 || options.Upscale {
		return trans(img, options.Width, options.Height, imaging.Lanczos)
	}

	return img
}

func imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})
	return pm
}

func encode(w io.Writer, img image.Image, format imagefile.Format, quality int) error {
	var err error
	switch format {
	case imagefile.JPEG:
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

	case imagefile.PNG:
		err = png.Encode(w, img)
	case imagefile.GIF:
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	case imagefile.TIFF:
		err = tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case imagefile.BMP:
		err = bmp.Encode(w, img)
	case imagefile.WEBP:
		err = webp.Encode(w, img, &webp.Options{Quality: float32(quality)})
	default:
		err = imaging.ErrUnsupportedFormat
	}
	return err
}
