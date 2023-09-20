package picfit

import (
	"context"
	"log/slog"

	"github.com/thoas/picfit/constants"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/store"
)

// NewProcessor returns a Processor instance from a config.Config instance
func NewProcessor(ctx context.Context, cfg *config.Config) (*Processor, error) {
	cfg.Logger.ContextKeys = []string{constants.RequestIDCtx}
	log := logger.New(cfg.Logger)

	sourceStorage, destinationStorage, err := storage.New(ctx,
		log.With(slog.String("logger", "storage")), cfg.Storage)
	if err != nil {
		return nil, err
	}

	s, err := store.New(ctx,
		log.With(slog.String("logger", "store")),
		cfg.KVStore)
	if err != nil {
		return nil, err
	}

	e := engine.New(*cfg.Engine, log.With(slog.String("logger", "engine")))

	log.InfoContext(ctx, "Image engine configured",
		slog.String("engine", e.String()))

	return &Processor{
		Logger: log,

		config:             cfg,
		destinationStorage: destinationStorage,
		engine:             e,
		sourceStorage:      sourceStorage,
		store:              s,
	}, nil
}
