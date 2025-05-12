package operations

import (
	"context"
	"fmt"
	"log"
	"math-ia/internal/ia/vectorstore"
)

type Operations struct {
	instance *vectorstore.Milvus
}

func NewOperations(instance *vectorstore.Milvus) *Operations {
	return &Operations{
		instance: instance,
	}
}

func (o *Operations) ListCollections(ctx context.Context) {
	collections, err := o.instance.Client.ListCollections(ctx)

	if err != nil {
		log.Fatalf("erro ao listar coleções: %v", err)
	}

	if len(collections) == 0 {
		println("Nenhuma coleção encontrada.")
		return
	}

	for _, coll := range collections {
		fmt.Println("Coleção:", coll.Name)
	}
}
