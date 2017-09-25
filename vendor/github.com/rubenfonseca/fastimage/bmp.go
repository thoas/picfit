package fastimage

type imageBMP struct{}

func (b imageBMP) Type() ImageType {
	return BMP
}

func (b imageBMP) Detect(buffer []byte) bool {
	firstTwoBytes := buffer[:2]
	return string(firstTwoBytes) == "BM"
}

func (b imageBMP) GetSize(buffer []byte) (*ImageSize, error) {
	// TODO: We currently don't detect BMP size, so just return nothing
	return nil, nil
}

func init() {
	register(&imageBMP{})
}
