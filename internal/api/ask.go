package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math-ia/internal/db/loader"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/selector"
	"math-ia/internal/ia/vectorstore"
	"strconv"
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
	log.Println("Selected model:", model)

	embedding, err := a.OllamaClient.GenerateEmbedding(r.Context(), "nomic-embed-text", req.Prompt)
	if err != nil {
		http.Error(w, "Erro ao gerar embedding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Tamanho do embedding gerado:", len(embedding))

	similarDocs, err := a.Vector.SearchSimilar(r.Context(), embedding, 10)
	if err != nil {
		http.Error(w, "Erro na busca vetorial: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var contextBuilder strings.Builder
	visited := make(map[string]bool)

	println("Encontrados", len(similarDocs), "documentos similares")

	for _, doc := range similarDocs {
		println("Processando documento:", doc.Content)
		parts := strings.Split(doc.Source, ":")
		if len(parts) != 2 {
			contextBuilder.WriteString(fmt.Sprintf("Conteúdo: %s\n", doc.Content))
			contextBuilder.WriteString(fmt.Sprintf("Fonte desconhecida: %s\n", doc.Source))
			continue
		}
		entity := parts[0]
		idStr := parts[1]

		if visited[doc.Source] {
			continue
		}
		visited[doc.Source] = true

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}

		relContext, err := loader.GetEntityContext(entity, id)
		if err != nil {
			log.Printf("Erro ao buscar contexto de %s:%d => %v", entity, id, err)
			continue
		}

		contextBuilder.WriteString(fmt.Sprintf("Fonte: %s:%d\n", entity, id))
		for _, line := range relContext {
			contextBuilder.WriteString(line + "\n")
		}
		contextBuilder.WriteString("\n---\n")
	}

	finalContext := contextBuilder.String()

	answer, err := a.OllamaClient.Generate(r.Context(), model, req.Prompt, finalContext)
	if err != nil {
		http.Error(w, "Erro ao gerar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AskResponse{
		Model:    model,
		Response: answer,
	})
}
