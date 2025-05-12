package api

import (
	"context"
	"encoding/json"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/selector"

	"net/http"
)

type AskRequest struct {
	Prompt string `json:"prompt"`
}

type AskResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
}

type Handler struct {
	OllamaClient *ollama.Client
}

func NewHandler(client *ollama.Client) *Handler {
	return &Handler{OllamaClient: client}
}

func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	model := selector.SelectModel(req.Prompt)
	println("Selected model:", model)
	resp, err := h.OllamaClient.Generate(context.Background(), model, req.Prompt)
	if err != nil {
		http.Error(w, "Erro ao gerar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AskResponse{
		Model:    model,
		Response: resp,
	})
}
