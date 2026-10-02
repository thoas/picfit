package image

import (
	"strings"
	"testing"
)

type closeRecorder struct {
	*strings.Reader
	closed bool
}

func (c *closeRecorder) Close() error {
	c.closed = true
	return nil
}

func TestImageFileClose(t *testing.T) {
	t.Run("without stream", func(t *testing.T) {
		// get and redirect return cached files without loading their stream
		file := &ImageFile{HTTPStream: strings.NewReader("data")}
		file.Close()
		if file.HTTPStream != nil {
			t.Fatal("HTTPStream not reset")
		}
	})

	t.Run("with stream", func(t *testing.T) {
		stream := &closeRecorder{Reader: strings.NewReader("data")}
		file := &ImageFile{Stream: stream, HTTPStream: stream}
		file.Close()
		if !stream.closed {
			t.Fatal("stream not closed")
		}
		if file.HTTPStream != nil {
			t.Fatal("HTTPStream not reset")
		}
	})
}
