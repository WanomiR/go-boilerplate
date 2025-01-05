package app

import (
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(a.zapLogger)
	r.Use(a.rateLimiter)
	r.Use(a.requestsCounter)

	r.Get("/hello", a.http.HelloWorld)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", strings.Split(a.config.Target.Addr, ":")[1])),
	))

	return r
}
