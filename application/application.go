package application

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/kvstores"
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

	filename, err := a.Storage.Save(fmt.Sprintf("%s.%s", a.ShardFilename(i.Key), i.Format()), i.ContentType, content)

	if err != nil {
		a.Logger.Error.Print(err)
	} else {
		err = con.Set(i.Key, filename)

		a.Logger.Info.Printf("Save thumbnail %s to storage", filename)
	}
}
