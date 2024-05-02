package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/fx"
)

const (
	ColumnNameID         = "id"
	ColumnNameDocumentID = "document_id"
	ColumnNameChunk      = "chunk"
	ColumnNameContent    = "content"
	ColumnNameVector     = "vector"

	VectorDimension = 512
)

var responseColumns = []string{ColumnNameID, ColumnNameDocumentID, ColumnNameChunk, ColumnNameContent}

type (
	Config struct {
		Address string `env:"MILVUS_URL"`
	}
	Payload struct {
		ID         int64     `json:"id"`
		DocumentID int64     `json:"document_id"`
		Chunk      int64     `json:"chunk"`
		Content    string    `json:"content"`
		Vector     []float32 `json:"vector"`
	}
	MilvusClient interface {
		CreateSchema(ctx context.Context, name string) error
		Save(ctx context.Context, collection string, payloads ...*Payload) error
		Load(ctx context.Context, collection string, vector []float32) ([]*Payload, error)
	}
	milvusClient struct {
		client milvus.Client
	}
)

func (c *milvusClient) Load(ctx context.Context, collection string, vector []float32) ([]*Payload, error) {
	vs := []entity.Vector{entity.FloatVector(vector)}
	sp, _ := entity.NewIndexFlatSearchParam()
	result, err := c.client.Search(ctx, collection, []string{}, "", responseColumns, vs, ColumnNameVector, entity.L2, 10, sp)
	if err != nil {
		return nil, err
	}
	var payload []*Payload
	for _, row := range result {
		var pr Payload
		if err = pr.FromResult(row); err != nil {
			return nil, err
		}
		payload = append(payload, &pr)
	}
	return payload, nil
}

var MilvusModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	},
		NewMilvusClient,
	),
)

func NewMilvusClient(cfg *Config) (MilvusClient, error) {
	client, err := milvus.NewClient(context.Background(), milvus.Config{
		Address: cfg.Address,
	})
	if err != nil {
		return nil, err
	}
	return &milvusClient{
		client: client,
	}, nil
}

func (c *milvusClient) Save(ctx context.Context, collection string, payloads ...*Payload) error {
	var ids, documentIDs, chunks []int64
	var contents []string
	var vectors [][]float32

	for _, payload := range payloads {
		ids = append(ids, payload.ID)
		documentIDs = append(documentIDs, payload.DocumentID)
		chunks = append(chunks, payload.Chunk)
		contents = append(contents, payload.Content)
		vectors = append(vectors, payload.Vector)
	}
	if _, err := c.client.Insert(ctx, collection, "",
		entity.NewColumnInt64(ColumnNameID, ids),
		entity.NewColumnInt64(ColumnNameDocumentID, documentIDs),
		entity.NewColumnInt64(ColumnNameChunk, chunks),
		entity.NewColumnString(ColumnNameContent, contents),
		entity.NewColumnFloatVector(ColumnNameVector, VectorDimension, vectors),
	); err != nil {
		return err
	}
	return nil
}

func (c *milvusClient) CreateSchema(ctx context.Context, name string) error {
	collExists, err := c.client.HasCollection(ctx, name)
	if err != nil {
		return err
	}
	if !collExists {
		schema := entity.NewSchema().WithName(name).
			WithField(entity.NewField().WithName(ColumnNameID).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
			WithField(entity.NewField().WithName(ColumnNameDocumentID).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(ColumnNameChunk).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(ColumnNameContent).WithDataType(entity.FieldTypeString)).
			WithField(entity.NewField().WithName(ColumnNameVector).WithDataType(entity.FieldTypeFloatVector).WithDim(8))
		if err = c.client.CreateCollection(ctx, schema, 2, milvus.WithAutoID(true)); err != nil {
			return err
		}
	}
	return nil
}

func (p *Payload) FromResult(res milvus.SearchResult) error {
	var err error
	for _, field := range res.Fields {
		switch field.Name() {
		case ColumnNameID:
			p.ID, err = field.GetAsInt64(0)
		case ColumnNameDocumentID:
			p.DocumentID, err = field.GetAsInt64(0)
		case ColumnNameChunk:
			p.Chunk, err = field.GetAsInt64(0)
		case ColumnNameContent:
			p.Content, err = field.GetAsString(0)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
