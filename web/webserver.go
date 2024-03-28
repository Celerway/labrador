package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/celerway/labrador/broker"
	"log/slog"
	"net/http"
	"time"
)

type WebServer struct {
	logger *slog.Logger
	server *http.Server
	broker *broker.State
}

func New(addr string, br *broker.State, logger *slog.Logger) *WebServer {

	ws := &WebServer{
		logger: logger,
		broker: br,
	}

	routes := ws.routes()
	mw := LoggingMiddleware(logger)
	routes = mw(routes)

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

// LoggingMiddleware is a middleware function that logs the request path, method, remote address, user agent and duration of the request
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			w = newLoggingResponseWriter(w)
			next.ServeHTTP(w, r)
			// if the request is not for /healthz or /metrics, log the request:
			if r.URL.Path != "/healthz" && r.URL.Path != "/metrics" {
				logger.Info("request", "status", w.(*loggingResponseWriter).statusCode,
					"path", r.URL.Path, "method", r.Method,
					"remote", r.RemoteAddr, "user-agent", r.UserAgent(), "duration_ms", time.Since(start).Milliseconds())
			}
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// UnixToISO - converts an int64 Unix timestamp to an ISO 8601 formatted string.
func UnixToISO(unixTime int64) string {
	// Convert the int64 Unix timestamp to a time.Time
	t := time.Unix(unixTime, 0)
	// Format the time in ISO 8601 format
	return t.Format(time.RFC3339)
}
