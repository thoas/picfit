package backend

import (
	"github.com/h2non/bimg"
	"github.com/thoas/picfit/image"
)

type Libvips struct {
}

func (b *Libvips) String() string {
	return "libvips"
}

func (b *Libvips) Resize(img *image.ImageFile, options *Options) ([]byte, error) {
	res, err := bimg.NewImage(img.Source).Process(bimg.Options{
		Width:   options.Width,
		Height:  options.Height,
		Quality: options.Quality,
		Embed:   true,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *Libvips) Thumbnail(img *image.ImageFile, options *Options) ([]byte, error) {
	res, err := bimg.NewImage(img.Source).Process(bimg.Options{
		Width:   options.Width,
		Height:  options.Height,
		Quality: options.Quality,
		Crop:    true,
		Embed:   true,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *Libvips) Flip(img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (b *Libvips) Rotate(img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (b *Libvips) Fit(img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (b *Libvips) Flat(background *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

var _ Backend = (*Libvips)(nil)
