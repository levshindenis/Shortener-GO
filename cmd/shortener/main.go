// пакеты исполняемых приложений должны называться main
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
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

var storage Storage

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	storage.EmptyStorage()
	return http.ListenAndServe(`:8080`, http.HandlerFunc(ChoiceHandler))
}

func ChoiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		PostHandler(w, r)
	case http.MethodGet:
		GetHandler(w, r)
	default:
		http.Error(w, "Unsupported request method", http.StatusBadRequest)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	if err := r.Body.Close(); err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")

	if key := storage.ValueIn(string(body)); key != "" {
		if _, err := w.Write([]byte("http://localhost:8080/" + key)); err != nil {
			return
		}
	} else {
		shortkey := GenerateShortKey()
		fmt.Println("ShortKey: ", shortkey)
		for {
			if _, in := storage[shortkey]; !in {
				storage[shortkey] = string(body)
				break
			}
			shortkey = GenerateShortKey()
		}
		if _, err := w.Write([]byte("http://localhost:8080/" + shortkey)); err != nil {
			return
		}
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Location", storage[r.URL.Path[1:]])
}
