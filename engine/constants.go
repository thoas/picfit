package engine

import (
	"github.com/disintegration/imaging"
)

var ContentTypes = map[string]string{
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"bmp":  "image/bmp",
	"gif":  "image/gif",
}

var Formats = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"jpg":  imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
}
