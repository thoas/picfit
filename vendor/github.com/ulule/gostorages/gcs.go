package gostorages

import (
	"context"
	"io/ioutil"
	"mime"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var ctx = context.Background()

func NewGCSStorage(credentialsFile, bucket, location, baseURL, cacheControl string) (Storage, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &GCSStorage{
		NewBaseStorage(location, baseURL),
		client.Bucket(bucket),
		cacheControl,
	}, nil
}

type GCSStorage struct {
	*BaseStorage
	bucket       *storage.BucketHandle
	cacheControl string
}

type GCSStorageFile struct {
	*storage.Reader
}

func (f *GCSStorageFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(f)
}

// Open returns the file content in a dedicated bucket
func (gcs *GCSStorage) Open(filepath string) (File, error) {
	object := gcs.bucket.Object(gcs.Path(filepath))

	body, err := object.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSStorageFile{
		body,
	}, nil
}

// Delete the file from the bucket
func (gcs *GCSStorage) Delete(filepath string) error {
	object := gcs.bucket.Object(gcs.Path(filepath))

	return object.Delete(ctx)
}

// Exists checks if the given file is in the bucket
func (gcs *GCSStorage) Exists(filepath string) bool {
	object := gcs.bucket.Object(gcs.Path(filepath))

	_, err := object.Attrs(ctx)

	if err != nil {
		return false
	}

	return true
}

// IsNotExist returns a boolean indicating whether the error is known
// to report that a file or directory does not exist.
func (gcs *GCSStorage) IsNotExist(err error) bool {
	thr, ok := err.(*googleapi.Error)
	if !ok {
		return false
	}
	return thr.Code == 404
}

// ModifiedTime returns the last update time
func (gcs *GCSStorage) ModifiedTime(filepath string) (time.Time, error) {
	object := gcs.bucket.Object(gcs.Path(filepath))

	body, err := object.NewReader(ctx)

	if err != nil {
		return time.Time{}, err
	}

	return body.LastModified()
}

// Save saves a file at the given path in the bucket
func (gcs *GCSStorage) SaveWithContentType(filepath string, file File, contentType string) error {
	object := gcs.bucket.Object(gcs.Path(filepath))

	content, err := file.ReadAll()

	if err != nil {
		return err
	}

	w := object.NewWriter(ctx)

	w.ContentType = contentType
	w.CacheControl = gcs.cacheControl
	_, err = w.Write(content)
	if err != nil {
		return err
	}
	return w.Close()
}

// Save saves a file at the given path in the bucket
func (gcs *GCSStorage) Save(path string, file File) error {
	return gcs.SaveWithContentType(path, file, mime.TypeByExtension(filepath.Ext(path)))
}

// Size returns the size of the given file
func (gcs *GCSStorage) Size(filepath string) int64 {
	object := gcs.bucket.Object(gcs.Path(filepath))

	attrs, err := object.Attrs(ctx)

	if err != nil {
		return 0
	}

	return attrs.Size
}
