package http

import (
	"net/http"
	"reactiveNews/tickers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() {
	r := chi.NewRouter()
	r.Use(
		middleware.AllowContentType("application/json"),
		middleware.Logger,
		middleware.Recoverer,
	)
}

type Handler struct {
	service tickers.Service
}

func NewHandler(router chi.Router, service tickers.Service) Handler {
	h := Handler{
		service: service,
	}

	router.Route("/ticker", func(r chi.Router) {
		r.Put("/{tcname}", h.AddTicker)
		r.Delete("/{tcname}", h.RemoveTicker)
	})

	return h
}

func (h *Handler) AddTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "tcname")
	if ticker == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.Add(r.Context(), ticker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "tcname")
	if ticker == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.Remove(r.Context(), ticker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
