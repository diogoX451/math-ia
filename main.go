package main

import (
	"log"
	"math-ia/internal/ia/ollama"
	"math-ia/internal/ia/vectorstore"
	"math-ia/internal/ia/vectorstore/config"
	"math-ia/internal/router"
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

	r := router.NewRouter(ollama)
	log.Println("Starting server on :8081")
	http.ListenAndServe(":8081", r)
}
