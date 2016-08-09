package config

import "github.com/spf13/viper"

// Shard is a struct to allow shard location when files are uploaded
type Shard struct {
	Depth int
	Width int
}

// Options is a struct to add options to the application
type Options struct {
	EnableUpload  bool
	EnableDelete  bool
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
	BaseURL         string
	Region          string
	ACL             string
	AccessKeyID     string
	BucketName      string
	SecretAccessKey string
}

// Storages is a struct to represent a section of storage (src, fst)
type Storages struct {
	Src *Storage
	Dst *Storage
}

// Sentry is a struct to configure sentry using a dsn
type Sentry struct {
	Dsn  string
	Tags map[string]string
}

// Config is a struct to load configuration flags
type Config struct {
	Debug          bool
	Sentry         *Sentry
	SecretKey      string
	Shard          *Shard
	Port           int
	Options        *Options
	AllowedOrigins []string
	AllowedMethods []string
	Storage        *Storages
	KVStore        *KVStore
}

// Load creates a Config struct from a config file path
func Load(path string) (*Config, error) {
	config := &Config{}

	viper.SetDefault("options", &Options{
		EnableDelete:  false,
		EnableUpload:  false,
		DefaultFormat: DefaultFormat,
		Quality:       DefaultQuality,
		Format:        "",
	})

	viper.SetDefault("port", DefaultPort)

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
