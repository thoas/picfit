package application

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	"github.com/jmoiron/jsonq"
	"github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	"github.com/thoas/picfit/middleware"
	"github.com/thoas/picfit/util"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var App = NewApplication()

var KVStores = map[string]KVStoreParameter{
	"redis": RedisKVStoreParameter,
	"cache": CacheKVStoreParameter,
}

var Storages = map[string]StorageParameter{
	"http+s3": HTTPS3StorageParameter,
	"s3":      S3StorageParameter,
	"http+fs": HTTPFileSystemStorageParameter,
	"fs":      FileSystemStorageParameter,
}

func Run(path string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	err = dec.Decode(&data)

	if err != nil {
		return err
	}

	jq := jsonq.NewQuery(data)

	for _, initializer := range Initializers {
		err = initializer(jq)

		util.PanicIf(err)
	}

	App.Router = mux.NewRouter()
	App.Router.NotFoundHandler = NotFoundHandler()

	methods := map[string]Handler{
		"redirect": RedirectHandler,
		"display":  ImageHandler,
		"get":      GetHandler,
	}

	for name, handler := range methods {
		App.Router.Handle(fmt.Sprintf("/%s", name), handler)
		App.Router.Handle(fmt.Sprintf("/%s/{sig}/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handler)
		App.Router.Handle(fmt.Sprintf("/%s/{sig}/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", name), handler)
		App.Router.Handle(fmt.Sprintf("/%s/{sig}/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handler)
		App.Router.Handle(fmt.Sprintf("/%s/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handler)
		App.Router.Handle(fmt.Sprintf("/%s/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", name), handler)
		App.Router.Handle(fmt.Sprintf("/%s/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), handler)
	}

	allowedOrigins, err := jq.ArrayOfStrings("allowed_origins")
	allowedMethods, err := jq.ArrayOfStrings("allowed_methods")

	debug, err := jq.Bool("debug")

	if err != nil {
		debug = false
	}

	n := negroni.New(&middleware.Recovery{
		Raven:      App.Raven,
		Logger:     App.Logger.Error,
		PrintStack: debug,
		StackAll:   false,
		StackSize:  1024 * 8,
	}, &middleware.Logger{App.Logger.Info})
	n.Use(cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
	}))
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(App.Router)

	port, err := jq.Int("port")

	if err != nil {
		port = DefaultPort
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", strconv.Itoa(port)), Handler: n}

	gracehttp.Serve(server)

	return nil
}
