package application

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/thoas/kvstores"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/storages"
	"net/url"
	"strings"
)

type Shard struct {
	Depth int
	Width int
}

type Logger struct {
	Info  *logrus.Logger
	Error *logrus.Logger
}

type Application struct {
	Prefix        string
	SecretKey     string
	Format        string
	KVStore       kvstores.KVStore
	SourceStorage storages.Storage
	DestStorage   storages.Storage
	Router        *mux.Router
	Shard         Shard
	Logger        Logger
	Raven         *raven.Client
}

func NewApplication() *Application {
	var ErrorLogger = logrus.New()
	ErrorLogger.Level = logrus.ErrorLevel

	return &Application{
		Logger: Logger{
			Info:  logrus.New(),
			Error: ErrorLogger,
		},
	}
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

		return err
	}

	a.Logger.Info.Printf("Save thumbnail %s to storage", i.Filepath)

	key := a.WithPrefix(i.Key)

	err = con.Set(key, i.Filepath)

	if err != nil {
		a.Logger.Error.Print(err)

		return err
	}

	a.Logger.Info.Printf("Save key %s=%s to kvstore", key, i.Filepath)

	return nil
}

func (a *Application) WithPrefix(str string) string {
	return a.Prefix + str
}

func (a *Application) ImageFileFromRequest(req *Request, async bool, load bool) (*image.ImageFile, error) {
	var file *image.ImageFile = &image.ImageFile{
		Key:     req.Key,
		Storage: a.DestStorage,
	}
	var err error

	// Image from the KVStore found
	stored, err := kvstores.String(req.Connection.Get(a.WithPrefix(req.Key)))

	file.Filepath = stored

	if stored != "" {
		if load {
			file, err = image.FromStorage(a.DestStorage, stored)

			if err != nil {
				return nil, err
			}
		}
	} else {
		// Image not found from the KVStore, we need to process it
		// URL available in Query String
		if req.URL != nil {
			file, err = image.FromURL(req.URL)
		} else {
			// URL provided we use http protocol to retrieve it
			file, err = image.FromStorage(a.SourceStorage, req.Filepath)
		}

		if err != nil {
			return nil, err
		}

		file, err = file.Transform(req.Operation, req.QueryString)

		if err != nil {
			return nil, err
		}

		format := req.Format

		if format == "" {
			format = a.Format
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
func (a *Application) IsValidSign(qs map[string]string) bool {
	if a.SecretKey == "" {
		return true
	}

	params := url.Values{}
	for k, v := range qs {
		params.Set(k, v)
	}

	return signature.VerifySign(a.SecretKey, params.Encode())
}
