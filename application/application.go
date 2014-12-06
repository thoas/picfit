package application

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/kvstores"
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
	Format      string
	ContentType string
	BaseURL     string
	KVStore     kvstores.KVStore
	Storage     storages.Storage
	Router      *mux.Router
	Shard       Shard
	Logger      Logger
}

func (a *Application) URL(str ...string) string {
	var results []string

	results = append(results, a.BaseURL)

	for _, value := range str {
		results = append(results, value)
	}

	return strings.Join(results, "/")
}

func (a *Application) ShardFilename(filename string) string {
	results := shard(filename, a.Shard.Width, a.Shard.Depth, true)

	return strings.Join(results, "/")
}

func (a *Application) Store(i *image.ImageResponse) {
	con := App.KVStore.Connection()
	defer con.Close()

	content, err := i.ToBytes()

	if err != nil {
		a.Logger.Error.Print(err)
		return
	}

	filename := fmt.Sprintf("%s.%s", a.ShardFilename(i.Key), i.Format())

	err = a.Storage.Save(filename, content)

	if err != nil {
		a.Logger.Error.Print(err)
	} else {
		a.Logger.Info.Printf("Save thumbnail %s to storage", filename)

		err = con.Set(i.Key, filename)

		if err != nil {
			a.Logger.Info.Printf("Save key %s=%s to kvstore", i.Key, filename)
		} else {
			a.Logger.Error.Print(err)
		}
	}
}

func (a *Application) ImageResponseFromStorage(filename string) (*image.ImageResponse, error) {
	body, err := a.Storage.Open(filename)

	if err != nil {
		return nil, err
	}

	modifiedTime, err := a.Storage.ModifiedTime(filename)

	if err != nil {
		return nil, err
	}

	contentType := mime.TypeByExtension(filename)

	headers := map[string]string{
		"Last-Modified": modifiedTime.Format(storages.LastModifiedFormat),
		"Content-Type":  contentType,
	}

	imageResponse, err := image.ImageResponseFromBytes(body, contentType, headers)

	if err != nil {
		return nil, err
	}

	return imageResponse, nil
}
