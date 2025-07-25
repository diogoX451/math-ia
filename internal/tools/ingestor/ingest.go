package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"math-ia/internal/db/loader"
	"math-ia/internal/db/models"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/vectorstore"
	"os"
	"strings"
)

const maxContentLength = 4096

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

func GenerateDocsFromEntity(entity string) ([]Document, error) {
	ids, err := loader.GetLast10IDs(entity)
	if err != nil {
		return nil, err
	}

	var docs []Document
	for _, id := range ids {
		ctxLines, err := loader.GetEntityContext(entity, id)
		if err != nil {
			continue
		}

		text := fmt.Sprintf("Informações do checklist %d", id)
		content := strings.Join(ctxLines, "\n")

		docs = append(docs, Document{
			ID:      id,
			Text:    text,
			Content: content,
			Source:  fmt.Sprintf("%s:%d", entity, id),
		})
	}

	return docs, nil
}

func RunIngestFromDB(ctx context.Context, milvus *vectorstore.Milvus, ollama *ollama.Client, model string) error {
	for entity := range models.Registry {
		ids, err := loader.GetLast10IDs(entity)
		if err != nil {
			fmt.Printf("Erro ao buscar IDs da entidade %s: %v\n", entity, err)
			continue
		}

		for _, id := range ids {
			visited := make(map[string]bool)
			contextLines, err := loader.ContextRecursivo(entity, id, visited)
			if err != nil {
				fmt.Printf("Erro ao montar contexto para %s:%d: %v\n", entity, id, err)
				continue
			}

			content := strings.Join(contextLines, "\n")
			embedding, err := ollama.GenerateEmbedding(ctx, model, content)
			if err != nil {
				fmt.Printf("Erro ao gerar embedding para %s:%d: %v\n", entity, id, err)
				continue
			}

			if len(content) > maxContentLength {
				content = content[:maxContentLength]
			}

			err = milvus.UpsertVector(ctx, id, fmt.Sprintf("Contexto de %s %d", entity, id), content, embedding, map[string]string{
				"source": fmt.Sprintf("%s:%d", entity, id),
			})

			if err != nil {
				fmt.Printf("Erro ao inserir vetor para %s:%d: %v\n", entity, id, err)
				continue
			}

			fmt.Printf("Ingestão completa para %s:%d\n", entity, id)
		}
	}

	return nil
}
