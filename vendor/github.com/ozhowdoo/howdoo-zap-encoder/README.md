# Howdoo zap encoder

## Install

```sh
$ go get github.com/ozhowdoo/howdoo-zap-encoder
```

## Usage encoder

```go

config := encoder.Config{
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

logger := config.Build()

logger.Info("Test");

```


## Usage handlers

```go

r := gin.New()

// Logs all requests
r.Use(encoder.LoggerRequest(logger))

// Logs all panic errors
r.Use(encoder.LoggerRecovery(logger))

```