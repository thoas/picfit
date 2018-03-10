package server

import (
	"context"
	"net/http"

	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/config"
)

type Server struct {
	http *HTTPServer
}

func New(ctx context.Context) (*Server, error) {
	cfg := config.FromContext(ctx)

	httpServer, err := NewHTTPServer(cfg, WithContext(ctx))
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
	s.http.ServeHTTP(w, req)
}

// Run runs the application and launch servers
func Run(path string) error {
	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	ctx, err := application.Load(cfg)
	if err != nil {
		return err
	}

	server, err := New(ctx)
	if err != nil {
		return err
	}

	return server.Run()
}
