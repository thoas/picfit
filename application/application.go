package application

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thoas/kvstores"
	"github.com/thoas/picfit/image"
	"github.com/thoas/storages"
	"log"
	"mime"
	"strings"
)

type Shard struct {
	Depth int
	Width int
}

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

type Application struct {
	Format        string
	ContentType   string
	BaseURL       string
	KVStore       kvstores.KVStore
	SourceStorage storages.Storage
	DestStorage   storages.Storage
	Router        *mux.Router
	Shard         Shard
	Logger        Logger
}

func (a *Application) ShardFilename(filename string) string {
	results := shard(filename, a.Shard.Width, a.Shard.Depth, true)

	return strings.Join(results, "/")
}

func (a *Application) Store(i *image.ImageFile) error {
	con := App.KVStore.Connection()
	defer con.Close()

	content, err := i.ToBytes()

	if err != nil {
		a.Logger.Error.Print(err)
		return err
	}

	err = a.DestStorage.Save(i.Filename, content)

	if err != nil {
		a.Logger.Error.Print(err)
	} else {
		a.Logger.Info.Printf("Save thumbnail %s to storage", i.Filename)

		err = con.Set(i.Key, i.Filename)

		if err != nil {
			a.Logger.Info.Printf("Save key %s=%s to kvstore", i.Key, i.Filename)
		} else {
			a.Logger.Error.Print(err)
		}
	}

	return err
}

func (a *Application) ImageFileFromStorage(filename string) (*image.ImageFile, error) {
	var file *image.ImageFile
	var err error

	// URL provided we use http protocol to retrieve it
	if a.SourceStorage.HasBaseURL() {
		file, err = image.ImageFileFromURL(a.SourceStorage.URL(filename))

		if err != nil {
			return nil, err
		}
	} else {
		body, err := a.SourceStorage.Open(filename)

		if err != nil {
			return nil, err
		}

		modifiedTime, err := a.SourceStorage.ModifiedTime(filename)

		if err != nil {
			return nil, err
		}

		contentType := mime.TypeByExtension(filename)

		headers := map[string]string{
			"Last-Modified": modifiedTime.Format(storages.LastModifiedFormat),
			"Content-Type":  contentType,
		}

		file, err = image.ImageFileFromBytes(body, contentType, headers)

		if err != nil {
			return nil, err
		}

	}

	file.Filename = a.SourceStorage.Path(filename)

	return file, err
}

func (a *Application) ImageFileFromRequest(req *Request, async bool, load bool) (*image.ImageFile, error) {
	var file *image.ImageFile = &image.ImageFile{Key: req.Key}
	var err error

	// Image from the KVStore found
	stored := req.Connection.Get(req.Key)

	file.Filename = stored

	if stored != "" {
		file, err = a.ImageFileFromStorage(stored)
	} else {
		// Image not found from the KVStore, we need to process it
		// URL available in Query String
		if req.URL != nil {
			file, err = image.ImageFileFromURL(req.URL.String())

		} else {
			// URL provided we use http protocol to retrieve it
			file, err = a.ImageFileFromStorage(req.Filename)
		}

		file, err = file.Transform(req.Method, req.QueryString)

		if err != nil {
			return nil, err
		}

		file.Filename = fmt.Sprintf("%s.%s", a.ShardFilename(req.Key), file.Format())
	}

	file.Key = req.Key

	if stored == "" {
		if async == true {
			go a.Store(file)
		} else {
			err = a.Store(file)
		}
	}

	return file, err
}
