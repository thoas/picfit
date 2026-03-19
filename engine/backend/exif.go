package backend

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/go-spectest/imaging"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	_ "golang.org/x/image/webp"
)

// Decode is image.Decode handling orientation in EXIF tags if exists.
// Requires io.ReadSeeker instead of io.Reader.
func decode(reader io.ReadCloser) (image.Image, error) {
	header := make([]byte, 65536)
	n, err := io.ReadFull(reader, header)
	defer reader.Close()

	if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, errors.WithStack(err)
	}

	orientation := getOrientation(bytes.NewReader(header[:n]))

	// rebuild full stream
	fullStream := io.MultiReader(bytes.NewReader(header[:n]), reader)

	img, _, err := image.Decode(fullStream)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch orientation {
	case "1":
		return img, nil
	case "2":
		result := imaging.FlipH(img)
		img = nil
		return result, nil
	case "3":
		result := imaging.Rotate180(img)
		img = nil
		return result, nil
	case "4":
		result := imaging.Rotate180(imaging.FlipH(img))
		img = nil
		return result, nil
	case "5":
		result := imaging.Rotate270(imaging.FlipV(img))
		img = nil
		return result, nil
	case "6":
		result := imaging.Rotate270(img)
		img = nil
		return result, nil
	case "7":
		result := imaging.Rotate90(imaging.FlipV(img))
		img = nil
		return result, nil
	case "8":
		result := imaging.Rotate90(img)
		img = nil
		return result, nil
	default:
		return img, nil
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
