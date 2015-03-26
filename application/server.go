package application

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
	"runtime"
	"strconv"
)

func Run(path string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app, err := NewFromConfig(path)

	if err != nil {
		return err
	}

	n := app.InitRouter()

	server := &http.Server{Addr: fmt.Sprintf(":%s", strconv.Itoa(app.Port())), Handler: n}

	gracehttp.Serve(server)

	return nil
}
