package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"url-shortener/internal/server/handlers"
)

func NewRouter(handler *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handler.ShortenURL).Methods(http.MethodPost)
	router.HandleFunc("/{shortURL}", handler.GetOriginalURL).Methods(http.MethodGet)

	return router
}
