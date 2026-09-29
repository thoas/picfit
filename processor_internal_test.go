package picfit

import (
	"bytes"
	imagepkg "image"
	"image/png"
	"io"
	"testing"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/image"
)

func TestCheckImageMaxDimension(t *testing.T) {
	tests := []struct {
		name          string
		max           config.AllowedSize
		width, height int
		wantErr       bool
	}{
		{"within limits", config.AllowedSize{Width: 10, Height: 10}, 10, 10, false},
		{"width exceeded", config.AllowedSize{Width: 10, Height: 10}, 20, 5, true},
		{"height exceeded", config.AllowedSize{Width: 10, Height: 10}, 5, 20, true},
		{"both exceeded", config.AllowedSize{Width: 10, Height: 10}, 20, 20, true},
		{"no height limit", config.AllowedSize{Width: 10}, 5, 20, false},
		{"no width limit", config.AllowedSize{Height: 10}, 20, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Processor{maxImageDimensions: &tt.max}
			err := p.checkImageMaxDimension(imagepkg.Config{Width: tt.width, Height: tt.height})
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkImageMaxDimension(%dx%d) error = %v, wantErr %v", tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

func TestDecodeConfigKeepsStream(t *testing.T) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, imagepkg.NewRGBA(imagepkg.Rect(0, 0, 30, 20))); err != nil {
		t.Fatal(err)
	}
	original := buf.Bytes()
	file := &image.ImageFile{Stream: io.NopCloser(bytes.NewReader(original))}

	imageconfig, format, err := decodeConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	if imageconfig.Width != 30 || imageconfig.Height != 20 || format != "png" {
		t.Fatalf("got %dx%d %s, want 30x20 png", imageconfig.Width, imageconfig.Height, format)
	}

	data, err := io.ReadAll(file.Stream)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, original) {
		t.Fatalf("stream altered: got %d bytes, want %d", len(data), len(original))
	}
	if err := file.Stream.Close(); err != nil {
		t.Fatal(err)
	}
}
