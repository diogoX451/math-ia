package ollama

import (
	"context"
	"encoding/json"
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

func (c *Client) Generate(ctx context.Context, model, prompt string, system string) (string, error) {
	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	if system != "" {
		req.System = system
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
