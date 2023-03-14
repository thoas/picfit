package picfit

import (
	"context"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/store"
)

// NewProcessor returns a Processor instance from a config.Config instance
func NewProcessor(ctx context.Context, cfg *config.Config) (*Processor, error) {
	log := logger.New(cfg.Logger)

	sourceStorage, destinationStorage, err := storage.New(ctx,
		log.With(logger.String("logger", "storage")), cfg.Storage)
	if err != nil {
		return nil, err
	}

	s, err := store.New(ctx,
		log.With(logger.String("logger", "store")),
		cfg.KVStore)
	if err != nil {
		return nil, err
	}

	e := engine.New(*cfg.Engine, log.With(logger.String("logger", "engine")))

	log.Info("Image engine configured",
		logger.String("engine", e.String()))

	return &Processor{
		Logger: log,

		config:             cfg,
		destinationStorage: *destinationStorage,
		engine:             e,
		sourceStorage:      *sourceStorage,
		store:              s,
	}, nil
}
