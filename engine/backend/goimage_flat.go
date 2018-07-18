package backend

import (
	"image"
	"image/draw"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	colorful "github.com/lucasb-eyer/go-colorful"

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

	dst := positionForeground(bg, options.Position)
	fg := foregroundImage(dst, options.Color)
	fg = drawForeground(fg, images, options)

	draw.Draw(bg, dst, fg, fg.Bounds().Min, draw.Over)

	return e.ToBytes(bg, options.Format, options.Quality)
}

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

func foregroundImage(rec image.Rectangle, c string) *image.RGBA {
	fg := image.NewRGBA(image.Rectangle{image.ZP, rec.Size()})
	if c == "" {
		return fg
	}

	col, err := colorful.Hex(c)
	if err != nil {
		return fg
	}

	draw.Draw(fg, fg.Bounds(), &image.Uniform{col}, fg.Bounds().Min, draw.Src)
	return fg
}

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
