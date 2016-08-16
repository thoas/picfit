package engine

type Operation struct {
	Name string
}

var Resize = &Operation{
	"resize",
}

var Thumbnail = &Operation{
	"thumbnail",
}

var Rotate = &Operation{
	"rotate",
}

var Flip = &Operation{
	"flip",
}

var Fit = &Operation{
	"fit",
}

var Noop = &Operation{
	"noop",
}

var Operations = map[string]*Operation{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
	Flip.Name:      Flip,
	Rotate.Name:    Rotate,
	Fit.Name:       Fit,
	Noop.Name:      Noop,
}
