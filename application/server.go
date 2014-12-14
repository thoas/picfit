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
	"s3": S3StorageParameter,
	"fs": FileSystemStorageParameter,
}

func Run(path string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	for _, initializer := range Initializers {
		err = initializer(jq)

		panicIf(err)
	}

	App.Router = mux.NewRouter()
	App.Router.NotFoundHandler = NotFoundHandler()

	App.Router.Handle("/redirect", RedirectHandler)
	App.Router.Handle("/redirect/{sig}/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", RedirectHandler)
	App.Router.Handle("/redirect/{sig}/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", RedirectHandler)
	App.Router.Handle("/redirect/{sig}/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", RedirectHandler)
	App.Router.Handle("/redirect/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", RedirectHandler)
	App.Router.Handle("/redirect/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", RedirectHandler)
	App.Router.Handle("/redirect/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", RedirectHandler)

	App.Router.Handle("/display", ImageHandler)
	App.Router.Handle("/display/{sig}/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", ImageHandler)
	App.Router.Handle("/display/{sig}/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", ImageHandler)
	App.Router.Handle("/display/{sig}/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", ImageHandler)
	App.Router.Handle("/display/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", ImageHandler)
	App.Router.Handle("/display/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", ImageHandler)
	App.Router.Handle("/display/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", ImageHandler)

	App.Router.Handle("/get", GetHandler)
	App.Router.Handle("/get/{sig}/op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", GetHandler)
	App.Router.Handle("/get/{sig}/op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", GetHandler)
	App.Router.Handle("/get/{sig}/op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", GetHandler)
	App.Router.Handle("/get/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", GetHandler)
	App.Router.Handle("/get/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", GetHandler)
	App.Router.Handle("/get/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", GetHandler)

	allowedOrigins, err := jq.ArrayOfStrings("allowed_origins")
	allowedMethods, err := jq.ArrayOfStrings("allowed_methods")

	debug, err := jq.Bool("debug")

	if err != nil {
		debug = true
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
