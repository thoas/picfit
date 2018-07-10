package config

// DefaultFormat is the default image format
const DefaultFormat = "png"

// DefaultQuality is the default quality for processed images
const DefaultQuality = 95

// DefaultUserAgent is the default user-agent header to fetch images from URL with.
// n.b. application version later appended to this.
const DefaultUserAgent = "picfit"

// DefaultMimetypeDetector method to use
const DefaultMimetypeDetector = "extension"

// DefaultPort is the default port of the application server
const DefaultPort = 3001

// DefaultShardWidth is the default shard width
const DefaultShardWidth = 0

// DefaultShardDepth is the default shard depth
const DefaultShardDepth = 0

// DefaultShardRestOnly is the default shard rest behaviour
const DefaultShardRestOnly = true

// DefaultPngCompression is the default compression for png.
const DefaultPngCompression = 0

// DefaultMaxBufferSize is the maximum size of buffer for lilliput
const DefaultMaxBufferSize = 8192

// DefaultImageBufferSize is the default image buffer size for lilliput
const DefaultImageBufferSize = 50 * 1024 * 1024
