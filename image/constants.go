package image

import (
	"github.com/disintegration/imaging"
)

var Formats = map[string]imaging.Format{
	"image/jpeg": imaging.JPEG,
	"image/png":  imaging.PNG,
	"image/gif":  imaging.GIF,
	"image/bmp":  imaging.BMP,
}

var ContentTypes = map[string]string{
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"bmp":  "image/bmp",
	"gif":  "image/gif",
}

var Extensions = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/bmp":  "bmp",
	"image/gif":  "gif",
}
