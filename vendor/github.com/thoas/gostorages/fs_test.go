package gostorages

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func Test(t *testing.T) {
	baseURL := "http://example.com"

	tmp := os.TempDir()

	file, _ := ioutil.TempFile(tmp, "prefix")

	storage := NewFileSystemStorage(tmp, baseURL)

	filename := path.Base(file.Name())

	assert.True(t, storage.Exists(filename))

	storageFile, err := storage.Open(filename)

	if err != nil {
		t.Fatal(err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fileInfo.Size(), storageFile.Size())

	err = storage.Save("test", NewContentFile([]byte("a content example")))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, storage.URL("test"), baseURL+"/test")

	modified, err := storage.ModifiedTime("test")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modified.Format(time.RFC822), time.Now().Format(time.RFC822))

	err = storage.Delete("test")

	if err != nil {
		t.Fatal(err)
	}
}
