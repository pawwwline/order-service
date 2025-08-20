package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order-service/internal/config"
	"order-service/internal/controller/http/handlers"
	mid "order-service/internal/controller/http/middleware"
	"order-service/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg         *config.HTTPConfig
	httpHandler *handlers.HTTPHandler
	logger      *slog.Logger
	httpServer  *http.Server
}

func NewServer(cfg *config.HTTPConfig, uc *usecase.OrderUseCase, l *slog.Logger) *Server {
	return &Server{
		cfg:         cfg,
		httpHandler: handlers.NewHTTPHandler(uc),
		logger:      l,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
		},
	}
}

func (s *Server) Run() {
	s.httpServer.Handler = s.Routes()
	s.logger.Info("HTTP server starting", "addr", fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port))
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("server failed", "error", err)

		return
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(mid.RequestLogger(s.logger))
	r.Use(middleware.Recoverer)
	s.httpHandler.RegisterStaticRoutes(r)
	
	s.httpHandler.RegisterRoutes(r)

	return r
}
