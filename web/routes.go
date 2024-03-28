package web

import (
	_ "embed"
	"net/http"
)

//go:embed assets/index.html
var indexHTML []byte

func (ws *WebServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", index)
	return mux
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(indexHTML)
}
