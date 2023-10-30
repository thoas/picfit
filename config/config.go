package config

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"

	"github.com/thoas/picfit/constants"
	engineconfig "github.com/thoas/picfit/engine/config"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/store"
)

// Shard is a struct to allow shard location when files are uploaded
type Shard struct {
	Depth    int
	RestOnly bool
	Width    int
}

// AllowedSize is a struct used in the allowed_sizes option
type AllowedSize struct {
	Height int
	Width  int
}

// Options is a struct to add options to the application
type Options struct {
	AllowedIPAddresses  []string      `mapstructure:"allowed_ip_addresses"`
	AllowedSizes        []AllowedSize `mapstructure:"allowed_sizes"`
	DefaultUserAgent    string        `mapstructure:"default_user_agent"`
	EnableCascadeDelete bool          `mapstructure:"enable_cascade_delete"`
	EnableDelete        bool          `mapstructure:"enable_delete"`
	EnableHealth        bool          `mapstructure:"enable_health"`
	EnablePprof         bool          `mapstructure:"enable_pprof"`
	EnableStats         bool          `mapstructure:"enable_stats"`
	EnableUpload        bool          `mapstructure:"enable_upload"`
	EnablePrometheus    bool          `mapstructure:"enable_prometheus"`
	MimetypeDetector    string        `mapstructure:"mimetype_detector"`
	FreeMemoryInterval  int           `mapstructure:"free_memory_interval"`
	TransformTimeout    int           `mapstructure:"transform_timeout"`
}

// Sentry is a struct to configure sentry using a dsn
type Sentry struct {
	DSN  string
	Tags map[string]string
}

// Config is a struct to load configuration flags
type Config struct {
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	Debug          bool
	Engine         *engineconfig.Config
	KVStore        *store.Config
	Logger         logger.Config
	Options        *Options
	Port           int
	SecretKey      string `mapstructure:"secret_key"`
	Sentry         *Sentry
	Shard          *Shard
	Storage        *storage.Config
}

// DefaultConfig returns a default config instance
func DefaultConfig() *Config {
	return &Config{
		Engine: &engineconfig.Config{
			DefaultFormat:   DefaultFormat,
			Format:          "",
			ImageBufferSize: engineconfig.DefaultImageBufferSize,
			JpegQuality:     DefaultQuality,
			MaxBufferSize:   engineconfig.DefaultMaxBufferSize,
			PngCompression:  engineconfig.DefaultPngCompression,
			Quality:         DefaultQuality,
			WebpQuality:     DefaultQuality,
		},
		Options: &Options{
			DefaultUserAgent: fmt.Sprint(DefaultUserAgent, "/", constants.Version),
			EnableDelete:     false,
			EnableUpload:     false,
			MimetypeDetector: DefaultMimetypeDetector,
		},
		Port: DefaultPort,
		KVStore: &store.Config{
			Type: "dummy",
		},
		Shard: &Shard{
			Depth:    DefaultShardDepth,
			RestOnly: DefaultShardRestOnly,
			Width:    DefaultShardWidth,
		},
	}
}

func load(content string, isPath bool) (*Config, error) {
	config := &Config{}

	defaultConfig := DefaultConfig()

	viper.SetDefault("options", defaultConfig.Options)
	viper.SetDefault("shard", defaultConfig.Shard)
	viper.SetDefault("port", defaultConfig.Port)
	viper.SetDefault("kvstore", defaultConfig.KVStore)
	viper.SetDefault("engine", defaultConfig.Engine)
	viper.SetEnvPrefix("picfit")

	var err error

	if isPath == true {
		viper.SetConfigFile(content)
		err = viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		viper.SetConfigType("json")

		err = viper.ReadConfig(bytes.NewBuffer([]byte(content)))

		if err != nil {
			return nil, err
		}
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.Options.FreeMemoryInterval == 0 {
		config.Options.FreeMemoryInterval = 10
	}

	if config.Options.TransformTimeout == 0 {
		config.Options.TransformTimeout = 10
	}

	return config, nil
}

// Load creates a Config struct from a config file path
func Load(path string) (*Config, error) {
	return load(path, true)
}

// LoadFromContent creates a Config struct from a config content
func LoadFromContent(content string) (*Config, error) {
	return load(content, false)
}
