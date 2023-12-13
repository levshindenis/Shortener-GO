// пакеты исполняемых приложений должны называться main
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/levshindenis/sprint1/cmd/config"
)

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
