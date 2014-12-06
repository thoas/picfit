package image

import (
	"github.com/disintegration/imaging"
)

type Method struct {
	Name           string
	Params         []string
	Transformation Transformation
}

var Resize = &Method{
	"resize",
	[]string{"w", "h"},
	imaging.Resize,
}

var Thumbnail = &Method{
	"thumbnail",
	[]string{"w", "h"},
	imaging.Thumbnail,
}

var Methods = map[string]*Method{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
}

func (m *Method) HasEnoughParams(params map[string]string) bool {
	for _, param := range m.Params {
		if _, ok := params[param]; !ok {
			return false
		}
	}

	return true
}
