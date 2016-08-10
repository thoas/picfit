package context

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/thoas/gokvstores"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/kvstore"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
)

// SetEngine adds an engine in the gin context
func SetEngine(e engine.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		engine.ToContext(c, e)
		c.Next()
	}
}

// Engine extracts an engine from the gin context
func Engine(c *gin.Context) engine.Engine {
	return c.MustGet("engine").(engine.Engine)
}

// SetConfig adds a config to the gin context
func SetConfig(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ToContext(c, cfg)
		c.Next()
	}
}

// Config extracts a config from the gin context
func Config(c *gin.Context) config.Config {
	return c.MustGet("config").(config.Config)
}

// SetDestinationStorage adds a destionation storage to the gin context
func SetDestinationStorage(s gostorages.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		storage.DestinationToContext(c, s)
		c.Next()
	}
}

// SetSourceStorage adds a source storage in the gin context
func SetSourceStorage(s gostorages.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		storage.SourceToContext(c, s)
		c.Next()
	}
}

// SourceStorage extracts a source storage from the gin context
func SourceStorage(c *gin.Context) gostorages.Storage {
	return c.MustGet("srcStorage").(gostorages.Storage)
}

// DestinationStorage extracts a destination storage from the gin context
func DestinationStorage(c *gin.Context) gostorages.Storage {
	return c.MustGet("dstStorage").(gostorages.Storage)
}

// SetKVStore adds a kvstore to the gin context
func SetKVStore(s gokvstores.KVStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		kvstore.ToContext(c, s)
		c.Next()
	}
}

// KVStore extracts a kvstore from the gin context
func KVStore(c *gin.Context) gokvstores.KVStore {
	return c.MustGet("kvstore").(gokvstores.KVStore)
}

// SetLogger adds a logger to the gin context
func SetLogger(l logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.ToContext(c, l)
		c.Next()
	}
}

// Logger extracts a logger from the gin context
func Logger(c *gin.Context) logrus.Logger {
	return c.MustGet("logger").(logrus.Logger)
}
