package router

import (
	"math-ia/internal/api"
	"math-ia/internal/ia/ollama"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(ollamaClient *ollama.Client) http.Handler {
	r := chi.NewRouter()

	handler := api.NewHandler(ollamaClient)

	r.Post("/ask", handler.Ask)

	return r
}
