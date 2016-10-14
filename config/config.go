package config

import (
	"bytes"

	"github.com/spf13/viper"
)

// Shard is a struct to allow shard location when files are uploaded
type Shard struct {
	Depth int
	Width int
	RestOnly bool
}

// Options is a struct to add options to the application
type Options struct {
	EnableUpload  bool `mapstructure:"enable_upload"`
	EnableDelete  bool `mapstructure:"enable_delete"`
	DefaultFormat string
	Format        string
	Quality       int
}

// KVStore is a struct to represent a key/value store (redis, cache)
type KVStore struct {
	Type       string
	Host       string
	Port       int
	Password   string
	Db         int
	Prefix     string
	MaxEntries int
}

// Storage is a struct to represent a Storage (fs, s3)
type Storage struct {
	Type            string
	Location        string
	BaseURL         string `mapstructure:"base_url"`
	Region          string
	ACL             string
	AccessKeyID     string `mapstructure:"access_key_id"`
	BucketName      string `mapstructure:"bucket_name"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

// Storages is a struct to represent a section of storage (src, fst)
type Storages struct {
	Src *Storage
	Dst *Storage
}

// Sentry is a struct to configure sentry using a dsn
type Sentry struct {
	DSN  string
	Tags map[string]string
}

// Config is a struct to load configuration flags
type Config struct {
	Debug          bool
	Sentry         *Sentry
	SecretKey      string `mapstructure:"secret_key"`
	Shard          *Shard
	Port           int
	Options        *Options
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	Storage        *Storages
	KVStore        *KVStore
}

// DefaultConfig returns a default config instance
func DefaultConfig() *Config {
	return &Config{
		Options: &Options{
			EnableDelete:  false,
			EnableUpload:  false,
			DefaultFormat: DefaultFormat,
			Quality:       DefaultQuality,
			Format:        "",
		},
		Port: DefaultPort,
		KVStore: &KVStore{
			Type: "dummy",
		},
		Shard: &Shard{
			Width: DefaultShardWidth,
			Depth: DefaultShardDepth,
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

	if config.Options.Quality == 0 {
		config.Options.Quality = defaultConfig.Options.Quality
	}

	if config.Options.DefaultFormat == "" {
		config.Options.DefaultFormat = defaultConfig.Options.DefaultFormat
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
