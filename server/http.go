package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	api "gopkg.in/fukata/golang-stats-api-handler.v1"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/gin"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/middleware"
	"github.com/thoas/picfit/middleware/context"
	"github.com/thoas/picfit/server/handlers"
	"github.com/thoas/stats"
)

type handlerMethod func(string, ...gin.HandlerFunc) gin.IRoutes
type endpoint struct {
	pattern string
	handler gin.HandlerFunc
	method  handlerMethod
}

type HTTPServer struct {
	*gin.Engine
	config config.Config
}

func NewHTTPServer(cfg config.Config, opt ...Option) (*HTTPServer, error) {
	opts := NewOptions(opt...)

	server := &HTTPServer{
		config: cfg,
	}
	err := server.Init(opts)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *HTTPServer) Init(opts Options) error {
	var (
		router    = gin.New()
		endpoints = []endpoint{
			{
				pattern: "redirect",
				handler: handlers.Redirect,
				method:  router.GET,
			},
			{
				pattern: "display",
				handler: handlers.Display,
				method:  router.GET,
			},
			{
				pattern: "get",
				handler: handlers.Get,
				method:  router.GET,
			},
		}
	)

	if s.config.Debug {
		router.Use(gin.Recovery())
	}

	if s.config.Logger.GetLevel() == "debug" {
		router.Use(gin.Logger())
	}

	if s.config.Sentry != nil {
		client, err := raven.NewClient(s.config.Sentry.DSN, s.config.Sentry.Tags)
		if err != nil {
			return err
		}

		router.Use(sentry.Recovery(client, true))
	}

	router.Use(context.SetContext(opts.Context))

	if s.config.AllowedOrigins != nil && s.config.AllowedMethods != nil {
		allowedOrigins := s.config.AllowedOrigins

		allowAllOrigins := false

		if len(allowedOrigins) == 1 {
			if allowedOrigins[0] == "*" {
				allowAllOrigins = true
			}
		}

		if allowAllOrigins {
			allowedOrigins = nil
		}

		router.Use(cors.New(cors.Config{
			AllowAllOrigins: allowAllOrigins,
			AllowedOrigins:  allowedOrigins,
			AllowedMethods:  s.config.AllowedMethods,
			AllowedHeaders:  s.config.AllowedHeaders,
		}))
	}

	router.GET("/healthcheck", handlers.Healthcheck(time.Now().UTC()))

	restrictIPAddresses := middleware.RestrictIPAddresses(s.config.Options.AllowedIPAddresses)

	if s.config.Options.EnableStats {
		s := stats.New()

		router.Use(func() gin.HandlerFunc {
			return func(c *gin.Context) {
				beginning, recorder := s.Begin(c.Writer)
				c.Next()
				s.End(beginning, recorder)
			}
		}())

		router.GET("/sys/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, s.Data())
		})
	}

	if s.config.Options.EnableHealth {
		router.GET("/sys/health",
			restrictIPAddresses,
			func(c *gin.Context) {
				c.JSON(http.StatusOK,
					api.GetStats())
			})
	}

	for _, e := range endpoints {
		views := []gin.HandlerFunc{
			middleware.ParametersParser(),
			middleware.KeyParser(),
			middleware.Security(s.config.SecretKey),
			middleware.URLParser(s.config.Options.MimetypeDetector),
			middleware.OperationParser(),
			middleware.RestrictSizes(s.config.Options.AllowedSizes),
			e.handler,
		}

		e.method(fmt.Sprintf("/%s", e.pattern), views...)

		if s.config.Storage != nil && s.config.Storage.Source != nil {
			e.method(fmt.Sprintf("/%s/*parameters", e.pattern), views...)
		}
	}

	if s.config.Options.EnableUpload {
		router.POST("/upload",
			restrictIPAddresses,
			handlers.Upload)
	}

	if s.config.Options.EnableDelete {
		router.DELETE("/*parameters",
			restrictIPAddresses,
			middleware.ParametersParser(),
			middleware.KeyParser(),
			handlers.Delete)
	}

	if s.config.Options.EnablePprof {
		prefixRouter := router.Group("/debug/pprof")
		{
			prefixRouter.GET("/",
				restrictIPAddresses,
				handlers.Pprof(pprof.Index))
			prefixRouter.GET("/cmdline",
				restrictIPAddresses,
				handlers.Pprof(pprof.Cmdline))
			prefixRouter.GET("/profile",
				restrictIPAddresses,
				handlers.Pprof(pprof.Profile))
			prefixRouter.POST("/symbol",
				restrictIPAddresses,
				handlers.Pprof(pprof.Symbol))
			prefixRouter.GET("/symbol",
				restrictIPAddresses,
				handlers.Pprof(pprof.Symbol))
			prefixRouter.GET("/trace",
				restrictIPAddresses,
				handlers.Pprof(pprof.Trace))
			prefixRouter.GET("/block",
				restrictIPAddresses,
				handlers.Pprof(pprof.Handler("block").ServeHTTP))
			prefixRouter.GET("/goroutine",
				restrictIPAddresses,
				handlers.Pprof(pprof.Handler("goroutine").ServeHTTP))
			prefixRouter.GET("/heap",
				restrictIPAddresses,
				handlers.Pprof(pprof.Handler("heap").ServeHTTP))
			prefixRouter.GET("/mutex",
				restrictIPAddresses,
				handlers.Pprof(pprof.Handler("mutex").ServeHTTP))
			prefixRouter.GET("/threadcreate",
				restrictIPAddresses,
				handlers.Pprof(pprof.Handler("threadcreate").ServeHTTP))
		}
	}

	s.Engine = router

	return nil
}

// Run loads a new http server
func (s *HTTPServer) Run() error {
	s.Engine.Run(fmt.Sprintf(":%s", strconv.Itoa(s.config.Port)))

	return nil
}
