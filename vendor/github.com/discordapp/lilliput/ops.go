package lilliput

import (
	"io"
)

type ImageOpsSizeMethod int

const (
	ImageOpsNoResize ImageOpsSizeMethod = iota
	ImageOpsFit
	ImageOpsResize
)

// ImageOptions controls how ImageOps resizes and encodes the
// pixel data decoded from a Decoder
type ImageOptions struct {
	// FileType should be a string starting with '.', e.g.
	// ".jpeg"
	FileType string

	// Width controls the width of the output image
	Width int

	// Height controls the height of the output image
	Height int

	// ResizeMethod controls how the image will be transformed to
	// its output size. Notably, ImageOpsFit will do a cropping
	// resize, while ImageOpsResize will stretch the image.
	ResizeMethod ImageOpsSizeMethod

	// NormalizeOrientation will flip and rotate the image as necessary
	// in order to undo EXIF-based orientation
	NormalizeOrientation bool

	// EncodeOptions controls the encode quality options
	EncodeOptions map[int]int
}

// ImageOps is a reusable object that can resize and encode images.
type ImageOps struct {
	frames     []*Framebuffer
	frameIndex int
}

// NewImageOps creates a new ImageOps object that will operate
// on images up to maxSize on each axis.
func NewImageOps(maxSize int) *ImageOps {
	frames := make([]*Framebuffer, 2)
	frames[0] = NewFramebuffer(maxSize, maxSize)
	frames[1] = NewFramebuffer(maxSize, maxSize)
	return &ImageOps{
		frames:     frames,
		frameIndex: 0,
	}
}

func (o *ImageOps) active() *Framebuffer {
	return o.frames[o.frameIndex]
}

func (o *ImageOps) secondary() *Framebuffer {
	return o.frames[1-o.frameIndex]
}

func (o *ImageOps) swap() {
	o.frameIndex = 1 - o.frameIndex
}

// Clear resets all pixel data in ImageOps. This need not be called
// between calls to Transform. You may choose to call this to remove
// image data from memory.
func (o *ImageOps) Clear() {
	o.frames[0].Clear()
	o.frames[1].Clear()
}

// Close releases resources associated with ImageOps
func (o *ImageOps) Close() {
	o.frames[0].Close()
	o.frames[1].Close()
}

func (o *ImageOps) decode(d Decoder) error {
	active := o.active()
	return d.DecodeTo(active)
}

func (o *ImageOps) fit(d Decoder, width, height int) error {
	active := o.active()
	secondary := o.secondary()
	err := active.Fit(width, height, secondary)
	if err != nil {
		return err
	}
	o.swap()
	return nil
}

func (o *ImageOps) resize(d Decoder, width, height int) error {
	active := o.active()
	secondary := o.secondary()
	err := active.ResizeTo(width, height, secondary)
	if err != nil {
		return err
	}
	o.swap()
	return nil
}

func (o *ImageOps) normalizeOrientation(orientation ImageOrientation) {
	active := o.active()
	active.OrientationTransform(orientation)
}

func (o *ImageOps) encode(e Encoder, opt map[int]int) ([]byte, error) {
	active := o.active()
	return e.Encode(active, opt)
}

func (o *ImageOps) encodeEmpty(e Encoder, opt map[int]int) ([]byte, error) {
	return e.Encode(nil, opt)
}

// Transform performs the requested transform operations on the Decoder specified by d.
// The result is written into the output buffer dst. A new slice pointing to dst is returned
// with its length set to the length of the resulting image. Errors may occur if the decoded
// image is too large for ImageOps or if Encoding fails.
//
// It is important that .Decode() not have been called already on d.
func (o *ImageOps) Transform(d Decoder, opt *ImageOptions, dst []byte) ([]byte, error) {
	h, err := d.Header()
	if err != nil {
		return nil, err
	}

	enc, err := NewEncoder(opt.FileType, d, dst)
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	for {
		err = o.decode(d)
		emptyFrame := false
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			// io.EOF means we are out of frames, so we should signal to encoder to wrap up
			emptyFrame = true
		}

		o.normalizeOrientation(h.Orientation())

		if opt.ResizeMethod == ImageOpsFit {
			o.fit(d, opt.Width, opt.Height)
		} else if opt.ResizeMethod == ImageOpsResize {
			o.resize(d, opt.Width, opt.Height)
		}

		var content []byte
		if emptyFrame {
			content, err = o.encodeEmpty(enc, opt.EncodeOptions)
		} else {
			content, err = o.encode(enc, opt.EncodeOptions)
		}

		if err != nil {
			return nil, err
		}

		if content != nil {
			return content, nil
		}

		// content == nil and err == nil -- this is encoder telling us to do another frame

		// for mulitple frames/gifs we need the decoded frame to be active again
		o.swap()
	}
}
