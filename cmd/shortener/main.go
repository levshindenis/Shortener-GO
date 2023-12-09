// пакеты исполняемых приложений должны называться main
package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Storage map[string]string

func (storage *Storage) EmptyStorage() {
	*storage = make(map[string]string)
}

// ValueIn проверяет наличие значения в map
func (storage *Storage) ValueIn(s string) string {
	for key, value := range *storage {
		if value == s {
			return key
		}
	}
	return ""
}

// GenerateShortKey генерирует короткий URL
func GenerateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}
	return string(shortKey)
}

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func (storage *Storage) MyRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", storage.PostHandler)
		r.Get("/{id}", storage.GetHandler)
	})

	return r
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	var storage Storage
	storage.EmptyStorage()
	return http.ListenAndServe(`:8080`, storage.MyRouter())
}

func (storage *Storage) PostHandler(w http.ResponseWriter, r *http.Request) {
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

	if key := storage.ValueIn(string(body)); key != "" {
		if _, err := w.Write([]byte("http://localhost:8080/" + key)); err != nil {
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
		if _, err := w.Write([]byte("http://localhost:8080/" + shortkey)); err != nil {
			return
		}
	}
}

func (storage *Storage) GetHandler(w http.ResponseWriter, r *http.Request) {
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
