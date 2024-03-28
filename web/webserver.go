package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type WebServer struct {
	logger *slog.Logger
	server *http.Server
}

func New(addr string, logger *slog.Logger) *WebServer {

	ws := &WebServer{
		logger: logger,
	}
	routes := ws.routes()

	// Create a new web server
	ws.server = &http.Server{
		Addr:    addr,
		Handler: routes,
	}
	return ws
}

func (ws *WebServer) Run(ctx context.Context) error {
	// Run the web server
	go func() {
		<-ctx.Done()
		ws.logger.Info("web server shutting down")
		err := ws.server.Shutdown(ctx)
		if err != nil {
			ws.logger.Error("web server shutdown", "error", err)
		}
	}()
	ws.logger.Info("web server listening", "addr", ws.server.Addr)
	err := ws.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("web server ListenAndServe: %w", err)
	}
	return nil
}
