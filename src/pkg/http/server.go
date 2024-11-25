package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"cspage/pkg/config"
	"cspage/pkg/worker"
)

const (
	healthzURL = "/healthz"
	debugURL   = "/debug"
)

var (
	//nolint:gochecknoglobals // Azure VM Health extension requirement.
	healthzResponseHealthy = []byte(`{"ApplicationHealthState": "Healthy"}`)
	//nolint:gochecknoglobals // Azure VM Health extension requirement.
	healthzResponseUnhealthy = []byte(`{"ApplicationHealthState": "Unhealthy"}`)
)

type Server struct {
	ctx     context.Context // workerCtx
	log     *slog.Logger
	group   *worker.Group
	server  *http.Server
	timeout time.Duration
}

func NewSimpleServer(ctx context.Context, cfg *HTTPConfig, debug bool) *Server {
	return NewServer(ctx, cfg, debug, nil)
}

func NewServer(ctx context.Context, cfg *HTTPConfig, debug bool, handlers map[string]http.Handler) *Server {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestSize(cfg.HTTPMaxBodySize),
		newLogger(slog.Default()),
		middleware.Recoverer,
	)

	if cfg.HTTPCompressionLevel > 0 {
		router.Use(middleware.Compress(int(cfg.HTTPCompressionLevel)))
	}

	for path, handler := range handlers {
		router.Mount(path, handler)
	}

	if debug {
		router.Mount(debugURL, middleware.Profiler())
	}

	log := slog.With("listen-addr", cfg.HTTPListenAddr)
	server := &Server{
		ctx: ctx,
		log: log,
		server: &http.Server{
			Addr:           cfg.HTTPListenAddr,
			Handler:        router,
			ReadTimeout:    cfg.HTTPReadTimeout,
			WriteTimeout:   cfg.HTTPWriteTimeout,
			IdleTimeout:    cfg.HTTPIdleTimeout,
			MaxHeaderBytes: int(cfg.HTTPMaxHeaderSize),
			ErrorLog:       slog.NewLogLogger(log.Handler(), slog.LevelError),
		},
		timeout: cfg.HTTPShutdownTimeout,
	}
	router.Get(healthzURL, server.healthz)

	return server
}

func (s *Server) Enabled() bool {
	return true
}

func (s *Server) Run(runnerCtx context.Context, wg *worker.Group) {
	defer wg.Done()
	s.group = wg
	defer s.stop() //nolint:contextcheck // False positive?

	go func() {
		s.log.Info("HTTP server is starting...")
		wg.Start()
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			config.DieLog(s.log, "Could not start the HTTP server", "err", err)
		}
	}()

	<-runnerCtx.Done()
}

func (s *Server) stop() {
	s.log.Info("HTTP server is shutting down...")
	s.group.Stop()

	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Could not shut down the HTTP server", "err", err)
	}

	s.log.Info("HTTP server stopped")
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	var response []byte
	if s.group != nil && s.group.Healthy() {
		response = healthzResponseHealthy
		w.WriteHeader(http.StatusOK)
	} else {
		response = healthzResponseUnhealthy
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if _, err := w.Write(response); err != nil {
		s.log.Error("Could not write healthz response", "err", err)
	}
}
