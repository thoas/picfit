package encoder

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type ConfigBuilder interface {
	Build(opts ...zap.Option) (*zap.Logger, error)
}

type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	Level zap.AtomicLevel `json:"level" yaml:"level"`

	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling"`

	// EncoderConfig sets options for the chosen encoder.
	EncoderConfig EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	// OutputPaths is a list of URLs or file paths to write logging output to.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`

	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`

	// InitialFields is a collection of fields to add to the root logger.
	InitialFields []Field `json:"initialFields" yaml:"initialFields"`

	EncoderConstructor func(interface{}) (zapcore.Encoder, error)
}

// Build constructs a logger from the Config and Options.
func (cfg Config) Build(opts ...zap.Option) (*zap.Logger, error) {
	enc, err := cfg.EncoderConstructor(&cfg.EncoderConfig)
	if err != nil {
		return nil, err
	}

	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	log := zap.New(
		zapcore.NewCore(enc, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

func (cfg Config) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	opts = append(opts, zap.AddCaller())

	stackLevel := zap.ErrorLevel
	opts = append(opts, zap.AddStacktrace(stackLevel))

	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}

	if len(cfg.InitialFields) > 0 {
		opts = append(opts, zap.Fields(cfg.InitialFields...))
	}

	return opts
}

func (cfg Config) openSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := zap.Open(cfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}
	errSink, _, err := zap.Open(cfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}
	return sink, errSink, nil
}
