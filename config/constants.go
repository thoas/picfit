package config

const (
	// DefaultFormat is the default image format
	DefaultFormat = "png"

	// DefaultQuality is the default quality for processed images
	DefaultQuality = 95

	// DefaultUserAgent is the default user-agent header to fetch images from URL with.
	// n.b. application version later appended to this.
	DefaultUserAgent = "picfit"

	// DefaultMimetypeDetector method to use
	DefaultMimetypeDetector = "extension"

	// DefaultPort is the default port of the application server
	DefaultPort = 3001

	// DefaultShardWidth is the default shard width
	DefaultShardWidth = 0

	// DefaultShardDepth is the default shard depth
	DefaultShardDepth = 0

	// DefaultShardRestOnly is the default shard rest behaviour
	DefaultShardRestOnly = true
)
