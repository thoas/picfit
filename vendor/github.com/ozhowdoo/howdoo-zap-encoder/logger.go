package encoder

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

type Logger interface {
	Panic(string, ...Field)

	Info(string, ...Field)

	Error(string, ...Field)

	Sync() error

	Debug(string, ...Field)

	With(fields ...Field) *zap.Logger
}
