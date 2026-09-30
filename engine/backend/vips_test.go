package backend

import (
	"bytes"
	"context"
	"fmt"
	stdimage "image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	vips "github.com/cshum/vipsgen/vips816"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/bmp"

	"github.com/thoas/picfit/constants"
	imagefile "github.com/thoas/picfit/image"
)

type backendFunc func(Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error

var (
	resizeFunc backendFunc = func(b Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error {
		return b.Resize
	}
	thumbnailFunc backendFunc = func(b Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error {
		return b.Thumbnail
	}
	fitFunc    backendFunc = func(b Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error { return b.Fit }
	flipFunc   backendFunc = func(b Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error { return b.Flip }
	rotateFunc backendFunc = func(b Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error {
		return b.Rotate
	}
	effectFunc backendFunc = func(b Backend) func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error {
		return b.Effect
	}
)

func newTestFile(data []byte) *imagefile.ImageFile {
	return &imagefile.ImageFile{Stream: io.NopCloser(bytes.NewReader(data))}
}

func readFixture(t testing.TB, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", name))
	require.NoError(t, err)
	return data
}

func run(t *testing.T, b Backend, fn backendFunc, data []byte, opts Options) []byte {
	t.Helper()
	dst := &bytes.Buffer{}
	require.NoError(t, fn(b)(context.Background(), dst, newTestFile(data), &opts))
	return dst.Bytes()
}

func decodeConfig(t *testing.T, data []byte) (stdimage.Config, string) {
	t.Helper()
	cfg, format, err := stdimage.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	return cfg, format
}

func TestVipsParityWithGoImage(t *testing.T) {
	goimage := &GoImage{}
	vipsBackend := NewVips(0, nil)

	fixtures := []string{"schwarzy.jpg", "avatar.png", "giphy.gif"}
	formats := map[imagefile.Format]string{
		imagefile.JPEG: "jpeg",
		imagefile.PNG:  "png",
		imagefile.WEBP: "webp",
	}
	operations := []struct {
		name string
		fn   backendFunc
		opts Options
	}{
		{"resize w", resizeFunc, Options{Width: 100, Upscale: true}},
		{"resize h", resizeFunc, Options{Height: 100, Upscale: true}},
		{"resize wxh stretch", resizeFunc, Options{Width: 120, Height: 40, Upscale: true}},
		{"resize upscale", resizeFunc, Options{Width: 800, Upscale: true}},
		{"resize no upscale", resizeFunc, Options{Width: 800, Upscale: false}},
		{"resize no dimension", resizeFunc, Options{Upscale: true}},
		{"thumbnail", thumbnailFunc, Options{Width: 100, Height: 100, Upscale: true}},
		{"thumbnail wide", thumbnailFunc, Options{Width: 200, Height: 50, Upscale: true}},
		{"thumbnail upscale", thumbnailFunc, Options{Width: 900, Height: 600, Upscale: true}},
		{"thumbnail no upscale", thumbnailFunc, Options{Width: 900, Height: 600, Upscale: false}},
		{"fit", fitFunc, Options{Width: 100, Height: 100, Upscale: true}},
		{"fit wide", fitFunc, Options{Width: 300, Height: 50, Upscale: true}},
		{"fit bigger than source", fitFunc, Options{Width: 1000, Height: 1000, Upscale: true}},
		{"flip h", flipFunc, Options{Position: "h"}},
		{"flip v", flipFunc, Options{Position: "v"}},
		{"rotate 90", rotateFunc, Options{Degree: 90}},
		{"rotate 180", rotateFunc, Options{Degree: 180}},
		{"rotate 270", rotateFunc, Options{Degree: 270}},
		{"blur", effectFunc, Options{Filter: constants.FilterBlur}},
	}

	for _, fixture := range fixtures {
		data := readFixture(t, fixture)
		for format, formatName := range formats {
			for _, op := range operations {
				t.Run(fmt.Sprintf("%s/%s/%s", fixture, formatName, op.name), func(t *testing.T) {
					opts := op.opts
					opts.Format = format
					opts.Quality = 85

					expected, _ := decodeConfig(t, run(t, goimage, op.fn, data, opts))
					actual, actualFormat := decodeConfig(t, run(t, vipsBackend, op.fn, data, opts))

					assert.Equal(t, formatName, actualFormat)
					assert.Equal(t, expected.Width, actual.Width, "width")
					assert.Equal(t, expected.Height, actual.Height, "height")
				})
			}
		}
	}
}

// quadrant returns a 40x20 PNG with a red top-left corner, the rest is blue.
func quadrant(t *testing.T) []byte {
	t.Helper()
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 40, 20))
	for y := range 20 {
		for x := range 40 {
			c := color.RGBA{B: 255, A: 255}
			if x < 20 && y < 10 {
				c = color.RGBA{R: 255, A: 255}
			}
			img.Set(x, y, c)
		}
	}
	buf := &bytes.Buffer{}
	require.NoError(t, png.Encode(buf, img))
	return buf.Bytes()
}

func redCorner(t *testing.T, data []byte) string {
	t.Helper()
	img, _, err := stdimage.Decode(bytes.NewReader(data))
	require.NoError(t, err)
	b := img.Bounds()
	corners := map[string]stdimage.Point{
		"top-left":     {b.Min.X + 1, b.Min.Y + 1},
		"top-right":    {b.Max.X - 2, b.Min.Y + 1},
		"bottom-left":  {b.Min.X + 1, b.Max.Y - 2},
		"bottom-right": {b.Max.X - 2, b.Max.Y - 2},
	}
	for name, p := range corners {
		r, _, bl, _ := img.At(p.X, p.Y).RGBA()
		if r > bl {
			return name
		}
	}
	return "none"
}

func TestVipsOrientationParityWithGoImage(t *testing.T) {
	goimage := &GoImage{}
	vipsBackend := NewVips(0, nil)
	data := quadrant(t)

	tests := []struct {
		name string
		fn   backendFunc
		opts Options
	}{
		{"rotate 90", rotateFunc, Options{Degree: 90}},
		{"rotate 180", rotateFunc, Options{Degree: 180}},
		{"rotate 270", rotateFunc, Options{Degree: 270}},
		{"flip h", flipFunc, Options{Position: "h"}},
		{"flip v", flipFunc, Options{Position: "v"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.opts
			opts.Format = imagefile.PNG

			expected := redCorner(t, run(t, goimage, tt.fn, data, opts))
			actual := redCorner(t, run(t, vipsBackend, tt.fn, data, opts))
			assert.NotEqual(t, "none", expected)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestVipsExifOrientation(t *testing.T) {
	goimage := &GoImage{}
	vipsBackend := NewVips(0, nil)

	img, err := vips.NewImageFromBuffer(readFixture(t, "schwarzy.jpg"), nil)
	require.NoError(t, err)
	defer img.Close()
	require.NoError(t, img.SetOrientation(6))
	data, err := img.JpegsaveBuffer(&vips.JpegsaveBufferOptions{Keep: vips.KeepAll})
	require.NoError(t, err)

	opts := Options{Width: 100, Upscale: true, Format: imagefile.JPEG}
	expected, _ := decodeConfig(t, run(t, goimage, resizeFunc, data, opts))
	actual, _ := decodeConfig(t, run(t, vipsBackend, resizeFunc, data, opts))

	// schwarzy.jpg is 500x357, orientation 6 makes it portrait.
	assert.Equal(t, 100, actual.Width)
	assert.Greater(t, actual.Height, actual.Width)
	assert.Equal(t, expected.Width, actual.Width)
	assert.Equal(t, expected.Height, actual.Height)
}

func TestVipsNotImplementedKeepsStream(t *testing.T) {
	vipsBackend := NewVips(0, nil)

	bmpData := &bytes.Buffer{}
	require.NoError(t, bmp.Encode(bmpData, stdimage.NewRGBA(stdimage.Rect(0, 0, 10, 10))))

	tests := []struct {
		name string
		fn   func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error
		data []byte
		opts Options
	}{
		{"gif output", vipsBackend.Resize, readFixture(t, "giphy.gif"), Options{Width: 100, Format: imagefile.GIF}},
		{"unknown filter", vipsBackend.Effect, readFixture(t, "schwarzy.jpg"), Options{Filter: "sepia", Format: imagefile.JPEG}},
		{"flat", vipsBackend.Flat, readFixture(t, "schwarzy.jpg"), Options{Format: imagefile.JPEG}},
		{"bmp source", vipsBackend.Resize, bmpData.Bytes(), Options{Width: 5, Format: imagefile.JPEG}},
		{"not an image", vipsBackend.Resize, []byte("not an image"), Options{Width: 5, Format: imagefile.JPEG}},
		{"empty", vipsBackend.Resize, []byte{}, Options{Width: 5, Format: imagefile.JPEG}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := newTestFile(tt.data)
			dst := &bytes.Buffer{}
			err := tt.fn(context.Background(), dst, file, &tt.opts)
			assert.ErrorIs(t, err, MethodNotImplementedError)
			assert.Zero(t, dst.Len())

			rest, err := io.ReadAll(file.Stream)
			require.NoError(t, err)
			assert.Equal(t, tt.data, rest, "the next backend must get the full stream")
		})
	}
}

func TestVipsTruncatedImage(t *testing.T) {
	vipsBackend := NewVips(0, nil)
	data := readFixture(t, "schwarzy.jpg")
	truncated := data[:len(data)/2]

	tests := []struct {
		name string
		fn   func(context.Context, io.Writer, *imagefile.ImageFile, *Options) error
		opts Options
	}{
		{"resize", vipsBackend.Resize, Options{Width: 100, Upscale: true, Format: imagefile.JPEG}},
		{"thumbnail", vipsBackend.Thumbnail, Options{Width: 100, Height: 100, Upscale: true, Format: imagefile.JPEG}},
		{"rotate", vipsBackend.Rotate, Options{Degree: 90, Format: imagefile.JPEG}},
		{"no resize", vipsBackend.Resize, Options{Format: imagefile.PNG}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(context.Background(), io.Discard, newTestFile(truncated), &tt.opts)
			assert.Error(t, err)
			assert.NotErrorIs(t, err, MethodNotImplementedError)
		})
	}
}

func TestVipsInvalidParameters(t *testing.T) {
	vipsBackend := NewVips(0, nil)
	data := readFixture(t, "schwarzy.jpg")

	err := vipsBackend.Rotate(context.Background(), io.Discard, newTestFile(data), &Options{Degree: 45, Format: imagefile.JPEG})
	assert.ErrorContains(t, err, "degree=45")

	err = vipsBackend.Flip(context.Background(), io.Discard, newTestFile(data), &Options{Position: "x", Format: imagefile.JPEG})
	assert.ErrorContains(t, err, "x is not supported")
}

func TestVipsCanceledContext(t *testing.T) {
	vipsBackend := NewVips(0, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := vipsBackend.Resize(ctx, io.Discard, newTestFile(readFixture(t, "schwarzy.jpg")),
		&Options{Width: 100, Format: imagefile.JPEG})
	assert.ErrorIs(t, err, context.Canceled)
}

func TestVipsConcurrent(t *testing.T) {
	vipsBackend := NewVips(0, nil)
	data := readFixture(t, "schwarzy.jpg")

	var wg sync.WaitGroup
	errs := make(chan error, 32)
	for i := range 32 {
		wg.Go(func() {
			opts := &Options{Width: 50 + i, Height: 50, Upscale: true, Format: imagefile.WEBP, Quality: 80}
			dst := &bytes.Buffer{}
			if err := vipsBackend.Thumbnail(context.Background(), dst, newTestFile(data), opts); err != nil {
				errs <- err
				return
			}
			cfg, _, err := stdimage.DecodeConfig(dst)
			if err != nil {
				errs <- err
				return
			}
			if cfg.Width != opts.Width || cfg.Height != opts.Height {
				errs <- fmt.Errorf("got %dx%d, want %dx%d", cfg.Width, cfg.Height, opts.Width, opts.Height)
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

func BenchmarkBackends(b *testing.B) {
	data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "original.jpg"))
	if err != nil {
		b.Skip("tests/fixtures/original.jpg is missing")
	}

	backends := []Backend{&GoImage{}, NewVips(0, nil)}
	operations := []struct {
		name string
		fn   backendFunc
		opts Options
	}{
		{"resize", resizeFunc, Options{Width: 50, Height: 50, Upscale: true}},
		{"thumbnail", thumbnailFunc, Options{Width: 200, Height: 200, Upscale: true}},
		{"fit", fitFunc, Options{Width: 800, Height: 600, Upscale: true}},
		{"rotate", rotateFunc, Options{Degree: 90}},
		{"flip", flipFunc, Options{Position: "h"}},
		{"blur", effectFunc, Options{Filter: constants.FilterBlur}},
	}

	for _, op := range operations {
		for _, backend := range backends {
			b.Run(fmt.Sprintf("%s/%s", op.name, backend), func(b *testing.B) {
				opts := op.opts
				opts.Format = imagefile.JPEG
				opts.Quality = 95
				b.ReportAllocs()
				for b.Loop() {
					if err := op.fn(backend)(context.Background(), io.Discard, newTestFile(data), &opts); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
