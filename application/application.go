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
	Format      string
	ContentType string
	BaseURL     string
	KVStore     kvstores.KVStore
	Storage     storages.Storage
	Router      *mux.Router
	Shard       Shard
	Logger      Logger
}

func (a *Application) ShardFilename(filename string) string {
	results := shard(filename, a.Shard.Width, a.Shard.Depth, true)

	return strings.Join(results, "/")
}

func (a *Application) Store(i *image.ImageResponse) error {
	con := App.KVStore.Connection()
	defer con.Close()

	content, err := i.ToBytes()

	if err != nil {
		a.Logger.Error.Print(err)
		return err
	}

	err = a.Storage.Save(i.Filename, content)

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

func (a *Application) ImageResponseFromStorage(filename string) (*image.ImageResponse, error) {
	var imageResponse *image.ImageResponse
	var err error

	// URL provided we use http protocol to retrieve it
	if a.Storage.HasBaseURL() {
		imageResponse, err = image.ImageResponseFromURL(a.Storage.URL(filename))

		if err != nil {
			return nil, err
		}
	} else {
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

		imageResponse, err = image.ImageResponseFromBytes(body, contentType, headers)

		if err != nil {
			return nil, err
		}

	}

	imageResponse.Filename = a.Storage.Path(filename)

	return imageResponse, err
}

func (a *Application) ImageResponseFromRequest(req *Request, async bool) (*image.ImageResponse, error) {
	var imageResponse *image.ImageResponse
	var err error

	// Image from the KVStore found
	stored := req.Connection.Get(req.Key)

	if stored != "" {
		imageResponse, err = a.ImageResponseFromStorage(stored)
	} else {
		// Image not found from the KVStore, we need to process it
		// URL available in Query String
		if req.URL != nil {
			imageResponse, err = image.ImageResponseFromURL(req.URL.String())

		} else {
			// URL provided we use http protocol to retrieve it
			imageResponse, err = a.ImageResponseFromStorage(req.Filename)
		}

		file := image.NewImageFile(imageResponse.Image)

		dest, err := file.Transform(req.Method, req.QueryString)

		if err != nil {
			return nil, err
		}

		imageResponse.Image = dest
	}

	imageResponse.Key = req.Key
	imageResponse.Filename = fmt.Sprintf("%s.%s", a.ShardFilename(imageResponse.Key), imageResponse.Format())

	if stored == "" {
		if async == true {
			go a.Store(imageResponse)
		} else {
			err = a.Store(imageResponse)
		}
	}

	return imageResponse, err
}
