package server

import (
	"fmt"
	"strconv"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"github.com/thoas/picfit/application"
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

// Load loads the application and launch the webserver
func Load(path string) error {
	ctx, err := application.Load(path)

	if err != nil {
		return err
	}

	return Run(ctx)
}

// Run loads a new server
func Run(ctx netContext.Context) error {
	router := gin.Default()

	cfg := config.FromContext(ctx)

	methods := map[string]gin.HandlerFunc{
		"redirect": views.RedirectView,
		"display":  views.DisplayView,
		"get":      views.GetView,
	}

	kv := kvstore.FromContext(ctx)

	if cfg.Sentry != nil {
		client, err := raven.NewClient(cfg.Sentry.DSN, cfg.Sentry.Tags)

		if err != nil {
			return err
		}

		router.Use(sentry.Recovery(client, true))
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

	if cfg.AllowedOrigins != nil && cfg.AllowedMethods != nil {
		co := cors.New(cors.Options{
			AllowedOrigins: cfg.AllowedOrigins,
			AllowedMethods: cfg.AllowedMethods,
		})

		router.Use(func(c *gin.Context) {
			co.HandlerFunc(c.Writer, c.Request)

			c.Next()
		})
	}

	if cfg.Options.EnableUpload {
		router.POST("/upload", views.UploadView)
	}

	if cfg.Options.EnableDelete {
		router.DELETE("/{path:[\\w\\-/.]+}", views.DeleteView)
	}

	router.Run(fmt.Sprintf(":%s", strconv.Itoa(cfg.Port)))

	return nil
}
