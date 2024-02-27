package app

import (
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", templ.Handler(hello()).ServeHTTP)
	return r
}
