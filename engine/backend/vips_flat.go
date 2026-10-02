package backend

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	vips "github.com/cshum/vipsgen/vips816"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"

	"github.com/thoas/picfit/constants"
	imagefile "github.com/thoas/picfit/image"
)

// Flat implements Backend, it draws options.Images on the background
// with the same geometry as GoImage.Flat.
func (b *Vips) Flat(ctx context.Context, dst io.Writer, background *imagefile.ImageFile, options *Options) error {
	if !vipsCanSave(options.Format) {
		return MethodNotImplementedError
	}

	bgData, bg, err := b.load(ctx, background)
	if err != nil {
		return err
	}
	defer bg.Close()

	foregrounds := make([]*vips.Image, 0, len(options.Images))
	defer func() {
		for _, fg := range foregrounds {
			fg.Close()
		}
	}()
	fgData := make([][]byte, 0, len(options.Images))
	for i := range options.Images {
		data, fg, err := b.load(ctx, &options.Images[i])
		if err != nil {
			if errors.Is(err, MethodNotImplementedError) {
				// the next backend reads every stream again
				rewind(background, bgData)
				for j := range fgData {
					rewind(&options.Images[j], fgData[j])
				}
			}
			return err
		}
		foregrounds = append(foregrounds, fg)
		fgData = append(fgData, data)
	}

	if bg.Interpretation() != vips.InterpretationSrgb {
		if err := bg.Colourspace(vips.InterpretationSrgb, nil); err != nil {
			return errors.WithStack(err)
		}
	}
	hasAlpha := bg.HasAlpha()

	if options.Stick != "" {
		err = drawVipsStickForeground(bg, foregrounds, options)
	} else {
		err = drawVipsPosForeground(bg, foregrounds, options)
	}
	if err != nil {
		return err
	}

	// composite always outputs an alpha band, GoImage keeps opaque backgrounds opaque
	if !hasAlpha && bg.HasAlpha() {
		if err := bg.ExtractBand(0, &vips.ExtractBandOptions{N: bg.Bands() - 1}); err != nil {
			return errors.WithStack(err)
		}
	}

	return b.save(dst, bg, options)
}

// drawVipsStickForeground resizes each image to options.Width x options.Height
// and draws it in the options.Stick corner.
func drawVipsStickForeground(bg *vips.Image, foregrounds []*vips.Image, options *Options) error {
	// imaging.Resize returns an empty image without dimensions
	if options.Width == 0 && options.Height == 0 {
		return nil
	}

	for _, fg := range foregrounds {
		width, height := resizeSize(fg.Width(), fg.Height(), options.Width, options.Height)
		if err := resizeVipsImage(fg, width, height); err != nil {
			return err
		}

		var x, y int
		switch options.Stick {
		case constants.TopRight:
			x = bg.Width() - width
		case constants.BottomLeft:
			y = bg.Height() - height
		case constants.BottomRight:
			x, y = bg.Width()-width, bg.Height()-height
		}

		if err := bg.Composite2(fg, vips.BlendModeOver, &vips.Composite2Options{X: x, Y: y}); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// drawVipsPosForeground draws the images inside the options.Position rectangle,
// split in equal cells horizontally when the rectangle is wider than tall, vertically otherwise.
// Each image is fitted and centered in its cell.
func drawVipsPosForeground(bg *vips.Image, foregrounds []*vips.Image, options *Options) error {
	x0, y0, x1, y1 := flatPosition(bg.Width(), bg.Height(), options.Position)
	width, height := x1-x0, y1-y0
	if width < 0 || height < 0 {
		return fmt.Errorf("Invalid flat position %s", options.Position)
	}
	if width == 0 || height == 0 {
		return nil
	}

	if options.Color != "" {
		if col, err := colorful.Hex("#" + options.Color); err == nil {
			rect, err := vipsUniform(width, height, col)
			if err != nil {
				return err
			}
			defer rect.Close()

			if err := bg.Composite2(rect, vips.BlendModeOver, &vips.Composite2Options{X: x0, Y: y0}); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	n := len(foregrounds)
	if n == 0 {
		return nil
	}

	horizontal := width > height
	cellWidth, cellHeight := width, height/n
	if horizontal {
		cellWidth, cellHeight = width/n, height
	}
	// imaging.Fit returns an empty image without dimensions
	if cellWidth == 0 || cellHeight == 0 {
		return nil
	}

	for i, fg := range foregrounds {
		w, h := fitSize(fg.Width(), fg.Height(), cellWidth, cellHeight)
		if err := resizeVipsImage(fg, w, h); err != nil {
			return err
		}

		x, y := x0+(width-w)/2, y0+i*cellHeight+(cellHeight-h)/2
		if horizontal {
			x, y = x0+i*cellWidth+(cellWidth-w)/2, y0+(height-h)/2
		}

		if err := bg.Composite2(fg, vips.BlendModeOver, &vips.Composite2Options{X: x, Y: y}); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// flatPosition parses "x0.y0.x1.y1" percentages of the background like GoImage,
// missing or invalid values default to 100.
func flatPosition(bgWidth, bgHeight int, pos string) (int, int, int, int) {
	ratios := []int{100, 100, 100, 100}
	for i, val := range strings.Split(pos, ".") {
		if i >= len(ratios) {
			break
		}
		// GoImage ignores the error, Atoi returns 0
		ratios[i], _ = strconv.Atoi(val)
	}

	return bgWidth * ratios[0] / 100, bgHeight * ratios[1] / 100,
		bgWidth * ratios[2] / 100, bgHeight * ratios[3] / 100
}

// vipsUniform returns an opaque width x height sRGB image of the given color.
func vipsUniform(width, height int, col colorful.Color) (*vips.Image, error) {
	black, err := vips.NewBlack(width, height, &vips.BlackOptions{Bands: 3})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer black.Close()

	r, g, b := col.RGB255()
	// vipsgen passes len(a) as the length of both arrays, a must be as long as b
	if err := black.Linear([]float64{1, 1, 1}, []float64{float64(r), float64(g), float64(b)}, &vips.LinearOptions{Uchar: true}); err != nil {
		return nil, errors.WithStack(err)
	}

	rect, err := black.Copy(&vips.CopyOptions{Interpretation: vips.InterpretationSrgb})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rect, nil
}

func resizeVipsImage(img *vips.Image, width, height int) error {
	if img.Width() == width && img.Height() == height {
		return nil
	}

	return errors.WithStack(img.ThumbnailImage(width, &vips.ThumbnailImageOptions{
		Height: height,
		Size:   vips.SizeForce,
	}))
}

func rewind(img *imagefile.ImageFile, data []byte) {
	img.Stream = io.NopCloser(bytes.NewReader(data))
}
