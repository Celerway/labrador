package web

import (
	_ "embed"
	"net/http"
)

//go:embed assets/index.html
var indexHTML []byte

func (ws *WebServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", index)
	mux.HandleFunc("GET /clients/{$}", ws.clientList)
	mux.HandleFunc("GET /messages/{$}", ws.lastMessages)
	return mux
}
