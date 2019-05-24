package server

import (
	"net/http"

	"github.com/thoas/picfit"
	"github.com/thoas/picfit/config"
)

type Server struct {
	http *HTTPServer
}

func New(cfg *config.Config) (*Server, error) {
	processor, err := picfit.NewProcessor(cfg)
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
	s.http.ServeHTTP(w, req)
}

// Run runs the application and launch servers
func Run(path string) error {
	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	server, err := New(cfg)
	if err != nil {
		return err
	}

	return server.Run()
}
