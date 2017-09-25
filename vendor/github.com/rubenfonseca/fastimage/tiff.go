package fastimage

type imageTIFF struct{}

func (t imageTIFF) Type() ImageType {
	return TIFF
}

func (t imageTIFF) Detect(buffer []byte) bool {
	firstTwoBytes := string(buffer[:2])
	return firstTwoBytes == "II" || firstTwoBytes == "MM"
}

func (t imageTIFF) GetSize(buffer []byte) (*ImageSize, error) {
	// TODO: We currently don't detect TIFF size, so just return nothing
	return nil, nil
}

func init() {
	register(&imageTIFF{})
}
