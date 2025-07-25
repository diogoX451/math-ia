package selector

import (
	"strings"
)

var Models = map[string]string{
	"brunoconterato/Gemma-3-Gaia-PT-BR-4b-it:f16": "brunoconterato/Gemma-3-Gaia-PT-BR-4b-it:f16",
	"codellama":        "codellama",
	"nomic-embed-text": "nomic-embed-text",
	"hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF": "hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF",
}

var WeightedKeywordsByModel = map[string]map[string]int{
	"hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF": {
		"integral": 3, "derivada": 3, "cálculo": 2, "equação": 2,
		"raiz": 2, "soma": 1, "multiplicação": 1, "divisão": 1, "função": 2,
	},
	"codellama": {
		"python": 3, "golang": 3, "typescript": 3, "api": 2, "variável": 1,
	},
}

func SelectModel(prompt string) string {
	promptLower := strings.ToLower(prompt)

	scores := map[string]int{}
	for model, keywords := range WeightedKeywordsByModel {
		for kw, weight := range keywords {
			if strings.Contains(promptLower, kw) {
				scores[model] += weight
			}
		}
	}

	bestModel := "brunoconterato/Gemma-3-Gaia-PT-BR-4b-it:f16"
	highest := 0
	for model, score := range scores {
		if score > highest {
			bestModel = model
			highest = score
		}
	}

	return bestModel
}
