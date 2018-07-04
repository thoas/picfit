package lilliput

// #cgo CFLAGS: -msse -msse2 -msse3 -msse4.1 -msse4.2 -mavx
// #cgo darwin CFLAGS: -I${SRCDIR}/deps/osx/include
// #cgo linux CFLAGS: -I${SRCDIR}/deps/linux/include
// #cgo CXXFLAGS: -std=c++11
// #cgo darwin CXXFLAGS: -I${SRCDIR}/deps/osx/include
// #cgo linux CXXFLAGS: -I${SRCDIR}/deps/linux/include
// #cgo LDFLAGS:  -lopencv_core -lopencv_imgcodecs -lopencv_imgproc -ljpeg -lpng -lwebp -lippicv -lz
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/osx/lib -L${SRCDIR}/deps/osx/share/OpenCV/3rdparty/lib -framework Accelerate
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/linux/lib -L${SRCDIR}/deps/linux/share/OpenCV/3rdparty/lib
// #include "opencv.hpp"
import "C"

import (
	"io"
	"time"
	"unsafe"
)

// ImageOrientation describes how the decoded image is oriented according to its metadata.
type ImageOrientation int

const (
	JpegQuality    = int(C.CV_IMWRITE_JPEG_QUALITY)
	PngCompression = int(C.CV_IMWRITE_PNG_COMPRESSION)
	WebpQuality    = int(C.CV_IMWRITE_WEBP_QUALITY)

	JpegProgressive = int(C.CV_IMWRITE_JPEG_PROGRESSIVE)

	OrientationTopLeft     = ImageOrientation(C.CV_IMAGE_ORIENTATION_TL)
	OrientationTopRight    = ImageOrientation(C.CV_IMAGE_ORIENTATION_TR)
	OrientationBottomRight = ImageOrientation(C.CV_IMAGE_ORIENTATION_BR)
	OrientationBottomLeft  = ImageOrientation(C.CV_IMAGE_ORIENTATION_BL)
	OrientationLeftTop     = ImageOrientation(C.CV_IMAGE_ORIENTATION_LT)
	OrientationRightTop    = ImageOrientation(C.CV_IMAGE_ORIENTATION_RT)
	OrientationRightBottom = ImageOrientation(C.CV_IMAGE_ORIENTATION_RB)
	OrientationLeftBottom  = ImageOrientation(C.CV_IMAGE_ORIENTATION_LB)
)

// PixelType describes the base pixel type of the image.
type PixelType int

// ImageHeader contains basic decoded image metadata.
type ImageHeader struct {
	width       int
	height      int
	pixelType   PixelType
	orientation ImageOrientation
	numFrames   int
}

// Framebuffer contains an array of raw, decoded pixel data.
type Framebuffer struct {
	buf       []byte
	mat       C.opencv_mat
	width     int
	height    int
	pixelType PixelType
}

type openCVDecoder struct {
	decoder       C.opencv_decoder
	mat           C.opencv_mat
	hasReadHeader bool
	hasDecoded    bool
}

type openCVEncoder struct {
	encoder C.opencv_encoder
	dst     C.opencv_mat
	dstBuf  []byte
}

// Depth returns the number of bits in the PixelType.
func (p PixelType) Depth() int {
	return int(C.opencv_type_depth(C.int(p)))
}

// Channels returns the number of channels in the PixelType.
func (p PixelType) Channels() int {
	return int(C.opencv_type_channels(C.int(p)))
}

// Width returns the width of the image in number of pixels.
func (h *ImageHeader) Width() int {
	return h.width
}

// Height returns the height of the image in number of pixels.
func (h *ImageHeader) Height() int {
	return h.height
}

// PixelType returns a PixelType describing the image's pixels.
func (h *ImageHeader) PixelType() PixelType {
	return h.pixelType
}

// ImageOrientation returns the metadata-based image orientation.
func (h *ImageHeader) Orientation() ImageOrientation {
	return h.orientation
}

// NewFramebuffer creates the backing store for a pixel frame buffer.
func NewFramebuffer(width, height int) *Framebuffer {
	return &Framebuffer{
		buf: make([]byte, width*height*4),
		mat: nil,
	}
}

// Close releases the resources associated with Framebuffer.
func (f *Framebuffer) Close() {
	if f.mat != nil {
		C.opencv_mat_release(f.mat)
		f.mat = nil
	}
}

// Clear resets all of the pixel data in Framebuffer.
func (f *Framebuffer) Clear() {
	C.memset(unsafe.Pointer(&f.buf[0]), 0, C.size_t(len(f.buf)))
}

func (f *Framebuffer) resizeMat(width, height int, pixelType PixelType) error {
	if f.mat != nil {
		C.opencv_mat_release(f.mat)
		f.mat = nil
	}
	if pixelType.Depth() > 8 {
		pixelType = PixelType(C.opencv_type_convert_depth(C.int(pixelType), C.CV_8U))
	}
	newMat := C.opencv_mat_create_from_data(C.int(width), C.int(height), C.int(pixelType), unsafe.Pointer(&f.buf[0]), C.size_t(len(f.buf)))
	if newMat == nil {
		return ErrBufTooSmall
	}
	f.mat = newMat
	f.width = width
	f.height = height
	f.pixelType = pixelType
	return nil
}

// OrientationTransform rotates and/or mirrors the Framebuffer. Passing the
// orientation given by the ImageHeader will normalize the orientation of the Framebuffer.
func (f *Framebuffer) OrientationTransform(orientation ImageOrientation) {
	if f.mat == nil {
		return
	}

	C.opencv_mat_orientation_transform(C.CVImageOrientation(orientation), f.mat)
	f.width = int(C.opencv_mat_get_width(f.mat))
	f.height = int(C.opencv_mat_get_height(f.mat))
}

// ResizeTo performs a resizing transform on the Framebuffer and puts the result
// in the provided destination Framebuffer. This function does not preserve aspect
// ratio if the given dimensions differ in ratio from the source. Returns an error
// if the destination is not large enough to hold the given dimensions.
func (f *Framebuffer) ResizeTo(width, height int, dst *Framebuffer) error {
	err := dst.resizeMat(width, height, f.pixelType)
	if err != nil {
		return err
	}
	C.opencv_mat_resize(f.mat, dst.mat, C.int(width), C.int(height), C.CV_INTER_AREA)
	return nil
}

// Fit performs a resizing and cropping transform on the Framebuffer and puts the result
// in the provided destination Framebuffer. This function does preserve aspect ratio
// but will crop columns or rows from the edges of the image as necessary in order to
// keep from stretching the image content. Returns an error if the destination is
// not large enough to hold the given dimensions.
func (f *Framebuffer) Fit(width, height int, dst *Framebuffer) error {
	if f.mat == nil {
		return ErrFrameBufNoPixels
	}

	aspectIn := float64(f.width) / float64(f.height)
	aspectOut := float64(width) / float64(height)

	var widthPostCrop, heightPostCrop int
	if aspectIn > aspectOut {
		// input is wider than output, so we'll need to narrow
		// we preserve input height and reduce width
		widthPostCrop = int((aspectOut * float64(f.height)) + 0.5)
		heightPostCrop = f.height
	} else {
		// input is taller than output, so we'll need to shrink
		heightPostCrop = int((float64(f.width) / aspectOut) + 0.5)
		widthPostCrop = f.width
	}

	var left, top int
	left = int(float64(f.width-widthPostCrop) * 0.5)
	if left < 0 {
		left = 0
	}

	top = int(float64(f.height-heightPostCrop) * 0.5)
	if top < 0 {
		top = 0
	}

	newMat := C.opencv_mat_crop(f.mat, C.int(left), C.int(top), C.int(widthPostCrop), C.int(heightPostCrop))
	defer C.opencv_mat_release(newMat)

	err := dst.resizeMat(width, height, f.pixelType)
	if err != nil {
		return err
	}
	C.opencv_mat_resize(newMat, dst.mat, C.int(width), C.int(height), C.CV_INTER_AREA)
	return nil
}

// Width returns the width of the contained pixel data in number of pixels. This may
// differ from the capacity of the framebuffer.
func (f *Framebuffer) Width() int {
	return f.width
}

// Height returns the height of the contained pixel data in number of pixels. This may
// differ from the capacity of the framebuffer.
func (f *Framebuffer) Height() int {
	return f.height
}

// PixelType returns the PixelType information of the contained pixel data, if any.
func (f *Framebuffer) PixelType() PixelType {
	return f.pixelType
}

func newOpenCVDecoder(buf []byte) (*openCVDecoder, error) {
	mat := C.opencv_mat_create_from_data(C.int(len(buf)), 1, C.CV_8U, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))

	// this next check is sort of silly since this array is 1-dimensional
	// but if the create ever changes and we goof up, could catch a
	// buffer overwrite
	if mat == nil {
		return nil, ErrBufTooSmall
	}

	decoder := C.opencv_decoder_create(mat)
	if decoder == nil {
		C.opencv_mat_release(mat)
		return nil, ErrInvalidImage
	}

	return &openCVDecoder{
		mat:     mat,
		decoder: decoder,
	}, nil
}

func (d *openCVDecoder) Header() (*ImageHeader, error) {
	if !d.hasReadHeader {
		if !C.opencv_decoder_read_header(d.decoder) {
			return nil, ErrInvalidImage
		}
	}

	d.hasReadHeader = true

	return &ImageHeader{
		width:       int(C.opencv_decoder_get_width(d.decoder)),
		height:      int(C.opencv_decoder_get_height(d.decoder)),
		pixelType:   PixelType(C.opencv_decoder_get_pixel_type(d.decoder)),
		orientation: ImageOrientation(C.opencv_decoder_get_orientation(d.decoder)),
		numFrames:   1,
	}, nil
}

func (d *openCVDecoder) Close() {
	C.opencv_decoder_release(d.decoder)
	C.opencv_mat_release(d.mat)
}

func (d *openCVDecoder) Description() string {
	return C.GoString(C.opencv_decoder_get_description(d.decoder))
}

func (d *openCVDecoder) Duration() time.Duration {
	return time.Duration(0)
}

func (d *openCVDecoder) DecodeTo(f *Framebuffer) error {
	if d.hasDecoded {
		return io.EOF
	}
	h, err := d.Header()
	if err != nil {
		return err
	}
	err = f.resizeMat(h.Width(), h.Height(), h.PixelType())
	if err != nil {
		return err
	}
	ret := C.opencv_decoder_read_data(d.decoder, f.mat)
	if !ret {
		return ErrDecodingFailed
	}
	d.hasDecoded = true
	return nil
}

func newOpenCVEncoder(ext string, decodedBy Decoder, dstBuf []byte) (*openCVEncoder, error) {
	dstBuf = dstBuf[:1]
	dst := C.opencv_mat_create_empty_from_data(C.int(cap(dstBuf)), unsafe.Pointer(&dstBuf[0]))

	if dst == nil {
		return nil, ErrBufTooSmall
	}

	c_ext := C.CString(ext)
	defer C.free(unsafe.Pointer(c_ext))
	enc := C.opencv_encoder_create(c_ext, dst)
	if enc == nil {
		return nil, ErrInvalidImage
	}

	return &openCVEncoder{
		encoder: enc,
		dst:     dst,
		dstBuf:  dstBuf,
	}, nil
}

func (e *openCVEncoder) Encode(f *Framebuffer, opt map[int]int) ([]byte, error) {
	if f == nil {
		return nil, io.EOF
	}
	var optList []C.int
	var firstOpt *C.int
	for k, v := range opt {
		optList = append(optList, C.int(k))
		optList = append(optList, C.int(v))
	}
	if len(optList) > 0 {
		firstOpt = (*C.int)(unsafe.Pointer(&optList[0]))
	}
	if !C.opencv_encoder_write(e.encoder, f.mat, firstOpt, C.size_t(len(optList))) {
		return nil, ErrInvalidImage
	}

	ptrCheck := C.opencv_mat_get_data(e.dst)
	if ptrCheck != unsafe.Pointer(&e.dstBuf[0]) {
		// mat pointer got reallocated - the passed buf was too small to hold the image
		// XXX we should free? the mat here, probably want to recreate
		return nil, ErrBufTooSmall
	}

	length := int(C.opencv_mat_get_height(e.dst))

	return e.dstBuf[:length], nil
}

func (e *openCVEncoder) Close() {
	C.opencv_encoder_release(e.encoder)
	C.opencv_mat_release(e.dst)
}
