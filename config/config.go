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
	Width    int
	RestOnly bool
}

// AllowedSize is a struct used in the allowed_sizes option
type AllowedSize struct {
	Height int
	Width  int
}

// Options is a struct to add options to the application
type Options struct {
	AllowedIPAddresses  []string      `mapstructure:"allowed_ip_addresses"`
	EnablePprof         bool          `mapstructure:"enable_pprof"`
	EnableUpload        bool          `mapstructure:"enable_upload"`
	EnableDelete        bool          `mapstructure:"enable_delete"`
	EnableCascadeDelete bool          `mapstructure:"enable_cascade_delete"`
	EnableStats         bool          `mapstructure:"enable_stats"`
	EnableHealth        bool          `mapstructure:"enable_health"`
	AllowedSizes        []AllowedSize `mapstructure:"allowed_sizes"`
	DefaultUserAgent    string        `mapstructure:"default_user_agent"`
	MimetypeDetector    string        `mapstructure:"mimetype_detector"`
}

// Sentry is a struct to configure sentry using a dsn
type Sentry struct {
	DSN  string
	Tags map[string]string
}

// Config is a struct to load configuration flags
type Config struct {
	Debug          bool
	Engine         *engineconfig.Config
	Sentry         *Sentry
	SecretKey      string `mapstructure:"secret_key"`
	Shard          *Shard
	Port           int
	Options        *Options
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	Storage        *storage.Config
	KVStore        *store.Config
	Logger         logger.Config
}

// DefaultConfig returns a default config instance
func DefaultConfig() *Config {
	return &Config{
		Engine: &engineconfig.Config{
			DefaultFormat:   DefaultFormat,
			Quality:         DefaultQuality,
			JpegQuality:     DefaultQuality,
			WebpQuality:     DefaultQuality,
			PngCompression:  engineconfig.DefaultPngCompression,
			MaxBufferSize:   engineconfig.DefaultMaxBufferSize,
			ImageBufferSize: engineconfig.DefaultImageBufferSize,
			Format:          "",
		},
		Options: &Options{
			EnableDelete:     false,
			EnableUpload:     false,
			DefaultUserAgent: fmt.Sprint(DefaultUserAgent, "/", constants.Version),
			MimetypeDetector: DefaultMimetypeDetector,
		},
		Port: DefaultPort,
		KVStore: &store.Config{
			Type: "dummy",
		},
		Shard: &Shard{
			Width:    DefaultShardWidth,
			Depth:    DefaultShardDepth,
			RestOnly: DefaultShardRestOnly,
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

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
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
