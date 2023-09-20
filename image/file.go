package image

import (
	"bytes"
	"context"
	"mime"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/thoas/picfit/storage"
)

type ImageFile struct {
	Filepath  string
	Headers   map[string]string
	Key       string
	Processed []byte
	Source    []byte
	Storage   storage.Storage
}

func (i *ImageFile) URL() string {
	return i.Storage.URL(i.Filepath)
}

// Path joins the given file to the storage path
func (i *ImageFile) Path() string {
	return i.Storage.Path(i.Filepath)
}

func (i *ImageFile) Content() []byte {
	if i.Processed != nil {
		return i.Processed
	}

	return i.Source
}

func (i *ImageFile) Save(ctx context.Context) error {
	content := bytes.NewReader(i.Content())
	if err := i.Storage.Save(ctx, content, i.Filepath); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *ImageFile) Format() string {
	return Extensions[i.ContentType()]
}

func (i *ImageFile) ContentType() string {
	if _, ok := i.Headers["Content-Type"]; ok {
		return i.Headers["Content-Type"]
	}
	return mime.TypeByExtension(i.FilenameExt())
}

func (i *ImageFile) Filename() string {
	return i.Filepath[strings.LastIndex(i.Filepath, "/")+1:]
}

func (i *ImageFile) FilenameExt() string {
	return path.Ext(i.Filename())
}
