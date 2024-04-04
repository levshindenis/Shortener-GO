package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/levshindenis/sprint1/internal/app/tools"
)

func (serv *HStorage) SetLongURL(w http.ResponseWriter, r *http.Request) {
	var (
		body []byte
		err  error
	)

	cookie, _ := r.Cookie("UserID")
	http.SetCookie(w, &http.Cookie{Name: "UserID", Value: cookie.Value})

	if r.Header.Get("Content-Type") == "application/x-gzip" {
		body, err = tools.Unpacking(r.Body)
		if err != nil {
			http.Error(w, "Something bad with compression", http.StatusBadRequest)
			return
		}
	} else {
		body, _ = io.ReadAll(r.Body)
		if _, err = url.ParseRequestURI(string(body)); err != nil {
			http.Error(w, "There is not url", http.StatusBadRequest)
			return
		}
	}
	defer r.Body.Close()

	address, flag, err := serv.MakeShortURL(string(body))
	if err != nil {
		http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
		return
	}

	if flag {
		w.WriteHeader(http.StatusConflict)
	}

	w.WriteHeader(http.StatusCreated)
	if err = serv.GetStorage().SetData(address, string(body), cookie.Value); err != nil {
		http.Error(w, "Something bad with Save", http.StatusBadRequest)
		return
	}

	address = serv.GetServerConfig("baseURL") + "/" + address
	if _, err = w.Write([]byte(address)); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}
