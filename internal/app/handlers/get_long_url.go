package handlers

import (
	"net/http"
)

func (serv *HStorage) GetLongURL(w http.ResponseWriter, r *http.Request) {
	result, isdeleted, err := serv.GetStorage().GetData(r.URL.Path[1:], "key", "")
	if err != nil {
		http.Error(w, "Something bad with GetLongURL", http.StatusBadRequest)
		return
	}
	if result == "" {
		http.Error(w, "There is no such shortUrl", http.StatusBadRequest)
		return
	}
	if isdeleted[0] {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Add("Location", result)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
