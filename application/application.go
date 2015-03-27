package application

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/jsonq"
	"github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	"github.com/thoas/gokvstores"
	"github.com/thoas/gostorages"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/engines"
	"github.com/thoas/picfit/extractors"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/middleware"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/picfit/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var Extractors = map[string]extractors.Extractor{
	"op":   extractors.Operation,
	"url":  extractors.URL,
	"path": extractors.Path,
}

type Shard struct {
	Depth int
	Width int
}

type Application struct {
	Prefix        string
	SecretKey     string
	KVStore       gokvstores.KVStore
	SourceStorage gostorages.Storage
	DestStorage   gostorages.Storage
	Router        *mux.Router
	Shard         Shard
	Raven         *raven.Client
	Logger        *logrus.Logger
	Engine        engines.Engine
	Jq            *jsonq.JsonQuery
}

func NewApplication() *Application {
	return &Application{
		Logger: logrus.New(),
	}
}

func NewFromConfigPath(path string) (*Application, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Your config file %s cannot be loaded: %s", path, err)
	}

	return NewFromConfig(content)
}

func NewFromConfig(content []byte) (*Application, error) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	err := dec.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("Your config file %s cannot be parsed: %s", string(content), err)
	}

	jq := jsonq.NewQuery(data)

	return NewFromJsonQuery(jq)
}

func NewFromJsonQuery(jq *jsonq.JsonQuery) (*Application, error) {
	app := NewApplication()
	app.Jq = jq

	for _, initializer := range Initializers {
		err := initializer(app.Jq, app)

		if err != nil {
			return nil, fmt.Errorf("An error occured during init: %s", err)
		}
	}

	return app, nil
}

func (app *Application) ServeHTTP(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		con := app.KVStore.Connection()
		defer con.Close()

		request := muxer.NewRequest(req)

		for k, v := range request.Params {
			request.QueryString[k] = v
		}

		res := muxer.NewResponse(w)

		extracted := map[string]interface{}{}

		for key, extractor := range Extractors {
			result, err := extractor(key, request)

			if err != nil {
				app.Logger.Info(err)

				res.BadRequest()
				return
			}

			extracted[key] = result
		}

		sorted := util.SortMapString(request.QueryString)

		valid := app.IsValidSign(sorted)

		delete(sorted, "sig")

		serialized := hash.Serialize(sorted)

		key := hash.Tokey(serialized)

		app.Logger.Infof("Generating key %s from request: %s", key, serialized)

		var u *url.URL
		var path string

		value, ok := extracted["url"]

		if ok && value != nil {
			u = value.(*url.URL)
		}

		value, ok = extracted["path"]

		if ok && value != nil {
			path = value.(string)
		}

		if !valid || (u == nil && path == "") {
			res.BadRequest()
			return
		}

		h(res, &Request{
			request,
			extracted["op"].(*engines.Operation),
			con,
			key,
			u,
			path,
		}, app)
	})
}

func (a *Application) InitRouter() *negroni.Negroni {
	a.Router = mux.NewRouter()
	a.Router.NotFoundHandler = NotFoundHandler()

	methods := map[string]Handler{
		"redirect": RedirectHandler,
		"display":  ImageHandler,
		"get":      GetHandler,
	}

	for name, handler := range methods {
		handlerFunc := a.ServeHTTP(handler)

		a.Router.Handle(fmt.Sprintf("/%s", name), handlerFunc)
		a.Router.Handle(fmt.Sprintf("/%s/{sig}/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handlerFunc)
		a.Router.Handle(fmt.Sprintf("/%s/{sig}/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", name), handlerFunc)
		a.Router.Handle(fmt.Sprintf("/%s/{sig}/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handlerFunc)
		a.Router.Handle(fmt.Sprintf("/%s/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handlerFunc)
		a.Router.Handle(fmt.Sprintf("/%s/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", name), handlerFunc)
		a.Router.Handle(fmt.Sprintf("/%s/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handlerFunc)
	}

	allowedOrigins, err := a.Jq.ArrayOfStrings("allowed_origins")
	allowedMethods, err := a.Jq.ArrayOfStrings("allowed_methods")

	debug, err := a.Jq.Bool("debug")

	if err != nil {
		debug = false
	}

	n := negroni.New(&middleware.Recovery{
		Raven:      a.Raven,
		Logger:     a.Logger,
		PrintStack: debug,
		StackAll:   false,
		StackSize:  1024 * 8,
	}, &middleware.Logger{a.Logger})
	n.Use(cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
	}))
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(a.Router)

	return n
}

func (a *Application) Port() int {
	port, err := a.Jq.Int("port")

	if err != nil {
		port = DefaultPort
	}

	return port
}

func (a *Application) ShardFilename(filename string) string {
	results := hash.Shard(filename, a.Shard.Width, a.Shard.Depth, true)

	return strings.Join(results, "/")
}

func (a *Application) ImageHandler(res muxer.Response, req *Request) {
	file, err := a.ImageFileFromRequest(req, true, true)

	if err != nil {
		panic(err)
	}

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(file.Content())
}

func (a *Application) Store(i *image.ImageFile) error {
	con := a.KVStore.Connection()
	defer con.Close()

	err := i.Save()

	if err != nil {
		a.Logger.Fatal(err)
		return err
	}

	a.Logger.Infof("Save thumbnail %s to storage", i.Filepath)

	key := a.WithPrefix(i.Key)

	err = con.Set(key, i.Filepath)

	if err != nil {
		a.Logger.Fatal(err)

		return err
	}

	a.Logger.Infof("Save key %s => %s to kvstore", key, i.Filepath)

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

	key := a.WithPrefix(req.Key)

	// Image from the KVStore found
	stored, err := gokvstores.String(req.Connection.Get(key))

	file.Filepath = stored

	if stored != "" {
		a.Logger.Infof("Key %s found in kvstore: %s", key, stored)

		if load {
			file, err = image.FromStorage(a.DestStorage, stored)

			if err != nil {
				return nil, err
			}
		}
	} else {
		a.Logger.Infof("Key %s not found in kvstore", key)

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

		file, err = a.Engine.Transform(file, req.Operation, req.QueryString)

		if err != nil {
			return nil, err
		}

		file.Filepath = fmt.Sprintf("%s.%s", a.ShardFilename(req.Key), file.Format())
	}

	file.Key = req.Key
	file.Storage = a.DestStorage
	file.Headers["Content-Type"] = file.ContentType()

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
