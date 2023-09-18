package backend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/gif"
	"os/exec"

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
func (b *Gifsicle) Resize(ctx context.Context, imgfile *image.ImageFile, opts *Options) ([]byte, error) {
	img, err := gif.Decode(bytes.NewReader(imgfile.Source))
	if err != nil {
		return nil, err
	}
	factor := scalingFactorImage(img, opts.Width, opts.Height)
	if factor > 1 && !opts.Upscale {
		return imgfile.Source, nil
	}

	cmd := exec.CommandContext(ctx, b.Path,
		"--resize", fmt.Sprintf("%dx%d", opts.Width, opts.Height),
	)
	cmd.Stdin = bytes.NewReader(imgfile.Source)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	var target *exec.ExitError
	if err := cmd.Run(); errors.As(err, &target) && target.Exited() {
		return nil, errors.New(stderr.String())
	} else if err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

// Thumbnail implements Backend.
func (b *Gifsicle) Thumbnail(ctx context.Context, imgfile *image.ImageFile, opts *Options) ([]byte, error) {
	img, err := gif.Decode(bytes.NewReader(imgfile.Source))
	if err != nil {
		return nil, err
	}
	factor := scalingFactorImage(img, opts.Width, opts.Height)
	if factor > 1 && !opts.Upscale {
		return imgfile.Source, nil
	}

	bounds := img.Bounds()
	left, top, cropw, croph := computecrop(bounds.Dx(), bounds.Dy(), opts.Width, opts.Height)

	cmd := exec.CommandContext(ctx, b.Path,
		"--crop", fmt.Sprintf("%d,%d+%dx%d", left, top, cropw, croph),
		"--resize", fmt.Sprintf("%dx%d", opts.Width, opts.Height),
	)
	cmd.Stdin = bytes.NewReader(imgfile.Source)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	var target *exec.ExitError
	if err := cmd.Run(); errors.As(err, &target) && target.Exited() {
		return nil, errors.New(stderr.String())
	} else if err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

// Rotate implements Backend.
func (b *Gifsicle) Rotate(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Fit implements Backend.
func (b *Gifsicle) Fit(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Flat implements Backend.
func (b *Gifsicle) Flat(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Flip implements Backend.
func (b *Gifsicle) Flip(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
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
