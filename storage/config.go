package storage

// StorageConfig is a struct to represent a Storage (fs, s3)
type StorageConfig struct {
	ACL             string
	AccessKeyID     string `mapstructure:"access_key_id"`
	BaseURL         string `mapstructure:"base_url"`
	BucketName      string `mapstructure:"bucket_name"`
	CacheControl    string `mapstructure:"cache_control"`
	Location        string
	Region          string
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Type            string
	Endpoint        string `mapstructure:"endpoint"`
	Name            string `mapstructure:"name"`
}

// Config is a struct to represent a section of storage (src, fst)
type Config struct {
	Destination         *StorageConfig `mapstructure:"dst"`
	Source              *StorageConfig `mapstructure:"src"`
	DestinationReadOnly *StorageConfig `mapstructure:"dst-ro"`
}
