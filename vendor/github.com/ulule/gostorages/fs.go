package gostorages

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/djherbis/times"
)

// Storage is a file system storage handler
type FileSystemStorage struct {
	*BaseStorage
}

type FileSystemFile struct {
	*os.File
	Storage  Storage
	FileInfo os.FileInfo
}

// NewStorage returns a file system storage engine
func NewFileSystemStorage(location string, baseURL string) Storage {
	return &FileSystemStorage{
		&BaseStorage{
			Location: location,
			BaseURL:  baseURL,
		},
	}
}

func NewFileSystemFile(storage Storage, file *os.File) (*FileSystemFile, error) {
	fileInfo, err := file.Stat()

	if err != nil {
		return nil, err
	}

	return &FileSystemFile{
		file,
		storage,
		fileInfo,
	}, nil
}

func (s *FileSystemStorage) URL(filename string) string {
	if s.HasBaseURL() {
		return strings.Join([]string{s.BaseURL, filename}, "/")
	}

	return ""
}

func (f *FileSystemFile) Size() int64 {
	return f.FileInfo.Size()
}

func (f *FileSystemFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(f)
}

// Save saves a file at the given path
func (s *FileSystemStorage) Save(filepath string, file File) error {
	return s.SaveWithPermissions(filepath, file, DefaultFilePermissions)
}

// SaveWithPermissions saves a file with the given permissions to the storage
func (s *FileSystemStorage) SaveWithPermissions(filepath string, file File, perm os.FileMode) error {
	_, err := os.Stat(s.Location)

	if err != nil {
		return err
	}

	location := s.Path(filepath)

	basename := location[:strings.LastIndex(location, "/")+1]

	err = os.MkdirAll(basename, perm)

	if err != nil {
		return err
	}

	content, err := file.ReadAll()

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(location, content, perm)

	return err
}

// Open returns the file content
func (s *FileSystemStorage) Open(filepath string) (File, error) {
	file, err := os.Open(s.Path(filepath))

	if err != nil {
		return nil, err
	}

	return NewFileSystemFile(s, file)
}

// Delete the file from storage
func (s *FileSystemStorage) Delete(filepath string) error {
	return os.Remove(s.Path(filepath))
}

// ModifiedTime returns the last update time
func (s *FileSystemStorage) ModifiedTime(filepath string) (time.Time, error) {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}

	return fi.ModTime(), nil
}

// Exists checks if the given file is in the storage
func (s *FileSystemStorage) Exists(filepath string) bool {
	_, err := os.Stat(s.Path(filepath))
	return err == nil
}

func (s *FileSystemStorage) AccessedTime(filepath string) (time.Time, error) {
	t, err := times.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}

	return t.AccessTime(), nil
}

// CreatedTime returns the last access time.
func (s *FileSystemStorage) CreatedTime(filepath string) (time.Time, error) {
	t, err := times.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}

	return t.ChangeTime(), nil
}

// Size returns the size of the given file
func (s *FileSystemStorage) Size(filepath string) int64 {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return 0
	}
	return fi.Size()
}

// IsNotExist returns a boolean indicating whether the error is known
// to report that a file or directory does not exist.
func (s *FileSystemStorage) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
