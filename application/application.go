package application

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thoas/kvstores"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/storages"
	"log"
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
	SecretKey     string
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
	results := hash.Shard(filename, a.Shard.Width, a.Shard.Depth, true)

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

	err = a.DestStorage.Save(i.Filepath, content)

	if err != nil {
		a.Logger.Error.Print(err)
	} else {
		a.Logger.Info.Printf("Save thumbnail %s to storage", i.Filepath)

		err = con.Set(i.Key, i.Filepath)

		if err != nil {
			a.Logger.Info.Printf("Save key %s=%s to kvstore", i.Key, i.Filepath)
		} else {
			a.Logger.Error.Print(err)
		}
	}

	return err
}

func (a *Application) ImageFileFromRequest(req *Request, async bool, load bool) (*image.ImageFile, error) {
	var file *image.ImageFile = &image.ImageFile{Key: req.Key}
	var err error

	// Image from the KVStore found
	stored := req.Connection.Get(req.Key)

	file.Filepath = stored

	if stored != "" {
		if load {
			file, err = file.LoadFromStorage(a.DestStorage)

			if err != nil {
				return nil, err
			}
		}
	} else {
		// Image not found from the KVStore, we need to process it
		// URL available in Query String
		if req.URL != nil {
			file, err = image.ImageFileFromURL(req.URL)
		} else {
			// URL provided we use http protocol to retrieve it
			file.Filepath = req.Filepath

			file, err = file.LoadFromStorage(a.SourceStorage)
		}

		if err != nil {
			return nil, err
		}

		file, err = file.Transform(req.Operation, req.QueryString)

		if err != nil {
			return nil, err
		}

		format := file.Format()

		if format == "" {
			format = DefaultFormat
		}

		file.Filepath = fmt.Sprintf("%s.%s", a.ShardFilename(req.Key), format)
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

func (a *Application) ToJSON(file *image.ImageFile) ([]byte, error) {
	return json.Marshal(map[string]string{
		"filename": file.GetFilename(),
		"path":     a.DestStorage.Path(file.Filepath),
		"url":      a.DestStorage.URL(file.Filepath),
	})
}

func (a *Application) IsValidSign(qs string) bool {
	if a.SecretKey == "" {
		return true
	}

	return signature.VerifySign(a.SecretKey, qs)
}
