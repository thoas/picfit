package logger

import (
	"github.com/ozhowdoo/howdoo-zap-encoder"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

type Logger interface {
	Panic(string, ...Field)
	Fatal(string, ...Field)
	Info(string, ...Field)
	Error(string, ...Field)
	Sync() error
	Debug(string, ...Field)
	With(fields ...Field) *zap.Logger
}

func New(cfg Config) Logger {

	initialFields := []Field{
		String("app", cfg.App),
		String("channel", cfg.Channel),
		{Key: "extra", Type: zapcore.ObjectMarshalerType, Interface: encoder.ArrayFields([]Field{})},
	}
	sampling := &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}
	var config encoder.ConfigBuilder
	if cfg.GetType() == HowdooJsonType {

		config = encoder.Config{
			Level:            GetAtomicLevel(cfg.GetLevel()),
			Sampling:         sampling,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
			InitialFields:    initialFields,
			EncoderConfig: encoder.EncoderConfig{
				TimeKey:        "datetime",
				LevelKey:       "level_name",
				LevelIntKey:    "level",
				EnvKey:         "env",
				CallerKey:      "script_name",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				FieldsGroupKey: "context",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeCaller:   zapcore.FullCallerEncoder,
			},
			EncoderConstructor: func(encoderConfig interface{}) (zapcore.Encoder, error) {
				enc := encoder.NewEncoder(encoderConfig)
				return enc, nil
			},
		}
	} else {

		var encoderConfig zapcore.EncoderConfig
		if cfg.GetType() == JsonType {
			encoderConfig = zap.NewProductionEncoderConfig()
		} else {
			encoderConfig = zap.NewDevelopmentEncoderConfig()
		}

		config = zap.Config{
			Level:            GetAtomicLevel(cfg.GetLevel()),
			Development:      cfg.GetLevel() == DebugLevel,
			Sampling:         sampling,
			Encoding:         cfg.GetType(),
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	logger, _ := config.Build()

	return logger
}

func String(k, v string) Field {
	return zap.String(k, v)
}

func Duration(k string, d time.Duration) Field {
	return zap.Duration(k, d)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Int(k string, i int) Field {
	return zap.Int(k, i)
}

func Array(key string, val zapcore.ArrayMarshaler) Field {
	return zap.Array(key, val)
}

func Int64(k string, i int64) Field {
	return zap.Int64(k, i)
}

func Error(v error) Field {
	return zap.Error(v)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

func GetAtomicLevel(level string) zap.AtomicLevel {

	var atomicLevel zap.AtomicLevel
	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		atomicLevel.SetLevel(zap.DebugLevel)
	}

	return atomicLevel
}
