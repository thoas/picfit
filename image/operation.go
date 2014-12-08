package image

import (
	"github.com/disintegration/imaging"
)

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
