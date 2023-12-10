// пакеты исполняемых приложений должны называться main
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/levshindenis/sprint1/cmd/config"
	"math/rand"
	"net/http"
	"time"
)

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

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	var storage Storage
	storage.EmptyStorage()

	var sa config.ServerAddress
	config.ParseFlags(&sa)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", PostHandler(&storage, &sa))
		r.Get("/{id}", GetHandler(&storage))
	})
	return http.ListenAndServe(sa.GetStartAddress(), r)
}
