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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var App = &Application{
	Logger: Logger{
		Info:  log.New(os.Stdout, "", 0),
		Error: log.New(os.Stderr, "Error: ", 0),
	},
}

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
	App.Router.Handle("/image/{filepath}/method/{method}/display", ImageHandler)
	App.Router.Handle("/image/method/{method}/display", ImageHandler)
	App.Router.Handle("/image/{filepath}/method/{method}/get", GetHandler)
	App.Router.Handle("/image/method/{method}/get", GetHandler)

	allowedOrigins, err := jq.ArrayOfStrings("allowed_origins")
	allowedMethods, err := jq.ArrayOfStrings("allowed_methods")

	n := negroni.Classic()
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
