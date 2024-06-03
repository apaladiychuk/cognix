package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/fx"
	"strconv"
	"strings"
)

const (
	ColumnNameID         = "id"
	ColumnNameDocumentID = "document_id"
	ColumnNameChunk      = "chunk"
	ColumnNameContent    = "content"
	ColumnNameVector     = "vector"

	VectorDimension = 1536

	IndexStrategyDISKANN   = "DISKANN"
	IndexStrategyAUTOINDEX = "AUTOINDEX"
	IndexStrategyNoIndex   = "NOINDEX"
)

var responseColumns = []string{ColumnNameID, ColumnNameDocumentID, ColumnNameChunk, ColumnNameContent}

type (
	MilvusConfig struct {
		Address       string `env:"MILVUS_URL"`
		MetricType    string `env:"MILVUS_METRIC_TYPE" envDefault:"COSINE"`
		IndexStrategy string `env:"MILVUS_INDEX_STRATEGY" envDefault:"DISKANN"`
	}
	MilvusPayload struct {
		ID         int64     `json:"id"`
		DocumentID int64     `json:"document_id"`
		Chunk      int64     `json:"chunk"`
		Content    string    `json:"content"`
		Vector     []float32 `json:"vector"`
	}
	MilvusClient interface {
		CreateSchema(ctx context.Context, name string) error
		Save(ctx context.Context, collection string, payloads ...*MilvusPayload) error
		Load(ctx context.Context, collection string, vector []float32) ([]*MilvusPayload, error)
		Delete(ctx context.Context, collection string, documentID ...int64) error
	}
	milvusClient struct {
		client        milvus.Client
		MetricType    entity.MetricType
		IndexStrategy string
	}
)

func (c *milvusClient) Delete(ctx context.Context, collection string, documentID ...int64) error {
	ids := make([]string, 0, len(documentID))
	for _, id := range documentID {
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	return c.client.Delete(ctx, collection, "", fmt.Sprintf("document_id in [%s]", strings.Join(ids, ",")))
}

func (v MilvusConfig) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Address, validation.Required),
		validation.Field(&v.IndexStrategy, validation.Required,
			validation.In(IndexStrategyDISKANN, IndexStrategyAUTOINDEX, IndexStrategyNoIndex)),
		validation.Field(&v.MetricType, validation.Required,
			validation.In(string(entity.COSINE), string(entity.L2), string(entity.IP))),
	)
}

func (c *milvusClient) Load(ctx context.Context, collection string, vector []float32) ([]*MilvusPayload, error) {
	vs := []entity.Vector{entity.FloatVector(vector)}
	sp, _ := entity.NewIndexFlatSearchParam()
	result, err := c.client.Search(ctx, collection, []string{}, "", responseColumns, vs, ColumnNameVector, entity.L2, 10, sp)
	if err != nil {
		return nil, err
	}
	var payload []*MilvusPayload
	for _, row := range result {
		var pr MilvusPayload
		if err = pr.FromResult(row); err != nil {
			return nil, err
		}
		payload = append(payload, &pr)
	}
	return payload, nil
}

var MilvusModule = fx.Options(
	fx.Provide(func() (*MilvusConfig, error) {
		cfg := MilvusConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		NewMilvusClient,
	),
)

func NewMilvusClient(cfg *MilvusConfig) (MilvusClient, error) {
	client, err := milvus.NewClient(context.Background(), milvus.Config{
		Address: cfg.Address,
	})
	if err != nil {
		return nil, err
	}
	return &milvusClient{
		client:        client,
		MetricType:    entity.MetricType(cfg.MetricType),
		IndexStrategy: cfg.IndexStrategy,
	}, nil
}

func (c *milvusClient) Save(ctx context.Context, collection string, payloads ...*MilvusPayload) error {
	var ids, documentIDs, chunks []int64
	var contents [][]byte
	var vectors [][]float32

	for _, payload := range payloads {
		ids = append(ids, payload.ID)
		documentIDs = append(documentIDs, payload.DocumentID)
		chunks = append(chunks, payload.Chunk)
		contents = append(contents, []byte(fmt.Sprintf(`{"content":"%s"}`, payload.Content)))
		vectors = append(vectors, payload.Vector)
	}
	if _, err := c.client.Insert(ctx, collection, "",
		entity.NewColumnInt64(ColumnNameID, ids),
		entity.NewColumnInt64(ColumnNameDocumentID, documentIDs),
		entity.NewColumnInt64(ColumnNameChunk, chunks),
		entity.NewColumnJSONBytes(ColumnNameContent, contents),
		entity.NewColumnFloatVector(ColumnNameVector, VectorDimension, vectors),
	); err != nil {
		return err
	}
	return nil
}

func (c *milvusClient) indexStrategy() (entity.Index, error) {
	switch c.IndexStrategy {
	case IndexStrategyAUTOINDEX:
		return entity.NewIndexAUTOINDEX(c.MetricType)
	case IndexStrategyDISKANN:
		return entity.NewIndexDISKANN(c.MetricType)
	}
	return nil, fmt.Errorf("index strategy %s not supported yet", c.IndexStrategy)
}
func (c *milvusClient) CreateSchema(ctx context.Context, name string) error {

	collExists, err := c.client.HasCollection(ctx, name)
	if err != nil {
		return err
	}
	if collExists {
		if err = c.client.DropCollection(ctx, name); err != nil {
			return err
		}
		collExists = false
	}

	if !collExists {
		schema := entity.NewSchema().WithName(name).
			WithField(entity.NewField().WithName(ColumnNameID).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
			WithField(entity.NewField().WithName(ColumnNameDocumentID).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(ColumnNameChunk).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(ColumnNameContent).WithDataType(entity.FieldTypeJSON)).
			WithField(entity.NewField().WithName(ColumnNameVector).WithDataType(entity.FieldTypeFloatVector).WithDim(1536))
		if err = c.client.CreateCollection(ctx, schema, 2, milvus.WithAutoID(true)); err != nil {
			return err
		}

		if c.IndexStrategy != IndexStrategyNoIndex {
			indexStrategy, err := c.indexStrategy()
			if err != nil {
				return err
			}
			if err = c.client.CreateIndex(ctx, name, ColumnNameVector, indexStrategy, true); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *MilvusPayload) FromResult(res milvus.SearchResult) error {
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
