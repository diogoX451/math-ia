package main

import (
	"context"
	"log"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/vectorstore"
	"math-ia/internal/ia/vectorstore/config"
	"math-ia/internal/router"
	"math-ia/internal/tools"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	_ = godotenv.Load()

	config := config.NewMilvusConfig(
		os.Getenv("MILVUS_HOST"),
		os.Getenv("MILVUS_PORT"),
	)

	ctx := context.Background()

	milvus, err := vectorstore.NewMilvus(config)
	if err != nil {
		log.Fatal("Erro ao conectar Milvus:", err)
	}
	defer milvus.Client.Close()

	_ = milvus.Client.DropCollection(ctx, "docs")

	err = milvus.CreateCollectionIfNotExists(ctx, "docs", 768)
	if err != nil {
		log.Fatal("Erro ao criar collection:", err)
	}

	err = milvus.CreateIndexIfNotExists(ctx, "docs", "embedding")
	if err != nil {
		log.Fatal("Erro ao criar índice:", err)
	}

	err = milvus.Client.LoadCollection(ctx, "docs", false)
	if err != nil {
		log.Fatal("Erro ao carregar coleção:", err)
	}

	ollama := ollama.NewClient(os.Getenv("OLLAMA_HOST"))

	err = tools.RunIngest(ctx, milvus, ollama, "nomic-embed-text", "./examples")
	if err != nil {
		log.Fatalf("Erro ao ingerir contexto inicial: %v", err)
	}

	r := router.NewRouter(ollama, milvus)
	log.Println("Starting server on :8081")
	http.ListenAndServe(":8081", r)
}
