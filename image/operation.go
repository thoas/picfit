package image

import (
	"github.com/disintegration/imaging"
)

type Operation struct {
	Name           string
	Params         []string
	Transformation Transformation
}

var Resize = &Operation{
	"resize",
	[]string{"w", "h"},
	imaging.Resize,
}

var Thumbnail = &Operation{
	"thumbnail",
	[]string{"w", "h"},
	imaging.Thumbnail,
}

var Operations = map[string]*Operation{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
}

func (m *Operation) HasEnoughParams(params map[string]string) bool {
	for _, param := range m.Params {
		if _, ok := params[param]; !ok {
			return false
		}
	}

	return true
}
