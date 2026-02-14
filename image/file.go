package image

import (
	"context"
	"io"
	"mime"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/thoas/picfit/storage"
)

type ImageFile struct {
	Filepath string
	Headers  map[string]string
	Key      string

	Stream io.ReadCloser

	StorageStream io.Reader
	HTTPStream    io.Reader

	Storage *storage.Storage
}

func (i *ImageFile) URL() string {
	return i.Storage.URL(i.Filepath)
}

// Path joins the given file to the storage path
func (i *ImageFile) Path() string {
	return i.Storage.Path(i.Filepath)
}

func (i *ImageFile) Content() io.ReadCloser {
	return i.Stream
}

func (i *ImageFile) Close() {
	i.Stream.Close()
	i.HTTPStream = nil
}

func (i *ImageFile) Save(ctx context.Context) error {
	if err := i.Storage.Save(ctx, i.StorageStream, i.Filepath); err != nil {
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
