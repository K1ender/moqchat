package httptransport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/K1ender/moqchat/internal/config"
)

type Server struct {
	srv *http.Server
}

const ReadTimeout = 10 * time.Second
const WriteTimeout = 10 * time.Second

func NewServer(cfg config.HTTPConfig) Server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           mux,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		ReadHeaderTimeout: ReadTimeout,
	}

	return Server{
		srv: srv,
	}
}

func (s *Server) Run() error {
	err := s.srv.ListenAndServe()
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return fmt.Errorf("http server failed: %w", err)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
