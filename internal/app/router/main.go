package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/middleware"
)

func MyRouter(hs *handlers.HStorage) *chi.Mux {
	var ml middleware.MyLogger
	ml.Init()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", middleware.WithCompression(ml.WithLogging(middleware.CheckCookie(hs.SetLongURL, hs))))
		r.Get("/ping", middleware.WithCompression(ml.WithLogging(hs.Ping)))
		r.Get("/{id}", middleware.WithCompression(ml.WithLogging(hs.GetLongURL)))
		r.Get("/debug/pprof", pprof.Index)
		r.Get("/debug/pprof/profile", pprof.Profile)
		r.Get("/debug/pprof/cmdline", pprof.Cmdline)
		r.Get("/debug/pprof/symbol", pprof.Symbol)
		r.Get("/debug/pprof/trace", pprof.Trace)
		r.Method(http.MethodGet, "/debug/pprof/block", pprof.Handler("block"))
		r.Method(http.MethodGet, "/debug/pprof/heap", pprof.Handler("heap"))

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", middleware.WithCompression(ml.WithLogging(middleware.CheckCookie(hs.SetJSONLongURL, hs))))
			r.Post("/shorten/batch", middleware.WithCompression(ml.WithLogging(middleware.CheckCookie(hs.BatchURLs, hs))))
			r.Get("/user/urls", middleware.WithCompression(ml.WithLogging(middleware.CheckCookie(hs.GetURLs, hs))))
			r.Delete("/user/urls", ml.WithLogging(middleware.CheckCookie(hs.DelURLs, hs)))
		})
	})
	return r
}
