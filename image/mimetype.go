package image

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/rubenfonseca/fastimage"
)

const (
	MimetypeDetectorTypeFastimage = "fastimage"
	MimetypeDetectorTypeSniff     = "sniff"
)

func GetMimetypeDetector(mimetypeDetectorType string) MimetypeDetectorFunc {
	switch mimetypeDetectorType {
	case MimetypeDetectorTypeFastimage:
		return MimetypeDetectorFastimage
	case MimetypeDetectorTypeSniff:
		return MimetypeDetectorSniff
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

// Detect mimetype using go native content detection which implements
// the MIME Sniffing alorithm described here: https://mimesniff.spec.whatwg.org/
func MimetypeDetectorSniff(uri *url.URL) (string, error) {

	client := &http.Client{
		Timeout: 5000 * time.Millisecond,
	}

	resp, err := client.Get(uri.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buffer := make([]byte, 512)

	n, err := resp.Body.Read(buffer)

	// n.b. will ignore EOF errors if buffer < 512
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])

	return contentType, nil
}
