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
		r.Post("/", middleware.WithCompression(ml.WithLogging(middleware.WithCookie(hs.PostHandler, hs))))
		r.Get("/ping", middleware.WithCompression(ml.WithLogging(hs.GetPingHandler)))
		r.Get("/{id}", middleware.WithCompression(hs.GetHandler))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", middleware.WithCompression(ml.WithLogging(middleware.WithCookie(hs.JSONPostHandler, hs))))
			r.Post("/shorten/batch", middleware.WithCompression(ml.WithLogging(middleware.WithCookie(hs.BatchPostHandler, hs))))
			r.Get("/user/urls", middleware.WithCompression(ml.WithLogging(middleware.WithCookie(hs.GetURLs, hs))))
			r.Delete("/user/urls", hs.DelURLs)
		})
	})
	return r
}
