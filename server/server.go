package server

import (
	"fmt"
	"strconv"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/gin"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/kvstore"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/middleware"
	"github.com/thoas/picfit/middleware/context"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/views"
	netContext "golang.org/x/net/context"
)

// Run loads a new server
func Run(ctx netContext.Context) {
	router := gin.Default()

	cfg := config.FromContext(ctx)

	methods := map[string]gin.HandlerFunc{
		"redirect": views.RedirectView,
		"display":  views.DisplayView,
		"get":      views.GetView,
	}

	kv := kvstore.FromContext(ctx)

	if cfg.Sentry != nil {
		raven.SetDSN(cfg.Sentry.DSN)
		router.Use(sentry.Recovery(raven.DefaultClient, true))
	}

	router.Use(context.SetLogger(logger.FromContext(ctx)))
	router.Use(context.SetKVStore(kv))
	router.Use(context.SetSourceStorage(storage.SourceFromContext(ctx)))
	router.Use(context.SetDestinationStorage(storage.DestinationFromContext(ctx)))
	router.Use(context.SetEngine(engine.FromContext(ctx)))
	router.Use(context.SetConfig(cfg))

	for name, view := range methods {
		views := []gin.HandlerFunc{
			middleware.ParametersParser(),
			middleware.KeyParser(),
			middleware.URLParser(),
			middleware.OperationParser(),
			view,
		}

		router.GET(fmt.Sprintf("/%s", name), views...)
		router.GET(fmt.Sprintf("/%s/*parameters", name), views...)
	}

	if cfg.Options.EnableUpload {
		router.POST("/upload", views.UploadView)
	}

	if cfg.Options.EnableDelete {
		router.DELETE("/{path:[\\w\\-/.]+}", views.DeleteView)
	}

	router.Run(fmt.Sprintf(":%s", strconv.Itoa(cfg.Port)))
}
