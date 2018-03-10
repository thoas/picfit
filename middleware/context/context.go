package context

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/thoas/gokvstores"

	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/kvstore"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
)

// SetContext adds application context in the gin context
func SetContext(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		e := engine.FromContext(ctx)
		engine.ToContext(c, e)

		cfg := config.FromContext(ctx)
		config.ToContext(c, cfg)

		s := storage.DestinationFromContext(ctx)
		storage.DestinationToContext(c, s)

		s = storage.SourceFromContext(ctx)
		storage.SourceToContext(c, s)

		k := kvstore.FromContext(ctx)
		kvstore.ToContext(c, k)

		l := logger.FromContext(ctx)
		logger.ToContext(c, l)
		c.Next()
	}
}

// Engine extracts an engine from the gin context
func Engine(c *gin.Context) engine.Engine {
	return c.MustGet("engine").(engine.Engine)
}

// Config extracts a config from the gin context
func Config(c *gin.Context) config.Config {
	return c.MustGet("config").(config.Config)
}

// SourceStorage extracts a source storage from the gin context
func SourceStorage(c *gin.Context) gostorages.Storage {
	return c.MustGet("srcStorage").(gostorages.Storage)
}

// DestinationStorage extracts a destination storage from the gin context
func DestinationStorage(c *gin.Context) gostorages.Storage {
	return c.MustGet("dstStorage").(gostorages.Storage)
}

// KVStore extracts a kvstore from the gin context
func KVStore(c *gin.Context) gokvstores.KVStore {
	return c.MustGet("kvstore").(gokvstores.KVStore)
}

// SetLogger adds a logger to the gin context
func SetLogger(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// Logger extracts a logger from the gin context
func Logger(c *gin.Context) logger.Logger {
	return c.MustGet("logger").(logger.Logger)
}
