package image

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

//Decode is image.Decode handling orientation in EXIF tags if exists.
//Requires io.ReadSeeker instead of io.Reader.
func Decode(reader io.ReadSeeker) (image.Image, error) {
	img, err := imaging.Decode(reader)
	if err != nil {
		return img, err
	}
	reader.Seek(0, io.SeekStart)
	orientation := getOrientation(reader)
	switch orientation {
	case "1":
	case "2":
		img = imaging.FlipV(img)
	case "3":
		img = imaging.Rotate180(img)
	case "4":
		img = imaging.Rotate180(imaging.FlipV(img))
	case "5":
		img = imaging.Rotate270(imaging.FlipV(img))
	case "6":
		img = imaging.Rotate270(img)
	case "7":
		img = imaging.Rotate90(imaging.FlipV(img))
	case "8":
		img = imaging.Rotate90(img)
	}

	return img, err
}

func getOrientation(reader io.Reader) string {
	x, err := exif.Decode(reader)
	if err != nil {
		return "1"
	}
	if x != nil {
		orient, err := x.Get(exif.Orientation)
		if err != nil {
			return "1"
		}
		if orient != nil {
			return orient.String()
		}
	}

	return "1"
}
