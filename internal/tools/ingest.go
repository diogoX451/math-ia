package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/vectorstore"
	"os"
)

type Document struct {
	ID      int64  `json:"id"`
	Text    string `json:"text"`
	Source  string `json:"source"`
	Content string `json:"content"`
}

func LoadDocumentsFromFile(path string) ([]Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler JSON: %w", err)
	}
	var docs []Document
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON: %w", err)
	}
	return docs, nil
}

func RunIngest(ctx context.Context, milvus *vectorstore.Milvus, ollama *ollama.Client, model, jsonPath string) error {
	docs, err := LoadDocumentsFromFile(jsonPath)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		embedding, err := ollama.GenerateEmbedding(ctx, model, doc.Content)
		if err != nil {
			return fmt.Errorf("erro ao gerar embedding: %w", err)
		}

		err = milvus.UpsertVector(ctx, doc.ID, doc.Text, doc.Content, embedding, map[string]string{
			"source": doc.Source,
		})

		if err != nil {
			return fmt.Errorf("erro ao inserir/upsert: %w", err)
		}
	}

	return nil
}
