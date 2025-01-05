package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(a.ZapLogger)
	r.Use(a.RateLimiter)

	r.Get("/hello", a.http.HelloWorld)

	return r
}
