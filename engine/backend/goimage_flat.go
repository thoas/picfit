package backend

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"

	imagefile "github.com/thoas/picfit/image"
)

func (e *GoImage) Flat(backgroundFile *imagefile.ImageFile, options *Options) ([]byte, error) {
	if options.Format == imaging.GIF {
		return e.TransformGIF(backgroundFile, options, imaging.Resize)
	}

	background, err := e.Source(backgroundFile)
	if err != nil {
		return nil, err
	}

	images := make([]image.Image, len(options.Images))
	for i := range options.Images {
		images[i], err = e.Source(&options.Images[i])
		if err != nil {
			return nil, err
		}
	}

	bg := image.NewRGBA(image.Rectangle{image.Point{}, background.Bounds().Size()})
	draw.Draw(bg, background.Bounds(), background, image.Point{}, draw.Src)

	fg := drawForeground(foregroundImage(bg, options.Color), images, options)
	draw.Draw(bg, fg.Bounds(), fg, image.ZP, draw.Over)

	return e.ToBytes(bg, options.Format, options.Quality)
}

func foregroundImage(bg image.Image, c string) *image.RGBA {
	b := bg.Bounds()
	fg := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{b.Dx(), b.Dy() / 6},
	})
	draw.Draw(fg, fg.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	return fg
}

func drawForeground(fg *image.RGBA, images []image.Image, options *Options) *image.RGBA {
	n := len(images)
	if n == 0 {
		return fg
	}

	// resize images for foreground
	b := fg.Bounds()
	opts := &Options{
		Width:  b.Dx() / n,
		Height: b.Dy() / n,
	}

	for i := range images {
		images[i] = scale(images[i], opts, imaging.Fit)
	}

	if b.Dx() > b.Dy() {
		return foregroundHorizontal(fg, images, options)
	} else {
		return foregroundVertical(fg, images, options)
	}
}

func foregroundHorizontal(fg *image.RGBA, images []image.Image, options *Options) *image.RGBA {
	position := image.Point{0, 0}
	totalHeight := fg.Bounds().Dy()
	for i := range images {
		bounds := images[i].Bounds()
		position.Y = (totalHeight - bounds.Dy()) / 2
		r := image.Rectangle{
			position,
			position.Add(image.Point{bounds.Dx(), bounds.Dy()}),
		}
		draw.Draw(fg, r, images[i], image.Point{}, draw.Over)
		position.X = position.X + bounds.Dx()
	}
	return fg
}

func foregroundVertical(fg *image.RGBA, images []image.Image, options *Options) *image.RGBA {
	position := image.Point{0, 0}
	cellHeight := fg.Bounds().Dy() / len(images)
	for i := range images {
		bounds := images[i].Bounds()
		position.Y = (cellHeight - bounds.Dy()) / 2
		r := image.Rectangle{
			position,
			position.Add(image.Point{bounds.Dx(), bounds.Dy()}),
		}
		draw.Draw(fg, r, images[i], image.Point{}, draw.Over)
		position.Y = position.Y + bounds.Dy()
	}
	return fg
}
