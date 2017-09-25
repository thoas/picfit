package fastimage

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = uint(time.Duration(5000) * time.Millisecond)

// DetectImageType is the main function used to detect the type and size
// of a remote image represented by the url.
//
// Only check ImageType and ImageSize if error is not nil.
//
// If you want to set request timeout for uri, use DetectImageTypeWithTimeout() instead.
func DetectImageType(uri string) (ImageType, *ImageSize, error) {
	return DetectImageTypeWithTimeout(uri, defaultTimeout)
}

// DetectImageTypeWithTimeout acts same as DetectImageType(),
// except this function takes additional parameter timeout (millisecond)
// for setting request timeout.
//
// If timeout < 1, default timeout 5000 is used.
func DetectImageTypeWithTimeout(uri string, timeout uint) (ImageType, *ImageSize, error) {
	logger.Printf("Opening HTTP stream")
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, err := client.Get(uri)
	if err != nil {
		return Unknown, nil, err
	}
	defer closeHTTPStream(resp)

	return DetectImageTypeFromResponse(resp)
}

// DetectImageTypeFromResponse is a secondary function used to detect the type
// and size of a remote image represented by the resp.
// This way you can create your own request and then pass it here.
// Check examples from http://golang.org/pkg/net/http/
//
// Only check ImageType and ImageSize if error is not nil.
func DetectImageTypeFromResponse(resp *http.Response) (ImageType, *ImageSize, error) {
	logger.Printf("Response content-length: %v bytes", resp.ContentLength)
	return DetectImageTypeFromReader(resp.Body)
}

// DetectImageTypeFromReader detects the type and size from a stream of bytes.
//
// Only check ImageType and ImageSize if error is not nil.
func DetectImageTypeFromReader(r io.Reader) (ImageType, *ImageSize, error) {
	buffer := bytes.Buffer{}

	logger.Printf("Starting operation")
	defer logger.Printf("Ended after reading %v bytes", buffer.Len())

	for {
		err := readToBuffer(r, &buffer)
		if err != nil {
			logger.Printf("Bailing out because of err %v", err)
			return Unknown, nil, err
		}

		if buffer.Len() < 2 {
			continue
		}
		if buffer.Len() > 20000 {
			logger.Println("Type not found until buffer read exceeds 20000 bytes")
			return Unknown, nil, err
		}

		for _, ImageTypeParser := range imageTypeParsers {
			if ImageTypeParser.Detect(buffer.Bytes()) {
				t := ImageTypeParser.Type()
				size, err := ImageTypeParser.GetSize(buffer.Bytes())

				if err == nil {
					logger.Printf("Found image type %v with size %v", t, size)
					return t, size, nil
				}
				break
			}
		}
	}
}

func readToBuffer(body io.Reader, buffer *bytes.Buffer) error {
	chunk := make([]byte, 8)
	count, err := body.Read(chunk)

	logger.Printf("Read %v bytes", count)
	buffer.Write(chunk[:count])

	return err
}

func closeHTTPStream(http *http.Response) {
	logger.Printf("Closing HTTP Stream")
	http.Body.Close()
}
