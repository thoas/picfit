package context

import (
	"github.com/gin-gonic/gin"
	"github.com/thoas/gokvstores"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/kvstore"
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

// Storage extracts a storage from the gin context
func Storage(c *gin.Context) gostorages.Storage {
	return c.MustGet("storage").(gostorages.Storage)
}

// SetKVStore adds a kvstore to the gin context
func SetKVStore(s gokvstores.KVStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		kvstore.ToContext(c, s)
		c.Next()
	}
}

// Store extracts a kvstore from the gin context
func Store(c *gin.Context) gokvstores.KVStore {
	return c.MustGet("kvstore").(gokvstores.KVStore)
}
