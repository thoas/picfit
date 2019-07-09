package backend

import (
	"fmt"
	"image"
	"image/draw"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/thoas/picfit/constants"
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

	// draw background
	bg := image.NewRGBA(image.Rectangle{image.Point{}, background.Bounds().Size()})
	draw.Draw(bg, background.Bounds(), background, image.Point{}, draw.Src)

	if options.Stick != "" {
		drawStickForeground(bg, images, options)
	} else {
		drawPosForeground(bg, images, options)
	}

	return e.ToBytes(bg, options.Format, options.Quality)
}

func drawStickForeground(bg *image.RGBA, images []image.Image, options *Options) {
	for i := range images {
		opts := &Options{
			Upscale: true,
			Width:   options.Width,
			Height:  options.Height,
		}

		images[i] = scale(images[i], opts, imaging.Resize)

		bounds := images[i].Bounds()
		var position image.Point
		switch options.Stick {
		case constants.TopLeft:
			position = bounds.Min
		case constants.TopRight:
			position = image.Point{X: bg.Bounds().Dx() - bounds.Dx(), Y: 0}
		case constants.BottomLeft:
			position = image.Point{X: 0, Y: bg.Bounds().Dy() - bounds.Dy()}
		case constants.BottomRight:
			position = image.Point{
				X: bg.Bounds().Dx() - bounds.Dx(),
				Y: bg.Bounds().Dy() - bounds.Dy(),
			}
		}

		draw.Draw(bg, image.Rectangle{
			position,
			position.Add(bounds.Size()),
		}, images[i], bounds.Min, draw.Over)
	}
}

// drawPosForeground draw the given images on the given background inside the
// section delimited by the options position.
func drawPosForeground(bg *image.RGBA, images []image.Image, options *Options) {
	dst := positionForeground(bg, options.Position)
	fg := foregroundImage(dst, options.Color)
	fg = drawForeground(fg, images, options)

	draw.Draw(bg, dst, fg, fg.Bounds().Min, draw.Over)
}

// positionForeground creates a mask with the given position.
func positionForeground(bg image.Image, pos string) image.Rectangle {
	ratios := []int{100, 100, 100, 100}
	val := strings.Split(pos, ".")
	for i := range val {
		if i+1 > len(ratios) {
			break
		}
		ratios[i], _ = strconv.Atoi(val[i])
	}
	b := bg.Bounds()
	return image.Rectangle{
		image.Point{b.Dx() * ratios[0], b.Dy() * ratios[1]}.Div(100),
		image.Point{b.Dx() * ratios[2], b.Dy() * ratios[3]}.Div(100),
	}
}

// foregroundImage creates an Image with the given mask and the given color.
func foregroundImage(rec image.Rectangle, c string) *image.RGBA {
	fg := image.NewRGBA(image.Rectangle{image.ZP, rec.Size()})
	if c == "" {
		return fg
	}

	col, err := colorful.Hex(fmt.Sprintf("#%s", c))
	if err != nil {
		return fg
	}

	draw.Draw(fg, fg.Bounds(), &image.Uniform{col}, fg.Bounds().Min, draw.Src)
	return fg
}

// drawForeground draw the given images inside the destination foreground.
// if the foreground image has a height superior to its width, the images
// are vertically aligned, else they are horizontally aligned.
func drawForeground(fg *image.RGBA, images []image.Image, options *Options) *image.RGBA {
	n := len(images)
	if n == 0 {
		return fg
	}

	// resize images for foreground
	b := fg.Bounds()
	opts := &Options{Upscale: true}

	if b.Dx() > b.Dy() {
		opts.Width = b.Dx() / n
		opts.Height = b.Dy()
	} else {
		opts.Width = b.Dx()
		opts.Height = b.Dy() / n
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

// foregroundHorizontal splits the fg according to the number of images  in
// equal parts horizontally aligned and draw each images in the given order in
// the center of each of theses parts.
func foregroundHorizontal(fg *image.RGBA, images []image.Image, options *Options) *image.RGBA {
	position := fg.Bounds().Min
	totalHeight := fg.Bounds().Dy()
	cellWidth := fg.Bounds().Dx() / len(images)
	for i := range images {
		bounds := images[i].Bounds()
		position.Y = (totalHeight - bounds.Dy()) / 2
		position.X = fg.Bounds().Min.X + i*cellWidth + (cellWidth-bounds.Dx())/2
		r := image.Rectangle{
			position,
			position.Add(fg.Bounds().Size()),
		}
		draw.Draw(fg, r, images[i], bounds.Min, draw.Over)
	}
	return fg
}

// foregroundVertical splits the fg according to the number of images  in
// equal parts vertically aligned and draw each images in the given order in
// the center of each of theses parts.
func foregroundVertical(fg *image.RGBA, images []image.Image, options *Options) *image.RGBA {
	position := fg.Bounds().Min
	cellHeight := fg.Bounds().Dy() / len(images)
	totalWidth := fg.Bounds().Dx()
	for i := range images {
		bounds := images[i].Bounds()
		position.Y = fg.Bounds().Min.Y + i*cellHeight + (cellHeight-bounds.Dy())/2
		position.X = fg.Bounds().Min.X + (totalWidth-bounds.Dx())/2
		r := image.Rectangle{
			position,
			position.Add(image.Point{bounds.Dx(), bounds.Dy()}),
		}
		draw.Draw(fg, r, images[i], bounds.Min, draw.Over)
	}
	return fg
}
