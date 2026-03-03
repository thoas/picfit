package server

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/thoas/picfit"
	"github.com/thoas/picfit/config"
	"golang.org/x/sync/errgroup"
)

func New(ctx context.Context, cfg *config.Config) (*HTTPServer, error) {
	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpServer, err := NewHTTPServer(cfg, processor)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}

// Run runs the application and launch servers
func Run(ctx context.Context, path string) error {
	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	server, err := New(ctx, cfg)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		if err := server.Run(ctx); err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
