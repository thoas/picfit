package backend

import (
	"github.com/discordapp/lilliput"
	"github.com/pkg/errors"

	imagefile "github.com/thoas/picfit/image"
)

type LilliputEngine struct {
	MaxBufferSize int
}

func (e *LilliputEngine) Resize(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	opts := &lilliput.ImageOptions{
		Width:        options.Width,
		Height:       options.Height,
		ResizeMethod: lilliput.ImageOpsResize,
	}

	return e.transform(img, opts)
}

func (e *LilliputEngine) Rotate(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, nil
}

func (e *LilliputEngine) Flip(img *imagefile.ImageFile, options *Options) ([]byte, error) {

	return nil, nil
}

func (e *LilliputEngine) Thumbnail(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, nil
}

func (e *LilliputEngine) Fit(img *imagefile.ImageFile, options *Options) ([]byte, error) {
	opts := &lilliput.ImageOptions{
		Width:        options.Width,
		Height:       options.Height,
		ResizeMethod: lilliput.ImageOpsNoResize,
	}

	return e.transform(img, opts)
}

func (e *LilliputEngine) transform(img *imagefile.ImageFile, options *lilliput.ImageOptions) ([]byte, error) {
	decoder, err := lilliput.NewDecoder(img.Source)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer decoder.Close()

	ops := lilliput.NewImageOps(e.MaxBufferSize)
	defer ops.Close()

	outputImg := make([]byte, 50*1024*1024)

	return ops.Transform(decoder, options, outputImg)
}
