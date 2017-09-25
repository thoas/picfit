package fastimage

import (
	"bytes"
	"errors"
)

type imageGIF struct{}

func (g imageGIF) Type() ImageType {
	return GIF
}

func (g imageGIF) Detect(buffer []byte) bool {
	firstTwoBytes := buffer[:2]
	return bytes.Equal(firstTwoBytes, []byte{0x47, 0x49})
}

func (g imageGIF) GetSize(buffer []byte) (*ImageSize, error) {
	if len(buffer) <= 11 {
		return nil, errors.New("Insufficient data")
	}

	imageSize := ImageSize{}
	slice := buffer[6 : 6+4]

	imageSize.Width = uint32(readULint16(slice[:2]))
	imageSize.Height = uint32(readULint16(slice[2:4]))

	return &imageSize, nil
}

func init() {
	register(&imageGIF{})
}
