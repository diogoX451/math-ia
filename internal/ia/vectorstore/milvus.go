package vectorstore

import (
	"context"
	"fmt"
	"math-ia/internal/ia/vectorstore/config"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Milvus struct {
	config *config.MilvusConfig
	Client client.Client
}

type SearchResult struct {
	Text    string
	Content string
	Source  string
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

func (m *Milvus) InsertVector(ctx context.Context, text string, vector []float32, metadata map[string]string) error {
	collection := "chunks"

	textField := entity.NewColumnVarChar("text", []string{text})
	vectorField := entity.NewColumnFloatVector("embedding", len(vector), [][]float32{vector})
	sourceField := entity.NewColumnVarChar("source", []string{metadata["source"]})

	_, err := m.Client.Insert(ctx, collection, "", textField, vectorField, sourceField)
	return err
}

func (m *Milvus) UpsertVector(ctx context.Context, id int64, question, content string, vector []float32, metadata map[string]string) error {
	collection := "docs"

	textField := entity.NewColumnVarChar("text", []string{question})
	contentField := entity.NewColumnVarChar("content", []string{content})
	sourceField := entity.NewColumnVarChar("source", []string{metadata["source"]})
	vectorField := entity.NewColumnFloatVector("embedding", len(vector), [][]float32{vector})

	_, err := m.Client.Insert(ctx, collection, "", textField, contentField, sourceField, vectorField)
	return err
}

func (m *Milvus) CreateIndexIfNotExists(ctx context.Context, collectionName, vectorField string) error {
	index, err := entity.NewIndexFlat(entity.L2)
	if err != nil {
		return fmt.Errorf("erro ao criar índice: %w", err)
	}

	err = m.Client.CreateIndex(ctx, collectionName, vectorField, index, false)
	if err != nil {
		return fmt.Errorf("erro ao aplicar índice: %w", err)
	}

	return nil
}

func (m *Milvus) CreateCollectionIfNotExists(ctx context.Context, collectionName string, dim int) error {
	has, err := m.Client.HasCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência da collection: %w", err)
	}
	if has {
		return nil
	}

	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "Collection para armazenar chunks e embeddings",
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "text",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "1024",
				},
			},
			{
				Name:     "content",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "4096",
				},
			},
			{
				Name:     "source",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "256",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", dim),
				},
			},
		},
	}

	err = m.Client.CreateCollection(ctx, schema, 2)
	if err != nil {
		return fmt.Errorf("erro ao criar a collection: %w", err)
	}

	return nil
}

func (m *Milvus) SearchSimilar(ctx context.Context, query []float32, topK int) ([]SearchResult, error) {
	outputFields := []string{"text", "content", "source"}
	queryVectors := []entity.Vector{entity.FloatVector(query)}
	searchParam, _ := entity.NewIndexFlatSearchParam()

	results, err := m.Client.Search(
		ctx,
		"docs",
		[]string{},
		"",
		outputFields,
		queryVectors,
		"embedding",
		entity.L2,
		topK,
		searchParam,
	)
	if err != nil {
		return nil, err
	}

	var out []SearchResult
	for _, hit := range results {
		fields := map[string]*entity.ColumnVarChar{}
		for _, field := range hit.Fields {
			switch col := field.(type) {
			case *entity.ColumnVarChar:
				switch col.Name() {
				case "text":
					fields["text"] = col
				case "content":
					fields["content"] = col
				case "source":
					fields["source"] = col
				}
			}
		}

		textCol := fields["text"]
		contentCol := fields["content"]
		sourceCol := fields["source"]

		for i := 0; i < textCol.Len(); i++ {
			text, _ := textCol.ValueByIdx(i)
			content, _ := contentCol.ValueByIdx(i)
			source, _ := sourceCol.ValueByIdx(i)

			out = append(out, SearchResult{
				Text:    text,
				Content: content,
				Source:  source,
			})
		}
	}

	return out, nil
}
