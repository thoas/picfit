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
	httppkg "github.com/thoas/picfit/http"
)

const (
	MimetypeDetectorTypeFastimage = "fastimage"
	MimetypeDetectorTypeSniff     = "sniff"
)

func GetMimetypeDetector(mimetypeDetectorType string) MimetypeDetectorFunc {
	return GetMimetypeDetectorWithClient(mimetypeDetectorType, nil)
}

// GetMimetypeDetectorWithClient is GetMimetypeDetector with the HTTP client
// used by detectors that fetch the URL. A nil client keeps each detector's
// historical client. Pass the fetch-policy client so those requests are held
// to the same dial checks as the image download.
func GetMimetypeDetectorWithClient(mimetypeDetectorType string, client *httppkg.Client) MimetypeDetectorFunc {
	switch mimetypeDetectorType {
	case MimetypeDetectorTypeFastimage:
		if client == nil {
			return MimetypeDetectorFastimage
		}
		return func(uri *url.URL) (string, error) {
			return fastimageWithClient(uri, client)
		}
	case MimetypeDetectorTypeSniff:
		if client == nil {
			return MimetypeDetectorSniff
		}
		return func(uri *url.URL) (string, error) {
			return sniffWithClient(uri, client)
		}
	default:
		return MimetypeDetectorExtension
	}
}

type MimetypeDetectorFunc func(*url.URL) (string, error)

// Detect mimetype by looking at the URL extension
func MimetypeDetectorExtension(uri *url.URL) (string, error) {
	return mime.TypeByExtension(filepath.Ext(uri.Path)), nil
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

func fastimageWithClient(uri *url.URL, client *httppkg.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	imageType, _, err := fastimage.DetectImageTypeFromResponse(resp)
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

	buffer := make([]byte, 512)

	n, err := resp.Body.Read(buffer)

	// n.b. will ignore EOF errors if buffer < 512
	if err != nil && err != io.EOF {
		return "", err
	}

	if err := resp.Body.Close(); err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer[:n])

	return contentType, nil
}

func sniffWithClient(uri *url.URL, client *httppkg.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buffer := make([]byte, 512)
	n, err := resp.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	return http.DetectContentType(buffer[:n]), nil
}
