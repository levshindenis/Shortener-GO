package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/levshindenis/sprint1/cmd/config"
	"github.com/levshindenis/sprint1/cmd/funcs"
)

func PostHandler(storage *Storage, sa *config.ServerAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "There is not true method", http.StatusBadRequest)
			return
		}
		body, _ := io.ReadAll(r.Body)
		if err := r.Body.Close(); err != nil {
			return
		}

		if _, err := url.ParseRequestURI(string(body)); err != nil {
			http.Error(w, "There is not url", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		myAddress := sa.GetShortBaseURL() + "/"
		if value, ok := (*storage)[string(body)]; ok {
			if _, err := w.Write([]byte(myAddress + value)); err != nil {
				return
			}
		} else {
			shortKey := funcs.GenerateShortKey()
			for {
				if _, in := (*storage)[shortKey]; !in {
					(*storage)[shortKey] = string(body)
					break
				}
				shortKey = funcs.GenerateShortKey()
			}
			if _, err := w.Write([]byte(myAddress + shortKey)); err != nil {
				return
			}
		}
	}
}

func GetHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "There is not true method", http.StatusBadRequest)
		}
		//r.URL.Path[1:]
		//chi.URLParam(r, "id")
		if _, in := (*storage)[r.URL.Path[1:]]; in {
			w.Header().Add("Location", (*storage)[r.URL.Path[1:]])
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "There is no such shortUrl", http.StatusBadRequest)
		}
	}
}
