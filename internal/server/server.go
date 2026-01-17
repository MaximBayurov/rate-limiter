package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/app"
	"github.com/MaximBayurov/rate-limiter/internal/configuration"
	"github.com/MaximBayurov/rate-limiter/internal/logger"
	"github.com/MaximBayurov/rate-limiter/internal/server/handlers"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type HTTPServer struct {
	srv     *http.Server
	logger  logger.Logger
	app     app.App
	configs configuration.ServerConf
}

func New(logger logger.Logger, app app.App, configs configuration.ServerConf) Server {
	return HTTPServer{
		srv: &http.Server{
			ReadHeaderTimeout: 10 * time.Second,
		},
		logger:  logger,
		app:     app,
		configs: configs,
	}
}

func (s HTTPServer) Start(_ context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /ip/list", handlers.AddIPHandler(s.app, s.logger))
	mux.HandleFunc("DELETE /ip/list", handlers.DeleteIPHandler(s.app, s.logger))

	handler := loggingMiddleware(&s.logger)(mux)

	addr := fmt.Sprintf("%s:%d", s.configs.Host, s.configs.Port)
	s.logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))

	s.srv.Addr = addr
	s.srv.Handler = handler
	return s.srv.ListenAndServe()
}

func (s HTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
