package web

import (
	_ "embed"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed assets/index.html
var indexHTML []byte

func (ws *WebServer) routes(downloadFolder string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", index)
	mux.HandleFunc("GET /clients/{$}", ws.clientList)
	mux.HandleFunc("GET /messages/{$}", ws.lastMessages)
	mux.HandleFunc("GET /plugs/{$}", ws.plugs)
	mux.HandleFunc("GET /files/{$}", ws.fileList)
	mux.HandleFunc("GET /download/", makeDownloadHandler(downloadFolder))
	return mux
}

func makeDownloadHandler(downloadFolder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find the path to the file from the URL
		path := r.URL.Path[len("/download/"):]
		// sanitize the path
		if strings.Contains(path, "..") {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		// Serve the file
		fullPath := filepath.Join(downloadFolder, path)
		http.ServeFile(w, r, fullPath)
	}
}
