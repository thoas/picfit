package engine

import (
	"bytes"
	"context"
	"encoding/json"
	stdimage "image"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/engine/config"
	"github.com/thoas/picfit/image"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "tests", "fixtures", name))
	require.NoError(t, err)
	return data
}

// handledBy returns the backend of each "Engine handled image" log line.
func handledBy(t *testing.T, logs *bytes.Buffer) []string {
	t.Helper()
	var backends []string
	dec := json.NewDecoder(logs)
	for dec.More() {
		var line map[string]any
		require.NoError(t, dec.Decode(&line))
		if line["msg"] == "Engine handled image" {
			backends = append(backends, line["backend"].(string))
		}
	}
	return backends
}

func flatOperations(t *testing.T) []EngineOperation {
	return []EngineOperation{
		{
			Operation: Resize,
			Options:   &backend.Options{Width: 200, Height: 200, Upscale: true, Format: image.PNG, Quality: 95},
		},
		{
			Operation: Flat,
			Options: &backend.Options{
				Position: "60.10.180.130",
				Format:   image.PNG,
				Quality:  95,
				Images: []image.ImageFile{
					{Stream: io.NopCloser(bytes.NewReader(readFixture(t, "avatar.png")))},
				},
			},
		},
	}
}

// sepiaOperation is supported by no backend.
var sepiaOperation = EngineOperation{
	Operation: Effect,
	Options:   &backend.Options{Filter: "sepia", Format: image.PNG},
}

func transform(t *testing.T, cfg config.Config, source, contentType string, operations []EngineOperation) ([]byte, []string, error) {
	t.Helper()
	logs := &bytes.Buffer{}
	e := New(cfg, slog.New(slog.NewJSONHandler(logs, nil)))

	output := &image.ImageFile{
		Stream:  io.NopCloser(bytes.NewReader(readFixture(t, source))),
		Headers: map[string]string{"Content-Type": contentType},
	}
	dst := &bytes.Buffer{}
	_, err := e.Transform(context.Background(), dst, output, operations)
	return dst.Bytes(), handledBy(t, logs), err
}

func decodeConfig(t *testing.T, data []byte) stdimage.Config {
	t.Helper()
	cfg, _, err := stdimage.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	return cfg
}

var vipsFirst = config.Config{Backends: &config.Backends{
	Vips:    &config.VipsBackend{Mimetypes: MimeTypes},
	GoImage: &config.Backend{Weight: 1, Mimetypes: MimeTypes},
}}

func TestTransformFlat(t *testing.T) {
	actual, backends, err := transform(t, vipsFirst, "schwarzy.jpg", "image/png", flatOperations(t))
	require.NoError(t, err)
	assert.Equal(t, []string{"vips", "vips"}, backends)

	expected, backends, err := transform(t, config.Config{}, "schwarzy.jpg", "image/png", flatOperations(t))
	require.NoError(t, err)
	assert.Equal(t, []string{"goimage", "goimage"}, backends)

	assert.Equal(t, decodeConfig(t, expected), decodeConfig(t, actual))
}

func TestTransformFallsBackOnGoImage(t *testing.T) {
	// vips does not encode GIF
	operations := flatOperations(t)
	for i := range operations {
		operations[i].Options.Format = image.GIF
	}

	actual, backends, err := transform(t, vipsFirst, "giphy.gif", "image/gif", operations)
	require.NoError(t, err)
	assert.Equal(t, []string{"goimage", "goimage"}, backends)
	assert.Equal(t, 200, decodeConfig(t, actual).Width)
}

func TestTransformNoBackendForOperation(t *testing.T) {
	vipsOnly := config.Config{Backends: &config.Backends{
		Vips: &config.VipsBackend{Mimetypes: MimeTypes},
	}}
	resize := flatOperations(t)[0]

	t.Run("last operation", func(t *testing.T) {
		_, _, err := transform(t, vipsOnly, "schwarzy.jpg", "image/png", []EngineOperation{resize, sepiaOperation})
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
		assert.ErrorContains(t, err, "operation effect")
	})

	t.Run("intermediate operation", func(t *testing.T) {
		_, _, err := transform(t, vipsOnly, "schwarzy.jpg", "image/png", []EngineOperation{sepiaOperation, resize})
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
		assert.ErrorContains(t, err, "operation effect")
	})

	t.Run("intermediate operation with a fallback for the next one", func(t *testing.T) {
		// GoImage would fail on the empty stream left by the effect with a misleading error
		_, _, err := transform(t, vipsFirst, "schwarzy.jpg", "image/png", []EngineOperation{sepiaOperation, resize})
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
		assert.ErrorContains(t, err, "operation effect")
	})

	t.Run("no backend for the content type", func(t *testing.T) {
		vipsJPEGOnly := config.Config{Backends: &config.Backends{
			Vips: &config.VipsBackend{Mimetypes: []string{"image/jpeg"}},
		}}
		dst := &bytes.Buffer{}
		output := &image.ImageFile{
			Stream:  io.NopCloser(bytes.NewReader(readFixture(t, "giphy.gif"))),
			Headers: map[string]string{"Content-Type": "image/gif"},
		}
		_, err := New(vipsJPEGOnly, slog.New(slog.DiscardHandler)).Transform(context.Background(), dst, output,
			[]EngineOperation{{Operation: Resize, Options: &backend.Options{Width: 10, Format: image.GIF}}})
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
		assert.ErrorContains(t, err, "image/gif")
		assert.Zero(t, dst.Len())
	})
}
