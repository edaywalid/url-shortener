package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/edaywalid/url-shortner/internal/config"
	"github.com/edaywalid/url-shortner/internal/models"
	"github.com/edaywalid/url-shortner/internal/services"
)

type Handler struct {
	service *services.Service
	config  *config.Config
}

func NewHandler(service *services.Service, config *config.Config) *Handler {
	return &Handler{service, config}
}

func (h *Handler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req models.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	shortCode, err := h.service.GetShortURL(ctx, &req)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	resp := models.Response{ShortURL: h.config.BaseURL + shortCode}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpError(w, err, http.StatusInternalServerError)
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := r.URL.Path[1:]

	url, err := h.service.GetURL(ctx, shortCode)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func httpError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
