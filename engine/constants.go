package engine

var ContentTypes = map[string]string{
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"bmp":  "image/bmp",
	"gif":  "image/gif",
}

const (
	TopRight    = "top-right"
	TopLeft     = "top-left"
	BottomRight = "bottom-right"
	BottomLeft  = "bottom-left"
)

var StickPositions = []string{
	TopRight,
	TopLeft,
	BottomRight,
	BottomLeft,
}
