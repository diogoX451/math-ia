package ollama

import (
	"context"
	"encoding/json"
	"math-ia/internal/tools"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system,omitempty"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *Client) Generate(ctx context.Context, model string, question string, contextText string) (string, error) {
	newPrompt := tools.BuildPrompt([]string{contextText}, question)

	req := GenerateRequest{
		Model:  model,
		Prompt: newPrompt,
		Stream: false,
	}

	data, err := c.post(ctx, "/api/generate", req)
	if err != nil {
		return "", err
	}

	var res GenerateResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	return res.Response, nil
}
