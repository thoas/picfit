package image

import (
	"mime"
	"path"
	"strings"

	"github.com/thoas/gostorages"
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
