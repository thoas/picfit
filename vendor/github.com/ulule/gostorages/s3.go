package gostorages

import (
	"io"
	"io/ioutil"
	"mime"
	"path/filepath"
	"time"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var ACLs = map[string]s3.ACL{
	"private":                   s3.Private,
	"public-read":               s3.PublicRead,
	"public-read-write":         s3.PublicReadWrite,
	"authenticated-read":        s3.AuthenticatedRead,
	"bucket-owner-read":         s3.BucketOwnerRead,
	"bucket-owner-full-control": s3.BucketOwnerFull,
}

const LastModifiedFormat = time.RFC1123

func NewS3Storage(accessKeyId string, secretAccessKey string, bucketName string, location string, region aws.Region, acl s3.ACL, baseURL string) Storage {

	return &S3Storage{
		NewBaseStorage(location, baseURL),
		accessKeyId,
		secretAccessKey,
		bucketName,
		region,
		acl,
	}
}

type S3Storage struct {
	*BaseStorage
	AccessKeyId     string
	SecretAccessKey string
	BucketName      string
	Region          aws.Region
	ACL             s3.ACL
}

type S3StorageFile struct {
	io.ReadCloser
	Key     *s3.Key
	Storage Storage
}

func (f *S3StorageFile) Size() int64 {
	return f.Key.Size
}

func (f *S3StorageFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(f)
}

// Auth returns a Auth instance
func (s *S3Storage) Auth() (auth aws.Auth, err error) {
	return aws.GetAuth(s.AccessKeyId, s.SecretAccessKey)
}

// Client returns a S3 instance
func (s *S3Storage) Client() (*s3.S3, error) {
	auth, err := s.Auth()

	if err != nil {
		return nil, err
	}

	return s3.New(auth, s.Region), nil
}

// Bucket returns a bucket instance
func (s *S3Storage) Bucket() (*s3.Bucket, error) {
	client, err := s.Client()

	if err != nil {
		return nil, err
	}

	return client.Bucket(s.BucketName), nil
}

// Open returns the file content in a dedicated bucket
func (s *S3Storage) Open(filepath string) (File, error) {
	bucket, err := s.Bucket()

	if err != nil {
		return nil, err
	}

	key, err := s.Key(filepath)

	if err != nil {
		return nil, err
	}

	body, err := bucket.GetReader(s.Path(filepath))

	if err != nil {
		return nil, err
	}

	return &S3StorageFile{
		body,
		key,
		s,
	}, nil
}

// Delete the file from the bucket
func (s *S3Storage) Delete(filepath string) error {
	bucket, err := s.Bucket()

	if err != nil {
		return err
	}

	return bucket.Del(s.Path(filepath))
}

func (s *S3Storage) Key(filepath string) (*s3.Key, error) {
	bucket, err := s.Bucket()

	if err != nil {
		return nil, err
	}

	key, err := bucket.GetKey(s.Path(filepath))

	if err != nil {
		return nil, err
	}

	return key, nil
}

// Exists checks if the given file is in the bucket
func (s *S3Storage) Exists(filepath string) bool {
	_, err := s.Key(filepath)

	if err != nil {
		return false
	}

	return true
}

// Save saves a file at the given path in the bucket
func (s *S3Storage) SaveWithContentType(filepath string, file File, contentType string) error {
	bucket, err := s.Bucket()

	if err != nil {
		return err
	}

	content, err := file.ReadAll()

	if err != nil {
		return err
	}

	err = bucket.Put(s.Path(filepath), content, contentType, s.ACL)

	return err
}

// Save saves a file at the given path in the bucket
func (s *S3Storage) Save(path string, file File) error {
	return s.SaveWithContentType(path, file, mime.TypeByExtension(filepath.Ext(path)))
}

// Size returns the size of the given file
func (s *S3Storage) Size(filepath string) int64 {
	key, err := s.Key(filepath)

	if err != nil {
		return 0
	}

	return key.Size
}

// ModifiedTime returns the last update time
func (s *S3Storage) ModifiedTime(filepath string) (time.Time, error) {
	key, err := s.Key(filepath)

	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(LastModifiedFormat, key.LastModified)
}

// IsNotExist returns a boolean indicating whether the error is known
// to report that a file or directory does not exist.
func (s *S3Storage) IsNotExist(err error) bool {
	thr, ok := err.(*s3.Error)
	if !ok {
		return false
	}
	return thr.StatusCode == 404
}
