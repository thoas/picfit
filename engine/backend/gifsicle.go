package backend

import (
	"bytes"
	"context"
	"fmt"
	"image/gif"
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/thoas/picfit/image"
)

// Gifsicle is the gifsicle backend.
type Gifsicle struct {
	Path string
}

func (b *Gifsicle) String() string {
	return "gifsicle"
}

// Resize implements Backend.
func (b *Gifsicle) Resize(ctx context.Context, dst io.Writer, imgfile *image.ImageFile, opts *Options) error {
	data, err := io.ReadAll(imgfile.Stream)
	if err != nil {
		return errors.WithStack(err)
	}

	img, err := gif.Decode(bytes.NewReader(data))
	if err != nil {
		return errors.WithStack(err)
	}
	factor := scalingFactorImage(img, opts.Width, opts.Height)
	if factor > 1 && !opts.Upscale {
		return gif.Encode(dst, img, nil)
	}

	resizeOption := fmt.Sprintf("%dx%d", opts.Width, opts.Height)
	cmd := exec.CommandContext(ctx, b.Path,
		"--resize", resizeOption,
	)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = dst
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	var target *exec.ExitError
	if err := cmd.Run(); errors.As(err, &target) && target.Exited() {
		return errors.New(stderr.String())
	} else if err != nil {
		return errors.Wrap(err, "unable to resize")
	}
	return nil
}

// Thumbnail implements Backend.
func (b *Gifsicle) Thumbnail(ctx context.Context, dst io.Writer, imgfile *image.ImageFile, opts *Options) error {
	data, err := io.ReadAll(imgfile.Stream)
	if err != nil {
		return errors.WithStack(err)
	}

	img, err := gif.Decode(bytes.NewReader(data))
	if err != nil {
		return errors.WithStack(err)
	}
	factor := scalingFactorImage(img, opts.Width, opts.Height)
	if factor > 1 && !opts.Upscale {
		return gif.Encode(dst, img, nil)
	}

	bounds := img.Bounds()
	left, top, cropw, croph := computecrop(bounds.Dx(), bounds.Dy(), opts.Width, opts.Height)
	cropOption := fmt.Sprintf("%d,%d+%dx%d", left, top, cropw, croph)
	resizeOption := fmt.Sprintf("%dx%d", opts.Width, opts.Height)

	cmd := exec.CommandContext(ctx, b.Path,
		"--crop", cropOption,
		"--resize", resizeOption,
	)

	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = dst
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	var target *exec.ExitError
	if err := cmd.Run(); errors.As(err, &target) && target.Exited() {
		return errors.New(stderr.String())
	} else if err != nil {
		return errors.Wrap(err, "unable to thumbnail")
	}
	return nil
}

// Rotate implements Backend.
func (b *Gifsicle) Rotate(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error {
	return MethodNotImplementedError
}

// Fit implements Backend.
func (b *Gifsicle) Fit(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error {
	return MethodNotImplementedError
}

// Effect implements Backend.
func (b *Gifsicle) Effect(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error {
	return MethodNotImplementedError
}

// Flat implements Backend.
func (b *Gifsicle) Flat(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error {
	return MethodNotImplementedError
}

// Flip implements Backend.
func (b *Gifsicle) Flip(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error {
	return MethodNotImplementedError
}

func computecrop(srcw, srch, destw, desth int) (left, top, cropw, croph int) {
	srcratio := float64(srcw) / float64(srch)
	destratio := float64(destw) / float64(desth)

	if srcratio > destratio {
		cropw = int((destratio * float64(srch)) + 0.5)
		croph = srch
	} else {
		croph = int((float64(srcw) / destratio) + 0.5)
		cropw = srcw
	}

	left = int(float64(srcw-cropw) * 0.5)
	if left < 0 {
		left = 0
	}

	top = int(float64(srch-croph) * 0.5)
	if top < 0 {
		top = 0
	}
	return
}
