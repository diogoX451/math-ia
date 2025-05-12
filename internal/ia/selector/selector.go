package selector

import "strings"

var Models = map[string]string{
	"llama3":           "llama3",
	"codellama":        "codellama",
	"nomic-embed-text": "nomic-embed-text",
	"hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF": "hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF",
}

var KeywordsByModel = map[string][]string{
	"hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF": {
		"integral", "derivative", "limit", "cálculo", "matemática", "resolver", "função", "equação",
	},
	"codellama": {
		"typescript", "javascript", "python", "golang", "programar", "código", "API", "variável", "função",
	},
	"nomic-embed-text": {
		"embedding", "vetorial", "similaridade", "recuperação", "documento", "buscar contexto",
	},
}

var WeightedKeywordsByModel = map[string]map[string]int{
	"hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF": {
		"integral": 3, "derivada": 3, "cálculo": 2, "equação": 2,
		"raiz": 2, "soma": 1, "multiplicação": 1, "divisão": 1,
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

	bestModel := "llama3"
	highest := 0
	for model, score := range scores {
		if score > highest {
			bestModel = model
			highest = score
		}
	}

	return bestModel
}
