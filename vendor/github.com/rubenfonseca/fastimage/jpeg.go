package fastimage

import (
	"bytes"
	"errors"
)

type jpegHeaderSegment int

const (
	nextSegment jpegHeaderSegment = iota
	sofSegment
	skipSegment
	parseSegment
	eioSegment
)

type imageJPEG struct{}

func (j imageJPEG) Type() ImageType {
	return JPEG
}

func (j imageJPEG) Detect(buffer []byte) bool {
	firstTwoBytes := buffer[:2]
	return bytes.Equal(firstTwoBytes, []byte{0xFF, 0xD8})
}

func (j imageJPEG) GetSize(buffer []byte) (*ImageSize, error) {
	if len(buffer) <= 2 {
		return nil, errors.New("Insufficient data")
	}

	return parseJPEGData(buffer, 2, nextSegment)
}

func parseJPEGData(buffer []byte, offset int, segment jpegHeaderSegment) (*ImageSize, error) {
	logger.Printf("parseJPEGData: buffer size: %v offset: %v segment %v", len(buffer), offset, segment)
	if segment == eioSegment ||
		(len(buffer) <= offset+1) ||
		((len(buffer) <= offset+2) && segment == skipSegment) ||
		((len(buffer) <= offset+7) && segment == parseSegment) {
		logger.Printf("BAILING NOT ENOUGHTDATA")
		return nil, errors.New("Not enough data")
	}

	switch segment {
	case nextSegment:
		newOffset := offset + 1
		b := buffer[newOffset]
		if b == 0xFF {
			return parseJPEGData(buffer, newOffset, sofSegment)
		}
		return parseJPEGData(buffer, newOffset, nextSegment)
	case sofSegment:
		newOffset := offset + 1
		b := buffer[newOffset]

		if b >= 0xE0 && b <= 0xEF {
			return parseJPEGData(buffer, newOffset, skipSegment)
		}
		if (b >= 0xC0 && b <= 0xC3) || (b >= 0xC5 && b <= 0xC7) || (b >= 0xC9 && b <= 0xCB) || b >= 0xCD && b <= 0xCF {
			return parseJPEGData(buffer, newOffset, parseSegment)
		}
		if b == 0xFF {
			return parseJPEGData(buffer, newOffset, sofSegment)
		}
		if b == 0xD9 {
			return parseJPEGData(buffer, newOffset, eioSegment)
		}
		return parseJPEGData(buffer, newOffset, skipSegment)
	case skipSegment:
		length := readUint16(buffer[offset+1 : offset+3])

		newOffset := offset + int(length)
		return parseJPEGData(buffer, newOffset, nextSegment)
	case parseSegment:
		width := readUint16(buffer[offset+6 : offset+8])
		height := readUint16(buffer[offset+4 : offset+6])

		return &ImageSize{Width: uint32(width), Height: uint32(height)}, nil
	default:
		return nil, errors.New("Can't detect jpeg segment.")
	}
}

func init() {
	register(&imageJPEG{})
}
