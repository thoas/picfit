package backend

import (
	"bytes"
	"math"

	"github.com/discordapp/lilliput"
	"github.com/pkg/errors"

	"github.com/thoas/picfit/engine/config"
	imagefile "github.com/thoas/picfit/image"
)

type Lilliput struct {
	MaxBufferSize   int
	ImageBufferSize int
	EncodeOptions   map[int]int
}

func NewLilliput(cfg config.Config) *Lilliput {
	maxBufferSize := config.DefaultMaxBufferSize
	if cfg.MaxBufferSize != 0 {
		maxBufferSize = cfg.MaxBufferSize
	}

	imageBufferSize := config.DefaultImageBufferSize
	if cfg.ImageBufferSize != 0 {
		imageBufferSize = cfg.ImageBufferSize
	}

	jpegQuality := config.DefaultQuality
	if cfg.JpegQuality != 0 {
		jpegQuality = cfg.JpegQuality
	}

	webpQuality := config.DefaultQuality
	if cfg.WebpQuality != 0 {
		webpQuality = cfg.WebpQuality
	}

	pngCompression := config.DefaultPngCompression
	if cfg.PngCompression != 0 {
		pngCompression = cfg.PngCompression
	}

	return &Lilliput{
		MaxBufferSize:   maxBufferSize,
		ImageBufferSize: imageBufferSize,
		EncodeOptions: map[int]int{
			lilliput.JpegQuality:    jpegQuality,
			lilliput.PngCompression: pngCompression,
			lilliput.WebpQuality:    webpQuality,
		}}
}

// Resize resizes the image to the specified width and height and
// returns the transformed image. If one of width or height is 0,
// the image aspect ratio is preserved.
func (e *Lilliput) Resize(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	opts := &lilliput.ImageOptions{
		FileType:             img.FilenameExt(),
		Width:                options.Width,
		Height:               options.Height,
		NormalizeOrientation: true,
		ResizeMethod:         lilliput.ImageOpsResize,
		EncodeOptions:        e.EncodeOptions,
	}

	return e.transform(img, opts, options.Upscale)
}

func (e *Lilliput) Rotate(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Lilliput) Flip(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Thumbnail scales the image up or down using the specified resample filter, crops it
// to the specified width and hight and returns the transformed image.
func (e *Lilliput) Thumbnail(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	opts := &lilliput.ImageOptions{
		FileType:             img.FilenameExt(),
		Width:                options.Width,
		Height:               options.Height,
		NormalizeOrientation: true,
		// Lilliput ImageOpsFit is a thumbnail operation
		ResizeMethod:  lilliput.ImageOpsFit,
		EncodeOptions: e.EncodeOptions,
	}

	return e.transform(img, opts, options.Upscale)
}

func (e *Lilliput) Fit(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Lilliput) transform(img *imagefile.ImageFile, options *lilliput.ImageOptions, upscale bool) ([]byte, error) {
	decoder, err := lilliput.NewDecoder(img.Source)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		return nil, err
	}

	var (
		srcW int
		srcH int
	)

	same, err := sameInputAndOutputHeader(bytes.NewReader(img.Source))
	if err != nil {
		return nil, err
	}

	if same {
		srcW = header.Width()
		srcH = header.Height()
	} else {
		srcH = header.Width()
		srcW = header.Height()
	}

	if scalingFactor(srcW, srcH, options.Width, options.Height) > 1 && !upscale {
		return img.Source, nil
	}

	if options.Width == 0 {
		tmpW := float64(options.Height) * float64(srcW) / float64(srcH)
		options.Width = int(math.Max(1.0, math.Floor(tmpW+0.5)))
	}
	if options.Height == 0 {
		tmpH := float64(options.Width) * float64(srcH) / float64(srcW)
		options.Height = int(math.Max(1.0, math.Floor(tmpH+0.5)))
	}

	ops := lilliput.NewImageOps(e.MaxBufferSize)
	defer ops.Close()

	outputImg := make([]byte, e.ImageBufferSize)

	return ops.Transform(decoder, options, outputImg)
}

func (e *Lilliput) String() string {
	return "lilliput"
}
