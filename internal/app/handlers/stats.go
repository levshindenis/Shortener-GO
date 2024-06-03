package handlers

import (
	"encoding/json"
	"net/http"
	"slices"
)

// Stats - handler для получения статистики по пользователям и сокращенным URL
func (serv *HStorage) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Real-IP") == "" || !slices.Contains(serv.GetIps(), r.Header.Get("X-Real-IP")) {
		http.Error(w, "Bad IP", http.StatusForbidden)
		return
	}

	data, err := serv.Server.Stats()
	if err != nil {
		http.Error(w, "Something bad with Server Stats", http.StatusBadRequest)
		return
	}

	marsh, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Something bad with Marshal", http.StatusBadRequest)
		return
	}

	if _, err = w.Write(marsh); err != nil {
		http.Error(w, "Something bad with Write", http.StatusBadRequest)
		return
	}
}
