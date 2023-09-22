package interfaces

import (
	"net/http"

	"backend/internal/api/interfaces/helper"
	v1 "backend/internal/api/interfaces/v1/products"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
}

// NewApp factory returns an instantiated App
func NewApp() *App {
	return &App{}
}

// Router provides api router
func (a *App) Router() http.Handler {
	m := chi.NewRouter()

	m.Use(middleware.Compress(5, "gzip"))
	m.Use(middleware.RealIP)
	m.Use(middleware.RequestID)
	// m.Use(middleware.Logger)

	// status
	m.Get("/", healthCheck)

	// v1 interfaces
	m.Route("/v1/", func(r chi.Router) {
		r.Mount("/products", v1.NewProductResources().Routes())
	})

	return m
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	helper.Succeed(w, struct{}{})
}
