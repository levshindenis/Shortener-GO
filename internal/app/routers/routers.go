package routers

import (
	"github.com/go-chi/chi/v5"

	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/middleware"
)

func MyRouter(hs handlers.HStorage) *chi.Mux {
	var ml middleware.MyLogger
	ml.Init()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", middleware.WithCompression(ml.WithLogging(hs.PostHandler)))
		r.Get("/ping", middleware.WithCompression(ml.WithLogging(hs.GetPingHandler)))
		r.Get("/{id}", middleware.WithCompression(ml.WithLogging(hs.GetHandler)))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", middleware.WithCompression(ml.WithLogging(hs.JSONPostHandler)))
			r.Post("/shorten/batch", middleware.WithCompression(ml.WithLogging(hs.BatchPostHandler)))
			r.Get("/user/urls", middleware.WithCompression(ml.WithLogging(hs.GetURLS)))
		})
	})
	return r
}
