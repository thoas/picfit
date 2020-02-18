package encoder

import (
	"go.uber.org/zap/zapcore"
)

type EncoderConfig struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey     string `json:"messageKey" yaml:"messageKey"`
	LevelKey       string `json:"levelKey" yaml:"levelKey"`
	LevelIntKey    string `json:"levelIntKey" yaml:"levelIntKey"`
	TimeKey        string `json:"timeKey" yaml:"timeKey"`
	NameKey        string `json:"nameKey" yaml:"nameKey"`
	EnvKey         string `json:"envKey" yaml:"envKey"`
	CallerKey      string `json:"callerKey" yaml:"callerKey"`
	FieldsGroupKey string `json:"fieldsGroupKey" yaml:"fieldsGroupKey"`
	StacktraceKey  string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LineEnding     string `json:"lineEnding" yaml:"lineEnding"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    zapcore.LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     zapcore.TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration zapcore.DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   zapcore.CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName zapcore.NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
}
