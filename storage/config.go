package storage

// StorageConfig is a struct to represent a Storage (fs, s3)
type StorageConfig struct {
	Type            string
	Location        string
	BaseURL         string `mapstructure:"base_url"`
	Region          string
	ACL             string
	AccessKeyID     string `mapstructure:"access_key_id"`
	BucketName      string `mapstructure:"bucket_name"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

// Config is a struct to represent a section of storage (src, fst)
type Config struct {
	Source      *StorageConfig `mapstructure:"src"`
	Destination *StorageConfig `mapstructure:"dst"`
}
