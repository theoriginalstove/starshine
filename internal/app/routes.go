package app

import (
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", templ.Handler(hello()).ServeHTTP)
	r.Get("/on", h.on)
	r.Get("/off", h.off)
	return r
}
