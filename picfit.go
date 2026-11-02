package picfit

import (
	"context"
	"log/slog"

	"github.com/thoas/picfit/constants"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	httppkg "github.com/thoas/picfit/http"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/store"
)

// NewProcessor returns a Processor instance from a config.Config instance
func NewProcessor(ctx context.Context, cfg *config.Config) (*Processor, error) {
	cfg.Logger.ContextKeys = []string{constants.RequestIDCtx}
	log := logger.New(cfg.Logger)

	sourceStorage, destinationStorage, destinationReadOnlyStorage, err := storage.New(ctx,
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

	policy := httppkg.NewFetchPolicy(cfg.Options.BlockPrivateNetworks, cfg.Options.AllowedSourceHosts)
	processor := &Processor{
		Logger: log,

		config:                     cfg,
		destinationStorage:         destinationStorage,
		destinationReadOnlyStorage: destinationReadOnlyStorage,
		engine:                     e,
		sourceStorage:              sourceStorage,
		store:                      s,
		urlPolicy:                  policy,
		urlClient: httppkg.NewClient(
			httppkg.WithUserAgent(cfg.Options.DefaultUserAgent),
			httppkg.WithFetchPolicy(policy),
		),
	}
	if cfg.Options.MaxProcessorConcurrent != nil {
		processor.withSemaphore = true
		processor.semaphore = make(chan struct{}, *cfg.Options.MaxProcessorConcurrent)
		processor.semaphoreOperations = cfg.Options.MaxProcessorConcurrentOperations
	}
	if cfg.Options.MaxImageDimensions != nil {
		processor.maxImageDimensions = cfg.Options.MaxImageDimensions
	}
	if cfg.Options.MaxImagePixels > 0 {
		processor.maxImagePixels = cfg.Options.MaxImagePixels
	}
	return processor, nil
}
