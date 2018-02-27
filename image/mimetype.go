package image

import (
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/rubenfonseca/fastimage"
)

const MimetypeDetectorTypeFastimage = "fastimage"

func GetMimetypeDetector(mimetypeDetectorType string) MimetypeDetectorFunc {
	switch mimetypeDetectorType {
	case MimetypeDetectorTypeFastimage:
		return MimetypeDetectorFastimage
	default:
		return MimetypeDetectorExtension
	}
}

type MimetypeDetectorFunc func(*url.URL) (string, error)

// Detect mimetype by looking at the URL extension
func MimetypeDetectorExtension(uri *url.URL) (string, error) {
	return mime.TypeByExtension(filepath.Ext(uri.String())), nil
}

// Detect mimetype with third-party fastimage library.
// Overhead warning: fastimage makes a request (albeit a very small partial one) upon each detection.
func MimetypeDetectorFastimage(uri *url.URL) (string, error) {
	imageType, _, err := fastimage.DetectImageType(uri.String())
	if err != nil {
		return "", err
	}

	return fmt.Sprint("image/", strings.ToLower(imageType.String())), nil
}
