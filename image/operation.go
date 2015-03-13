package image

import (
	"github.com/thoas/imaging"
	"image"
)

type Transformation func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA

type Operation struct {
	Name           string
	Transformation Transformation
}

var Resize = &Operation{
	"resize",
	imaging.Resize,
}

var Thumbnail = &Operation{
	"thumbnail",
	imaging.Thumbnail,
}

var Operations = map[string]*Operation{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
}
