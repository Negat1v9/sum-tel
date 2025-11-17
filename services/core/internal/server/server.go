package server

import (
	"context"
	"net/http"

	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

type Server struct {
	log    *logger.Logger
	server http.Server
	cfg    *config.CoreConfig
}

func New(cfg *config.CoreConfig, log *logger.Logger) *Server {
	return &Server{
		log: log,
		server: http.Server{
			Addr: cfg.WebConfig.ListenAddress,
		},
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
