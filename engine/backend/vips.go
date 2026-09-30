package backend

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"slices"
	"sync"

	vips "github.com/cshum/vipsgen/vips816"
	"github.com/pkg/errors"

	"github.com/thoas/picfit/constants"
	imagefile "github.com/thoas/picfit/image"
)

var (
	vipsStartup sync.Once

	vipsFlipDirections = map[string]vips.Direction{
		"h": vips.DirectionHorizontal,
		"v": vips.DirectionVertical,
	}

	// imaging rotates counter-clockwise whereas libvips rotates clockwise.
	vipsRotateAngles = map[int]vips.Angle{
		90:  vips.AngleD270,
		180: vips.AngleD180,
		270: vips.AngleD90,
	}

	// vipsSourceTypes are the source formats handled by libvips,
	// the others are handed back to the next backend.
	vipsSourceTypes = []vips.ImageType{
		vips.ImageTypeJpeg,
		vips.ImageTypePng,
		vips.ImageTypeWebp,
		vips.ImageTypeGif,
		vips.ImageTypeTiff,
	}
)

type sizeFunc func(srcWidth, srcHeight, width, height int) (int, int)

// Vips is the libvips backend.
type Vips struct {
	logger *slog.Logger
}

// NewVips starts libvips once for the process and returns the backend.
// concurrency is the number of libvips worker threads per image, 0 uses the libvips default.
func NewVips(concurrency int, logger *slog.Logger) *Vips {
	if logger == nil {
		logger = slog.Default()
	}

	vipsStartup.Do(func() {
		// Read by vips_init: disables loaders not designed for untrusted input (magick, pdf...).
		if _, ok := os.LookupEnv("VIPS_BLOCK_UNTRUSTED"); !ok {
			os.Setenv("VIPS_BLOCK_UNTRUSTED", "1")
		}

		vips.SetLogging(func(domain string, level vips.LogLevel, message string) {
			logger.Warn(message, slog.String("domain", domain))
		}, vips.LogLevelWarning)

		vips.Startup(&vips.Config{
			ConcurrencyLevel: concurrency,
			VectorEnabled:    true,
		})
	})

	return &Vips{logger: logger}
}

func (b *Vips) String() string {
	return "vips"
}

// Resize implements Backend.
func (b *Vips) Resize(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	return b.scale(ctx, dst, img, options, resizeSize, vips.InterestingNone)
}

// Thumbnail implements Backend.
func (b *Vips) Thumbnail(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	return b.scale(ctx, dst, img, options, resizeSize, vips.InterestingCentre)
}

// Fit implements Backend.
func (b *Vips) Fit(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	return b.scale(ctx, dst, img, options, fitSize, vips.InterestingNone)
}

// Flip implements Backend.
func (b *Vips) Flip(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	direction, ok := vipsFlipDirections[options.Position]
	if !ok {
		return fmt.Errorf("Invalid flip transformation, %s is not supported", options.Position)
	}

	return b.transform(ctx, dst, img, options, func(image *vips.Image) error {
		return image.Flip(direction)
	})
}

// Rotate implements Backend.
func (b *Vips) Rotate(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	angle, ok := vipsRotateAngles[options.Degree]
	if !ok {
		return fmt.Errorf("Invalid rotate transformation degree=%d is not supported", options.Degree)
	}

	return b.transform(ctx, dst, img, options, func(image *vips.Image) error {
		return image.Rot(angle)
	})
}

// Effect implements Backend.
func (b *Vips) Effect(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options) error {
	if options.Filter != constants.FilterBlur {
		return MethodNotImplementedError
	}

	return b.transform(ctx, dst, img, options, func(image *vips.Image) error {
		const maxSigma = 50
		sigma := min(max(image.Width(), image.Height())/20, maxSigma)
		if sigma <= 0 {
			return nil
		}

		return image.Gaussblur(float64(sigma), nil)
	})
}

// scale follows the GoImage rules: no upscale unless requested, output dimensions computed like imaging.
func (b *Vips) scale(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options, size sizeFunc, crop vips.Interesting) error {
	if !vipsCanSave(options.Format) {
		return MethodNotImplementedError
	}

	data, image, err := b.load(ctx, img)
	if err != nil {
		return err
	}
	defer image.Close()

	srcWidth, srcHeight := image.Width(), image.Height()
	if options.Width == 0 && options.Height == 0 {
		return b.save(dst, image, options)
	}

	factor := scalingFactor(srcWidth, srcHeight, options.Width, options.Height)
	if factor >= 1 && !options.Upscale {
		return b.save(dst, image, options)
	}

	width, height := size(srcWidth, srcHeight, options.Width, options.Height)
	if width == srcWidth && height == srcHeight {
		return b.save(dst, image, options)
	}

	// SizeForce gives the exact computed dimensions but stretches instead of cropping.
	thumbnailSize := vips.SizeForce
	if crop != vips.InterestingNone {
		thumbnailSize = vips.SizeBoth
	}

	thumbnail, err := vips.NewThumbnailBuffer(data, width, &vips.ThumbnailBufferOptions{
		Height: height,
		Size:   thumbnailSize,
		Crop:   crop,
		// thumbnail_buffer does not forward fail_on to the loader, truncated images
		// would be rendered grey instead of failing like GoImage.
		OptionString: "fail_on=error",
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer thumbnail.Close()

	return b.save(dst, thumbnail, options)
}

func (b *Vips) transform(ctx context.Context, dst io.Writer, img *imagefile.ImageFile, options *Options, trans func(*vips.Image) error) error {
	if !vipsCanSave(options.Format) {
		return MethodNotImplementedError
	}

	_, image, err := b.load(ctx, img)
	if err != nil {
		return err
	}
	defer image.Close()

	if err := trans(image); err != nil {
		return errors.WithStack(err)
	}

	return b.save(dst, image, options)
}

// load reads the source and opens it lazily, only the header is decoded here.
// Sources libvips cannot open are handed back to the next backend with a rewound stream.
func (b *Vips) load(ctx context.Context, img *imagefile.ImageFile) ([]byte, *vips.Image, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	data, err := io.ReadAll(img.Stream)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	fallback := func(reason string) ([]byte, *vips.Image, error) {
		b.logger.WarnContext(ctx, "vips cannot handle image, falling back",
			slog.String("image", img.Filepath),
			slog.String("reason", reason),
		)
		rewind(img, data)
		return nil, nil, MethodNotImplementedError
	}

	// vipsgen panics on an empty buffer.
	if len(data) == 0 {
		return fallback("empty image")
	}

	image, err := vips.NewImageFromBuffer(data, &vips.LoadOptions{FailOn: vips.FailOnError})
	if err != nil {
		return fallback(err.Error())
	}

	if format := image.Format(); !slices.Contains(vipsSourceTypes, format) {
		image.Close()
		return fallback(fmt.Sprintf("unsupported source format %s", format))
	}

	if err := image.Autorot(nil); err != nil {
		image.Close()
		return nil, nil, errors.WithStack(err)
	}

	return data, image, nil
}

func (b *Vips) save(dst io.Writer, image *vips.Image, options *Options) error {
	if image.Interpretation() == vips.InterpretationCmyk {
		if err := image.Colourspace(vips.InterpretationSrgb, nil); err != nil {
			return errors.WithStack(err)
		}
	}

	var (
		buf []byte
		err error
	)
	switch options.Format {
	case imagefile.JPEG:
		buf, err = image.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
			Q:              options.Quality,
			OptimizeCoding: true,
			Keep:           vips.KeepNone,
		})
	case imagefile.PNG:
		buf, err = image.PngsaveBuffer(&vips.PngsaveBufferOptions{
			Keep: vips.KeepNone,
		})
	case imagefile.WEBP:
		buf, err = image.WebpsaveBuffer(&vips.WebpsaveBufferOptions{
			Q:    options.Quality,
			Keep: vips.KeepNone,
		})
	default:
		return MethodNotImplementedError
	}
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = dst.Write(buf)
	return errors.WithStack(err)
}

func vipsCanSave(format imagefile.Format) bool {
	return format == imagefile.JPEG || format == imagefile.PNG || format == imagefile.WEBP
}

// resizeSize computes the dimensions like imaging.Resize: a zero dimension preserves the aspect ratio.
func resizeSize(srcWidth, srcHeight, width, height int) (int, int) {
	if width == 0 {
		width = int(max(1.0, math.Floor(float64(height)*float64(srcWidth)/float64(srcHeight)+0.5)))
	}
	if height == 0 {
		height = int(max(1.0, math.Floor(float64(width)*float64(srcHeight)/float64(srcWidth)+0.5)))
	}

	return width, height
}

// fitSize computes the dimensions like imaging.Fit: the image is bounded by width x height, never upscaled.
func fitSize(srcWidth, srcHeight, width, height int) (int, int) {
	if width == 0 || height == 0 {
		return resizeSize(srcWidth, srcHeight, width, height)
	}
	if srcWidth <= width && srcHeight <= height {
		return srcWidth, srcHeight
	}

	srcRatio := float64(srcWidth) / float64(srcHeight)
	if srcRatio > float64(width)/float64(height) {
		return resizeSize(srcWidth, srcHeight, width, int(float64(width)/srcRatio))
	}

	return resizeSize(srcWidth, srcHeight, int(float64(height)*srcRatio), height)
}

var _ Backend = (*Vips)(nil)
