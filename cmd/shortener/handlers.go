package main

import (
	"github.com/levshindenis/sprint1/cmd/config"
	"io"
	"net/http"
	"net/url"
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
		if key := storage.ValueIn(string(body)); key != "" {
			if _, err := w.Write([]byte(myAddress + key)); err != nil {
				return
			}
		} else {
			shortkey := GenerateShortKey()
			for {
				if _, in := (*storage)[shortkey]; !in {
					(*storage)[shortkey] = string(body)
					break
				}
				shortkey = GenerateShortKey()
			}
			if _, err := w.Write([]byte(myAddress + shortkey)); err != nil {
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
