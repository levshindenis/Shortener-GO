// пакеты исполняемых приложений должны называться main
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/levshindenis/sprint1/internal/app/handlers"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	var server handlers.API
	server.Init()

	return http.ListenAndServe(server.GetStartSA(), MyRouter(server))
}

func MyRouter(api handlers.API) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", api.PostHandler)
		r.Get("/{id}", api.GetHandler)
	})
	return r
}
