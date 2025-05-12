package ollama

import (
	"context"
	"encoding/json"
)

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (c *Client) GenerateEmbedding(ctx context.Context, model, prompt string) ([]float32, error) {
	req := EmbeddingRequest{
		Model:  model,
		Prompt: prompt,
	}
	data, err := c.post(ctx, "/api/embeddings", req)
	if err != nil {
		return nil, err
	}

	var res EmbeddingResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Embedding, nil
}
