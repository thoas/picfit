package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	api "gopkg.in/fukata/golang-stats-api-handler.v1"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thoas/picfit"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/failure"
	"github.com/thoas/picfit/middleware"
	"github.com/thoas/stats"
)

type HTTPServer struct {
	config    *config.Config
	engine    *gin.Engine
	processor *picfit.Processor
}

func NewHTTPServer(cfg *config.Config, processor *picfit.Processor) (*HTTPServer, error) {
	server := &HTTPServer{
		config:    cfg,
		processor: processor,
	}
	if err := server.Init(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *HTTPServer) Init() error {
	var (
		router    = gin.New()
		handlers  = &handlers{s.processor}
		endpoints = []endpoint{
			{
				pattern: "redirect",
				handler: failure.Handle(handlers.redirect),
				method:  router.GET,
			},
			{
				pattern: "display",
				handler: failure.Handle(handlers.display),
				method:  router.GET,
			},
			{
				pattern: "get",
				handler: failure.Handle(handlers.get),
				method:  router.GET,
			},
		}
	)

	if s.config.Debug {
		router.Use(gin.Recovery())
	} else {
		router.Use(middleware.Recover)
	}

	router.Use(middleware.NewLogger(s.config, s.processor.Logger))
	router.Use(middleware.Metrics)

	if s.config.Sentry != nil {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: s.config.Sentry.DSN,
		}); err != nil {
			return err
		}

		router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}

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

	router.GET("/healthcheck", handlers.healthcheck(time.Now().UTC()))

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
			failure.Handle(handlers.upload))
	}

	if s.config.Options.EnableDelete {
		router.DELETE("/*parameters",
			restrictIPAddresses,
			middleware.ParametersParser(),
			middleware.KeyParser(),
			failure.Handle(handlers.delete))
	}

	router.GET("/error", handlers.internalError)

	if s.config.Options.EnablePrometheus {
		router.GET("/metrics", prometheusHandler())
	}

	if s.config.Options.EnablePprof {
		prefixRouter := router.Group("/debug/pprof")
		{
			prefixRouter.GET("/",
				restrictIPAddresses,
				pprofHandler(pprof.Index))
			prefixRouter.GET("/cmdline",
				restrictIPAddresses,
				pprofHandler(pprof.Cmdline))
			prefixRouter.GET("/profile",
				restrictIPAddresses,
				pprofHandler(pprof.Profile))
			prefixRouter.POST("/symbol",
				restrictIPAddresses,
				pprofHandler(pprof.Symbol))
			prefixRouter.GET("/symbol",
				restrictIPAddresses,
				pprofHandler(pprof.Symbol))
			prefixRouter.GET("/trace",
				restrictIPAddresses,
				pprofHandler(pprof.Trace))
			prefixRouter.GET("/block",
				restrictIPAddresses,
				pprofHandler(pprof.Handler("block").ServeHTTP))
			prefixRouter.GET("/goroutine",
				restrictIPAddresses,
				pprofHandler(pprof.Handler("goroutine").ServeHTTP))
			prefixRouter.GET("/heap",
				restrictIPAddresses,
				pprofHandler(pprof.Handler("heap").ServeHTTP))
			prefixRouter.GET("/mutex",
				restrictIPAddresses,
				pprofHandler(pprof.Handler("mutex").ServeHTTP))
			prefixRouter.GET("/threadcreate",
				restrictIPAddresses,
				pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
		}
	}

	s.engine = router

	return nil
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Run loads a new http server
func (s *HTTPServer) Run() error {
	s.engine.Run(fmt.Sprintf(":%s", strconv.Itoa(s.config.Port)))

	return nil
}
