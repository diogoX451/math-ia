package router

import (
	"math-ia/internal/api"
	"math-ia/internal/ia/ollama"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		println("Middleware - Origin:", origin)

		allowed := map[string]bool{
			"http://localhost:8080":         true,
			"https://app.swiftstock.com.br": true,
		}

		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func NewRouter(ollamaClient *ollama.Client) http.Handler {
	r := chi.NewRouter()

	r.Use(withCORS)

	handler := api.NewHandler(ollamaClient)

	r.Post("/ask", handler.Ask)

	return r
}
