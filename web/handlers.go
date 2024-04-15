package web

import "net/http"

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(indexHTML)
}

func (ws *WebServer) clientList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	clients := ws.broker.CurrentClients()
	component := clientsFragment(clients)
	err := component.Render(r.Context(), w)
	if err != nil {
		ws.logger.Error("clientList", "error", err)
	}
}

func (ws *WebServer) publishMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err := r.ParseForm()
	if err != nil {
		ws.logger.Error("publishMessage", "error", err)
		return
	}
	topic := r.FormValue("topic")
	payload := r.FormValue("payload")
	ws.logger.Info("publishMessage", "topic", topic, "payload", payload)
}

func (ws *WebServer) lastMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	msgs := ws.broker.LastMessages()
	component := messagesFragment(msgs)
	err := component.Render(r.Context(), w)
	if err != nil {
		ws.logger.Error("lastMessages", "error", err)
	}
}

func (ws *WebServer) plugs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	pds := ws.broker.HueBridge.GetPlugs()
	component := plugsFragment(pds)
	err := component.Render(r.Context(), w)
	if err != nil {
		ws.logger.Error("plugs", "error", err)
	}
}

func (ws *WebServer) fileList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	files := ws.downloadList()
	component := filesFragment(files)
	err := component.Render(r.Context(), w)
	if err != nil {
		ws.logger.Error("fileList", "error", err)
	}
}
