package backend

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/thoas/picfit/image"
)

// Gifsicle is the gifsicle backend.
type Gifsicle struct{}

func (b *Gifsicle) String() string {
	return "gifsicle"
}

// Fit implements Backend.
func (b *Gifsicle) Fit(*image.ImageFile, *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Flat implements Backend.
func (b *Gifsicle) Flat(*image.ImageFile, *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Flip implements Backend.
func (b *Gifsicle) Flip(*image.ImageFile, *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Resize implements Backend.
func (b *Gifsicle) Resize(*image.ImageFile, *Options) ([]byte, error) {
	return nil, MethodNotImplementedError

}

// Rotate implements Backend.
func (b *Gifsicle) Rotate(*image.ImageFile, *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

// Thumbnail implements Backend.
func (b *Gifsicle) Thumbnail(imgfile *image.ImageFile, opts *Options) ([]byte, error) {
	cmd := exec.Command("gifsicle", "--resize", fmt.Sprintf("%dx%d", opts.Width, opts.Height))
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
