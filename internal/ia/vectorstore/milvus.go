package vectorstore

import (
	"context"
	"fmt"
	"math-ia/internal/ia/vectorstore/config"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Milvus struct {
	config *config.MilvusConfig
	Client client.Client
}

func NewMilvus(cfg *config.MilvusConfig) (*Milvus, error) {
	ctx := context.Background()

	address := cfg.GetURL()

	mc, err := client.NewDefaultGrpcClient(ctx, address)

	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com Milvus: %w", err)
	}

	return &Milvus{
		config: cfg,
		Client: mc,
	}, nil
}
