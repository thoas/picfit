package server

import (
	"context"
	"net/http"

	"github.com/thoas/picfit"
	"github.com/thoas/picfit/config"
)

type Server struct {
	http *HTTPServer
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpServer, err := NewHTTPServer(cfg, processor)
	if err != nil {
		return nil, err
	}

	server := &Server{http: httpServer}
	return server, nil
}

func (s *Server) Run() error {
	return s.http.Run()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.http.engine.ServeHTTP(w, req)
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

	return server.Run()
}
