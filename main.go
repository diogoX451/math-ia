package main

import (
	"context"
	"log"
	"math-ia/internal/db"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/vectorstore"
	"math-ia/internal/ia/vectorstore/config"
	"math-ia/internal/router"
	"math-ia/internal/tools/ingestor"
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

	urlDatabase := os.Getenv("DATABASE_URL")
	if urlDatabase == "" {
		panic("DATABASE_URL is not set")
	}

	if err := db.Init(urlDatabase); err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}

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

	if err := milvus.CreateCollectionIfNotExists(ctx, "docs", 768); err != nil {
		log.Fatal("Erro ao criar collection:", err)
	}
	if err := milvus.CreateIndexIfNotExists(ctx, "docs", "embedding"); err != nil {
		log.Fatal("Erro ao criar índice:", err)
	}
	if err := milvus.Client.LoadCollection(ctx, "docs", false); err != nil {
		log.Fatal("Erro ao carregar coleção:", err)
	}

	ollama := ollama.NewClient(os.Getenv("OLLAMA_HOST"))

	// // Ingesta do JSON estático
	if err := ingestor.RunIngest(ctx, milvus, ollama, "nomic-embed-text", "context/produzindocerto.json"); err != nil {
		log.Fatalf("Erro ao ingerir contexto do JSON: %v", err)
	}

	// Ingesta dinâmica do banco de dados
	if err := ingestor.RunIngestFromDB(ctx, milvus, ollama, "nomic-embed-text"); err != nil {
		log.Fatalf("Erro ao ingerir contexto do banco de dados: %v", err)
	}

	r := router.NewRouter(ollama, milvus)
	log.Println("Starting server on :8081")
	http.ListenAndServe(":8081", r)
}
