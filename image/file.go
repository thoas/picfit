package image

import (
	"bytes"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/thoas/gostorages"
	"github.com/thoas/imaging"
	"image"
	"math"
	"mime"
	"path"
	"strconv"
	"strings"
)

type ImageFile struct {
	Source   []byte
	Image    image.Image
	Key      string
	Headers  map[string]string
	Filepath string
	Storage  gostorages.Storage
}

func (i *ImageFile) ImageSize() (int, int) {
	return i.Image.Bounds().Max.X, i.Image.Bounds().Max.Y
}

func (i *ImageFile) Scale(dstWidth int, dstHeight int, upscale bool, trans Transformation) *image.NRGBA {
	width, height := i.ImageSize()

	factor := scalingFactor(width, height, dstWidth, dstHeight)

	if factor < 1 || upscale {
		return trans(i.Image, dstWidth, dstHeight, imaging.Lanczos)
	}

	return imaging.Clone(i.Image)
}

func (i *ImageFile) Transform(operation *Operation, qs map[string]string) (*ImageFile, error) {

	if i.Image == nil {
		image, err := imaging.Decode(bytes.NewReader(i.Source))

		if err != nil {
			return nil, err
		}

		i.Image = image
	}

	params := map[string]string{
		"upscale": "1",
		"h":       "0",
		"w":       "0",
	}

	err := mergo.Merge(&qs, params)

	if err != nil {
		return nil, err
	}

	upscale, err := strconv.ParseBool(qs["upscale"])

	if err != nil {
		return nil, err
	}

	w, err := strconv.Atoi(qs["w"])

	if err != nil {
		return nil, err
	}

	h, err := strconv.Atoi(qs["h"])

	if err != nil {
		return nil, err
	}

	switch operation {
	case Resize, Thumbnail:
		dest := i.Scale(w, h, upscale, operation.Transformation)

		file := &ImageFile{
			Image:    dest,
			Source:   i.Source,
			Key:      i.Key,
			Headers:  i.Headers,
			Filepath: i.Filepath,
		}

		return file, err
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (i *ImageFile) ToBytes() ([]byte, error) {
	format, ok := Formats[i.ContentType()]

	if !ok {
		format = DefaultFormat
	}

	return i.ToBytesWithFormat(format)
}

func (i *ImageFile) URL() string {
	return i.Storage.URL(i.Filepath)
}

func (i *ImageFile) Path() string {
	return i.Storage.Path(i.Filepath)
}

func (i *ImageFile) Save() error {
	content, err := i.ToBytes()

	if err != nil {
		return err
	}

	return i.Storage.Save(i.Filepath, gostorages.NewContentFile(content))
}

func (i *ImageFile) SaveWithFormat(format imaging.Format) error {
	content, err := i.ToBytesWithFormat(format)

	if err != nil {
		return err
	}

	return i.Storage.Save(i.Filepath, gostorages.NewContentFile(content))
}

func (i *ImageFile) ToBytesWithFormat(format imaging.Format) ([]byte, error) {
	if i.Image != nil {
		buf := &bytes.Buffer{}

		err := imaging.Encode(buf, i.Image, format)

		if err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	return i.Source, nil
}

func (i *ImageFile) Format() string {
	return Extensions[i.ContentType()]
}

func (i *ImageFile) ContentType() string {
	return mime.TypeByExtension(i.FilenameExt())
}

func (i *ImageFile) Filename() string {
	return i.Filepath[strings.LastIndex(i.Filepath, "/")+1:]
}

func (i *ImageFile) FilenameExt() string {
	return path.Ext(i.Filename())
}

func scalingFactor(srcWidth int, srcHeight int, destWidth int, destHeight int) float64 {
	return math.Max(float64(destWidth)/float64(srcWidth), float64(destHeight)/float64(srcHeight))
}
