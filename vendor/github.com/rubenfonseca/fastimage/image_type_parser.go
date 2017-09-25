package fastimage

// ImageTypeParser is the interface each image type needs to implement to be
// able to detect images on a buffer stream.
type ImageTypeParser interface {
	// Type returns the type of the image
	Type() ImageType

	// Returns true if the image type can be detected on the byte slice
	Detect([]byte) bool

	// Returns the image size by inspecting the byte slice, or error if it
	// can't be detected (more data needed?)
	GetSize([]byte) (*ImageSize, error)
}

var imageTypeParsers = make(map[string]ImageTypeParser)

func register(imageTypeParser ImageTypeParser) {
	if imageTypeParser == nil {
		panic("ImageTypeParser: Register image type is nil")
	}

	name := imageTypeParser.Type().String()
	if _, dup := imageTypeParsers[name]; dup {
		panic("ImageTypeParser: Register called twice for type " + name)
	}

	logger.Printf("ImageTypeParser: Registering %s %v", name, imageTypeParser)
	imageTypeParsers[name] = imageTypeParser
}
