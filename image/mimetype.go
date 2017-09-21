package image

import (
	"mime"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/rubenfonseca/fastimage"
	"github.com/thoas/picfit/config"
)

const MimetypeDetectorTypeFastimage = "fastimage"

func GetMimetypeDetector(cfg *config.Options) MimetypeDetectorFunc {
	switch cfg.MimetypeDetector {
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
	if imageType, _, err := fastimage.DetectImageType(uri.String()); err != nil {
		return "", err
	} else {
		mimetype := "image/" + strings.ToLower(imageType.String())
		return mimetype, nil
	}
}
