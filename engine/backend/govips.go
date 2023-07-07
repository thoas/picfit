package backend

import (
	"bytes"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/thoas/picfit/image"
)

type Govips struct {
}

func (g *Govips) Fit(img *image.ImageFile, options *Options) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Govips) Flat(background *image.ImageFile, options *Options) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Govips) Flip(img *image.ImageFile, options *Options) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Govips) Resize(img *image.ImageFile, options *Options) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Govips) Rotate(img *image.ImageFile, options *Options) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Govips) String() string {
	return "govips"
}

func (g *Govips) Thumbnail(img *image.ImageFile, options *Options) ([]byte, error) {
	source := bytes.NewReader(img.Source)
	i, err := vips.NewImageFromReader(source)
	if err != nil {
		return nil, err
	}
	if err := i.Thumbnail(options.Width, options.Height, 0); err != nil {
		return nil, err
	}
	return i.ToBytes()
}
