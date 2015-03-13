package image

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/engines"
	"math"
	"mime"
	"path"
	"strconv"
	"strings"
)

type ImageFile struct {
	Source    []byte
	Processed []byte
	Key       string
	Headers   map[string]string
	Filepath  string
	Storage   gostorages.Storage
}

func (i *ImageFile) Content() []byte {
	if i.Processed != nil {
		return i.Processed
	}

	return i.Source
}

func (i *ImageFile) Transform(engine engines.Engine, operation *Operation, qs map[string]string, defaultOptions *engines.Options) (*ImageFile, error) {

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

	q, ok := qs["q"]

	var quality int
	var format string

	if ok {
		quality, err := strconv.Atoi(q)

		if err != nil {
			return nil, err
		}

		if quality > 100 {
			return nil, fmt.Errorf("Quality should be <= 100")
		}
	}

	format, ok = qs["fmt"]

	if ok {
		if _, ok := ContentTypes[format]; !ok {
			return nil, fmt.Errorf("Unknown format %s", format)
		}
	} else {
		format = defaultOptions.Format
	}

	file := &ImageFile{
		Source:   i.Source,
		Key:      i.Key,
		Headers:  i.Headers,
		Filepath: i.Filepath,
	}

	options := &engines.Options{
		Quality: quality,
		Format:  format,
		Upscale: upscale,
	}

	switch operation {
	case Resize:
		content, err := engine.Resize(i.Source, w, h, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	case Thumbnail:
		content, err := engine.Thumbnail(i.Source, w, h, options)

		if err != nil {
			return nil, err
		}

		file.Processed = content

		return file, err
	}

	return nil, fmt.Errorf("Operation not found for %s", operation)
}

func (i *ImageFile) URL() string {
	return i.Storage.URL(i.Filepath)
}

func (i *ImageFile) Path() string {
	return i.Storage.Path(i.Filepath)
}

func (i *ImageFile) Save() error {
	return i.Storage.Save(i.Filepath, gostorages.NewContentFile(i.Content()))
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
