package backend

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	_ "golang.org/x/image/webp"
)

//Decode is image.Decode handling orientation in EXIF tags if exists.
//Requires io.ReadSeeker instead of io.Reader.
func decode(reader io.ReadSeeker) (image.Image, error) {
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, io.SeekStart)
	orientation := getOrientation(reader)
	switch orientation {
	case "1":
		return img, nil
	case "2":
		return imaging.FlipV(img), nil
	case "3":
		return imaging.Rotate180(img), nil
	case "4":
		return imaging.Rotate180(imaging.FlipV(img)), nil
	case "5":
		return imaging.Rotate270(imaging.FlipV(img)), nil
	case "6":
		return imaging.Rotate270(img), nil
	case "7":
		return imaging.Rotate90(imaging.FlipV(img)), nil
	case "8":
		return imaging.Rotate90(img), nil
	default:
		return img, nil
	}
}

// SameInputAndOutputHeader return true if image width and height
// are not changed after exif correction.
func sameInputAndOutputHeader(reader io.ReadSeeker) (bool, error) {
	_, err := imaging.Decode(reader)
	if err != nil {
		return false, err
	}
	reader.Seek(0, io.SeekStart)
	orientation := getOrientation(reader)
	switch orientation {
	case "1":
		return true, nil
	case "2":
		return true, nil
	case "3":
		return true, nil
	case "4":
		return true, nil
	case "5":
		return false, nil
	case "6":
		return false, nil
	case "7":
		return false, nil
	case "8":
		return false, nil
	default:
		return true, nil
	}
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
