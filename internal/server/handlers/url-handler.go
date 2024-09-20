package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
}

type Handler struct {
	urlService URLService
	log        *slog.Logger
}

func NewHandler(urlService URLService, log *slog.Logger) *Handler {
	return &Handler{
		urlService: urlService,
		log:        log,
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginalURL string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	h.log.Info("Received request to shorten URL", slog.String("original_url", req.OriginalURL))

	shortURL, err := h.urlService.CreateShortURL(r.Context(), req.OriginalURL)
	if err != nil {
		http.Error(w, "Error creating short URL", http.StatusInternalServerError)
		return
	}

	h.log.Info("Short URL created", slog.String("short_url", shortURL))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shortURL"]
	if shortURL == "" {
		http.Error(w, "URL not provided", http.StatusBadRequest)
		return
	}

	h.log.Info("Received request to get original URL", slog.String("short_url", shortURL))

	originalURL, err := h.urlService.GetOriginalURL(r.Context(), shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	h.log.Info("Original URL found", slog.String("original_url", originalURL))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}
