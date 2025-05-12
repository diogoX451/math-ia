package main

import (
	"context"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/selector"
	"math-ia/internal/ia/vectorstore"
	"math-ia/internal/ia/vectorstore/config"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	config := config.NewMilvusConfig(
		os.Getenv("MILVUS_HOST"),
		os.Getenv("MILVUS_PORT"),
	)

	milvus, err := vectorstore.NewMilvus(config)
	if err != nil {
		panic(err)
	}

	defer milvus.Client.Close()

	ollama := ollama.NewClient(os.Getenv("OLLAMA_HOST"))

	prompt := "Como calcular a raiz quadrada de 16 e depois calcular o logaritmo na base 2 do resultado?"
	model := selector.SelectModel(prompt)

	res, err := ollama.Generate(context.Background(), model, prompt)
	if err != nil {
		panic(err)
	}

	println(res)
}
