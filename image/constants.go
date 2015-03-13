package image

var ContentTypes = map[string]string{
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"bmp":  "image/bmp",
	"gif":  "image/gif",
}

var Extensions = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/bmp":  "bmp",
	"image/gif":  "gif",
}

var HeaderKeys = []string{
	"Age",
	"Content-Type",
	"Last-Modified",
	"Date",
	"Etag",
}
