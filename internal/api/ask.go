package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/selector"
	"math-ia/internal/ia/vectorstore"
	"strings"

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
	Vector       *vectorstore.Milvus
}

func NewHandler(client *ollama.Client, vector *vectorstore.Milvus) *Handler {
	return &Handler{OllamaClient: client, Vector: vector}
}

func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	model := selector.SelectModel(req.Prompt)
	println("Selected model:", model)
	resp, err := h.OllamaClient.Generate(context.Background(), model, req.Prompt, "")
	if err != nil {
		http.Error(w, "Erro ao gerar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AskResponse{
		Model:    model,
		Response: resp,
	})
}

func (a *Handler) AskWithContext(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	model := selector.SelectModel(req.Prompt)
	println("Selected model:", model)

	embedding, err := a.OllamaClient.GenerateEmbedding(context.Background(), "nomic-embed-text", req.Prompt)
	if err != nil {
		http.Error(w, "Erro ao gerar embedding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Tamanho do embedding gerado:", len(embedding))

	similarDocs, err := a.Vector.SearchSimilar(context.Background(), embedding, 3)
	if err != nil {
		http.Error(w, "Erro na busca vetorial: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var contextBuilder strings.Builder
	for _, doc := range similarDocs {
		contextBuilder.WriteString(doc.Content)
		contextBuilder.WriteString("\n---\n")
	}
	context := contextBuilder.String()

	answer, err := a.OllamaClient.Generate(r.Context(), model, req.Prompt, context)
	if err != nil {
		http.Error(w, "Erro ao gerar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AskResponse{
		Model:    model,
		Response: answer,
	})
}
