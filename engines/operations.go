package engines

type Operation struct {
	Name string
}

var Resize = &Operation{
	"resize",
}

var Thumbnail = &Operation{
	"thumbnail",
}

var Operations = map[string]*Operation{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
}
