package fastimage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type imagePNG struct{}

func (p imagePNG) Type() ImageType {
	return PNG
}

func (p imagePNG) Detect(buffer []byte) bool {
	firstTwoBytes := buffer[:2]
	return bytes.Equal(firstTwoBytes, []byte{0x89, 0x50})
}

func (p imagePNG) GetSize(buffer []byte) (*ImageSize, error) {
	if len(buffer) < 25 {
		return nil, errors.New("Insufficient data")
	}

	imageSize := ImageSize{}
	slice := buffer[16 : 16+8]

	widthBuffer := bytes.NewReader(slice[:4])
	binary.Read(widthBuffer, binary.BigEndian, &imageSize.Width)

	heightBuffer := bytes.NewReader(slice[4:8])
	binary.Read(heightBuffer, binary.BigEndian, &imageSize.Height)

	return &imageSize, nil
}

func init() {
	register(&imagePNG{})
}
