package engine

import "github.com/thoas/picfit/engine/backend"

type Operation string

func (o Operation) String() string {
	return string(o)
}

const (
	Fit       = Operation("fit")
	Flat      = Operation("flat")
	Flip      = Operation("flip")
	Noop      = Operation("noop")
	Resize    = Operation("resize")
	Rotate    = Operation("rotate")
	Thumbnail = Operation("thumbnail")
)

var Operations = map[string]Operation{
	Fit.String():       Fit,
	Flat.String():      Flat,
	Flip.String():      Flip,
	Noop.String():      Noop,
	Resize.String():    Resize,
	Rotate.String():    Rotate,
	Thumbnail.String(): Thumbnail,
}

type EngineOperation struct {
	Operation Operation
	Options   *backend.Options
}
