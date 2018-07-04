package backend

import (
	"bytes"
	"math"

	"github.com/discordapp/lilliput"
	"github.com/pkg/errors"

	imagefile "github.com/thoas/picfit/image"
)

type LilliputEngine struct {
	MaxBufferSize int
}

func NewLilliputEngine(maxBufferSize int) *LilliputEngine {
	if maxBufferSize > 0 {
		return &LilliputEngine{MaxBufferSize: maxBufferSize}
	}
	return &LilliputEngine{MaxBufferSize: 8192}
}

// Resize resizes the image to the specified width and height and
// returns the transformed image. If one of width or height is 0,
// the image aspect ratio is preserved.
func (e *LilliputEngine) Resize(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	opts := &lilliput.ImageOptions{
		FileType:             img.FilenameExt(),
		Width:                options.Width,
		Height:               options.Height,
		NormalizeOrientation: true,
		ResizeMethod:         lilliput.ImageOpsResize,
	}

	return e.transform(img, opts, options.Upscale)
}

func (e *LilliputEngine) Rotate(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *LilliputEngine) Flip(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Thumbnail scales the image up or down using the specified resample filter, crops it
// to the specified width and hight and returns the transformed image.
func (e *LilliputEngine) Thumbnail(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	opts := &lilliput.ImageOptions{
		FileType:             img.FilenameExt(),
		Width:                options.Width,
		Height:               options.Height,
		NormalizeOrientation: true,
		// Lilliput ImageOpsFit is a thumbnail operation
		ResizeMethod: lilliput.ImageOpsFit,
	}

	return e.transform(img, opts, options.Upscale)
}

func (e *LilliputEngine) Fit(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *LilliputEngine) transform(img *imagefile.ImageFile, options *lilliput.ImageOptions, upscale bool) ([]byte, error) {
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

	outputImg := make([]byte, 50*1024*1024)

	return ops.Transform(decoder, options, outputImg)
}
