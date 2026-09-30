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
		{"flat gif output", vipsBackend.Flat, readFixture(t, "giphy.gif"), Options{Position: "10.10.50.50", Format: imagefile.GIF}},
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

// pixelDiff returns the mean absolute difference per channel (0-255) and the ratio
// of pixels differing by more than 64 on a channel, a misplaced overlay raises the latter.
func pixelDiff(t *testing.T, a, b []byte) (float64, float64) {
	t.Helper()
	imgA, _, err := stdimage.Decode(bytes.NewReader(a))
	require.NoError(t, err)
	imgB, _, err := stdimage.Decode(bytes.NewReader(b))
	require.NoError(t, err)
	require.Equal(t, imgA.Bounds().Size(), imgB.Bounds().Size())

	var sum float64
	var far, n int
	bounds := imgA.Bounds()
	for y := range bounds.Dy() {
		for x := range bounds.Dx() {
			r1, g1, b1, _ := imgA.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			r2, g2, b2, _ := imgB.At(imgB.Bounds().Min.X+x, imgB.Bounds().Min.Y+y).RGBA()
			worst := 0
			for _, d := range []int{int(r1>>8) - int(r2>>8), int(g1>>8) - int(g2>>8), int(b1>>8) - int(b2>>8)} {
				d = max(d, -d)
				sum += float64(d)
				worst = max(worst, d)
			}
			if worst > 64 {
				far++
			}
			n++
		}
	}
	return sum / float64(3*n), float64(far) / float64(n)
}

func TestVipsFlatParityWithGoImage(t *testing.T) {
	goimage := &GoImage{}
	vipsBackend := NewVips(0, nil)

	tests := []struct {
		name        string
		background  string
		foregrounds []string
		opts        Options
	}{
		{"pos one image", "schwarzy.jpg", []string{"avatar.png"}, Options{Position: "10.10.60.60"}},
		{"pos horizontal", "schwarzy.jpg", []string{"avatar.png", "giphy.gif"}, Options{Position: "10.20.90.60"}},
		{"pos vertical", "schwarzy.jpg", []string{"avatar.png", "giphy.gif"}, Options{Position: "70.5.95.95"}},
		{"pos color", "schwarzy.jpg", []string{"giphy.gif"}, Options{Position: "10.10.50.90", Color: "ff0000"}},
		{"pos color without image", "schwarzy.jpg", nil, Options{Position: "20.20.80.80", Color: "00ff00"}},
		{"pos invalid color", "schwarzy.jpg", []string{"avatar.png"}, Options{Position: "10.10.60.60", Color: "zz"}},
		{"pos beyond background", "schwarzy.jpg", []string{"avatar.png"}, Options{Position: "80.80.150.150", Color: "0000ff"}},
		{"pos empty", "schwarzy.jpg", []string{"avatar.png"}, Options{}},
		{"stick top-left", "schwarzy.jpg", []string{"avatar.png"}, Options{Stick: constants.TopLeft, Width: 100, Height: 80}},
		{"stick top-right", "schwarzy.jpg", []string{"giphy.gif"}, Options{Stick: constants.TopRight, Width: 100, Height: 80}},
		{"stick bottom-left", "schwarzy.jpg", []string{"avatar.png"}, Options{Stick: constants.BottomLeft, Width: 100, Height: 80}},
		{"stick bottom-right", "schwarzy.jpg", []string{"giphy.gif"}, Options{Stick: constants.BottomRight, Width: 100, Height: 80}},
		{"stick width only", "schwarzy.jpg", []string{"avatar.png"}, Options{Stick: constants.BottomRight, Width: 120}},
		{"stick without dimension", "schwarzy.jpg", []string{"avatar.png"}, Options{Stick: constants.TopLeft}},
		{"transparent background", "giphy.gif", []string{"avatar.png"}, Options{Position: "25.25.75.75"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flat := func(b Backend) []byte {
				opts := tt.opts
				opts.Format = imagefile.PNG
				for _, name := range tt.foregrounds {
					opts.Images = append(opts.Images, *newTestFile(readFixture(t, name)))
				}
				dst := &bytes.Buffer{}
				require.NoError(t, b.Flat(context.Background(), dst, newTestFile(readFixture(t, tt.background)), &opts))
				return dst.Bytes()
			}

			expected, actual := flat(goimage), flat(vipsBackend)
			mean, far := pixelDiff(t, expected, actual)
			assert.Less(t, mean, 3.0, "mean difference per channel")
			assert.Less(t, far, 0.01, "ratio of very different pixels")
		})
	}
}

func TestVipsFlatFallbackRewindsEveryStream(t *testing.T) {
	vipsBackend := NewVips(0, nil)

	bmpData := &bytes.Buffer{}
	require.NoError(t, bmp.Encode(bmpData, stdimage.NewRGBA(stdimage.Rect(0, 0, 10, 10))))
	background, avatar := readFixture(t, "schwarzy.jpg"), readFixture(t, "avatar.png")

	bgFile := newTestFile(background)
	opts := &Options{
		Position: "10.10.60.60",
		Format:   imagefile.JPEG,
		Images:   []imagefile.ImageFile{*newTestFile(avatar), *newTestFile(bmpData.Bytes())},
	}
	err := vipsBackend.Flat(context.Background(), io.Discard, bgFile, opts)
	require.ErrorIs(t, err, MethodNotImplementedError)

	for name, tt := range map[string]struct {
		file *imagefile.ImageFile
		data []byte
	}{
		"background": {bgFile, background},
		"avatar":     {&opts.Images[0], avatar},
		"bmp":        {&opts.Images[1], bmpData.Bytes()},
	} {
		rest, err := io.ReadAll(tt.file.Stream)
		require.NoError(t, err)
		assert.Equal(t, tt.data, rest, name)
	}
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
	sources := map[string][]byte{}
	for _, name := range []string{"original.jpg", "schwarzy.jpg"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", name))
		if err != nil {
			b.Skipf("tests/fixtures/%s is missing", name)
		}
		sources[name] = data
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

	for _, source := range []string{"original.jpg", "schwarzy.jpg"} {
		for _, op := range operations {
			for _, backend := range backends {
				// key=value names so that benchstat -col /backend compares the backends
				b.Run(fmt.Sprintf("src=%s/op=%s/backend=%s", source, op.name, backend), func(b *testing.B) {
					opts := op.opts
					opts.Format = imagefile.JPEG
					opts.Quality = 95
					b.ReportAllocs()
					for b.Loop() {
						if err := op.fn(backend)(context.Background(), io.Discard, newTestFile(sources[source]), &opts); err != nil {
							b.Fatal(err)
						}
					}
				})
			}
		}
	}
}

func BenchmarkFlat(b *testing.B) {
	background, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "original.jpg"))
	if err != nil {
		b.Skip("tests/fixtures/original.jpg is missing")
	}
	foreground, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "avatar.png"))
	if err != nil {
		b.Fatal(err)
	}

	for _, backend := range []Backend{&GoImage{}, NewVips(0, nil)} {
		b.Run(fmt.Sprintf("backend=%s", backend), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				opts := &Options{
					Position: "10.10.60.60",
					Color:    "ffffff",
					Format:   imagefile.JPEG,
					Quality:  95,
					Images:   []imagefile.ImageFile{*newTestFile(foreground)},
				}
				if err := backend.Flat(context.Background(), io.Discard, newTestFile(background), opts); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
