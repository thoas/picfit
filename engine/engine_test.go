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

func transform(t *testing.T, cfg config.Config, operations []EngineOperation) ([]byte, []string, error) {
	t.Helper()
	logs := &bytes.Buffer{}
	e := New(cfg, slog.New(slog.NewJSONHandler(logs, nil)))

	output := &image.ImageFile{
		Stream:  io.NopCloser(bytes.NewReader(readFixture(t, "schwarzy.jpg"))),
		Headers: map[string]string{"Content-Type": "image/png"},
	}
	dst := &bytes.Buffer{}
	_, err := e.Transform(context.Background(), dst, output, operations)
	return dst.Bytes(), handledBy(t, logs), err
}

func TestTransformFlatFallsBackOnGoImage(t *testing.T) {
	vipsFirst := config.Config{Backends: &config.Backends{
		Vips:    &config.VipsBackend{Mimetypes: MimeTypes},
		GoImage: &config.Backend{Weight: 1, Mimetypes: MimeTypes},
	}}

	actual, backends, err := transform(t, vipsFirst, flatOperations(t))
	require.NoError(t, err)
	assert.Equal(t, []string{"vips", "goimage"}, backends)

	expected, _, err := transform(t, config.Config{}, flatOperations(t))
	require.NoError(t, err)

	actualCfg, _, err := stdimage.DecodeConfig(bytes.NewReader(actual))
	require.NoError(t, err)
	expectedCfg, _, err := stdimage.DecodeConfig(bytes.NewReader(expected))
	require.NoError(t, err)
	assert.Equal(t, expectedCfg, actualCfg)
}

func TestTransformNoBackendForOperation(t *testing.T) {
	vipsOnly := config.Config{Backends: &config.Backends{
		Vips: &config.VipsBackend{Mimetypes: MimeTypes},
	}}

	t.Run("last operation", func(t *testing.T) {
		_, _, err := transform(t, vipsOnly, flatOperations(t))
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
	})

	t.Run("intermediate operation", func(t *testing.T) {
		operations := flatOperations(t)
		operations = append([]EngineOperation{operations[1]}, operations[0])

		_, _, err := transform(t, vipsOnly, operations)
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
		assert.ErrorContains(t, err, "operation flat")
	})

	t.Run("intermediate operation with a fallback for the next one", func(t *testing.T) {
		vipsWithGIFFallback := config.Config{Backends: &config.Backends{
			Vips:    &config.VipsBackend{Mimetypes: MimeTypes},
			GoImage: &config.Backend{Weight: 1, Mimetypes: []string{"image/gif"}},
		}}
		operations := flatOperations(t)
		operations = append([]EngineOperation{operations[1]}, operations[0])

		_, _, err := transform(t, vipsWithGIFFallback, operations)
		assert.ErrorIs(t, err, backend.MethodNotImplementedError)
		assert.ErrorContains(t, err, "operation flat")
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
