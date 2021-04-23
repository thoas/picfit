package image

var (
	Extensions = map[string]string{
		"image/bmp":  "bmp",
		"image/gif":  "gif",
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/webp": "webp",
	}

	HeaderKeys = []string{
		"Age",
		"Content-Type",
		"Date",
		"Etag",
		"Last-Modified",
	}
)
