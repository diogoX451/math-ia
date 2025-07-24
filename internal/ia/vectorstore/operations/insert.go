package operations

import (
	"context"
	"fmt"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/vectorstore"
)

func InsertChunksToMilvus(
	ctx context.Context,
	chunks []string, model string, meta map[string]string,
	instance vectorstore.Milvus,
	ollama *ollama.Client,
) error {
	if len(chunks) == 0 {
		return fmt.Errorf("nenhum chunk fornecido para inserção")
	}

	for _, chunk := range chunks {
		embedding, err := ollama.GenerateEmbedding(ctx, model, chunk)
		if err != nil {
			return fmt.Errorf("erro ao gerar embedding: %w", err)
		}

		err = instance.InsertVector(ctx, chunk, embedding, meta)
		if err != nil {
			return fmt.Errorf("erro ao inserir no Milvus: %w", err)
		}
	}

	return nil
}
