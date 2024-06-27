package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const (
	ColumnNameID         = "id"
	ColumnNameDocumentID = "document_id"
	ColumnNameContent    = "content"
	ColumnNameVector     = "vector"

	VectorDimension = 1536

	IndexStrategyDISKANN   = "DISKANN"
	IndexStrategyAUTOINDEX = "AUTOINDEX"
	IndexStrategyNoIndex   = "NOINDEX"
)

var responseColumns = []string{ColumnNameID, ColumnNameDocumentID, ColumnNameContent}

type (
	MilvusConfig struct {
		Address       string `env:"MILVUS_URL,required"`
		Username      string `env:"MILVUS_USERNAME,required"`
		Password      string `env:"MILVUS_PASSWORD,required"`
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
		client     milvus.Client
		cfg        *MilvusConfig
		MetricType entity.MetricType
	}
)

func (c *milvusClient) Delete(ctx context.Context, collection string, documentIDs ...int64) error {
	if err := c.checkConnection(); err != nil {
		return err
	}
	var docsID []string
	for _, id := range documentIDs {
		docsID = append(docsID, strconv.FormatInt(id, 10))
	}
	queryResult, err := c.client.Query(ctx, collection, []string{},
		fmt.Sprintf("document_id in [%s]", strings.Join(docsID, ",")),
		[]string{"id"},
	)
	if err != nil {
		return err
	}
	var ids []string
	for _, result := range queryResult {
		for i := 0; i < result.Len(); i++ {
			if id, err := result.GetAsInt64(i); err == nil {
				ids = append(ids, strconv.FormatInt(id, 10))
			}
		}
	}
	if len(ids) == 0 {
		return c.client.Delete(ctx, collection, "", fmt.Sprintf("id in [%s]", strings.Join(ids, ",")))
	}
	return nil
}

func (v MilvusConfig) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.IndexStrategy, validation.Required,
			validation.In(IndexStrategyDISKANN, IndexStrategyAUTOINDEX, IndexStrategyNoIndex)),
		validation.Field(&v.MetricType, validation.Required,
			validation.In(string(entity.COSINE), string(entity.L2), string(entity.IP))),
	)
}

func (c *milvusClient) Load(ctx context.Context, collection string, vector []float32) ([]*MilvusPayload, error) {
	if err := c.checkConnection(); err != nil {
		return nil, err
	}
	vs := []entity.Vector{entity.FloatVector(vector)}
	sp, _ := entity.NewIndexFlatSearchParam()
	result, err := c.client.Search(ctx, collection, []string{}, "", responseColumns, vs, ColumnNameVector, c.MetricType, 10, sp)
	if err != nil {
		return nil, err
	}
	var payload []*MilvusPayload
	for _, row := range result {
		for i := 0; i < row.ResultCount; i++ {
			var pr MilvusPayload
			if err = pr.FromResult(i, row); err != nil {
				return nil, err
			}
			payload = append(payload, &pr)
		}
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
	client, err := connect(cfg)
	if err != nil {
		zap.S().Errorf("connect to milvus error %s ", err.Error())
	}
	return &milvusClient{
		client:     client,
		cfg:        cfg,
		MetricType: entity.MetricType(cfg.MetricType),
	}, nil
}

func (c *milvusClient) Save(ctx context.Context, collection string, payloads ...*MilvusPayload) error {
	var ids, documentIDs, chunks []int64
	var contents [][]byte
	var vectors [][]float32
	if err := c.checkConnection(); err != nil {
		return err
	}
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
		entity.NewColumnJSONBytes(ColumnNameContent, contents),
		entity.NewColumnFloatVector(ColumnNameVector, VectorDimension, vectors),
	); err != nil {
		return err
	}
	return nil
}

func (c *milvusClient) indexStrategy() (entity.Index, error) {
	switch c.cfg.IndexStrategy {
	case IndexStrategyAUTOINDEX:
		return entity.NewIndexAUTOINDEX(c.MetricType)
	case IndexStrategyDISKANN:
		return entity.NewIndexDISKANN(c.MetricType)
	}
	return nil, fmt.Errorf("index strategy %s not supported yet", c.cfg.IndexStrategy)
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
	schema := entity.NewSchema().WithName(name).
		WithField(entity.NewField().WithName(ColumnNameID).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
		WithField(entity.NewField().WithName(ColumnNameDocumentID).WithDataType(entity.FieldTypeInt64)).
		WithField(entity.NewField().WithName(ColumnNameContent).WithDataType(entity.FieldTypeJSON)).
		WithField(entity.NewField().WithName(ColumnNameVector).WithDataType(entity.FieldTypeFloatVector).WithDim(1536))
	if err = c.client.CreateCollection(ctx, schema, 2, milvus.WithAutoID(true)); err != nil {
		return err
	}

	if c.cfg.IndexStrategy != IndexStrategyNoIndex {
		indexStrategy, err := c.indexStrategy()
		if err != nil {
			return err
		}
		if err = c.client.CreateIndex(ctx, name, ColumnNameVector, indexStrategy, true); err != nil {
			return err
		}
	}
	return nil
}

func (p *MilvusPayload) FromResult(i int, res milvus.SearchResult) error {
	var err error

	for _, field := range res.Fields {
		switch field.Name() {
		case ColumnNameID:
			p.ID, err = field.GetAsInt64(i)
		case ColumnNameDocumentID:
			p.DocumentID, err = field.GetAsInt64(i)
		case ColumnNameContent:
			row, err := field.GetAsString(i)
			if err != nil {
				continue
			}
			contentS := ""
			if err = json.Unmarshal([]byte(row), &contentS); err == nil {
				contentS = strings.ReplaceAll(contentS, "\n", "")
				content := make(map[string]string)
				if err = json.Unmarshal([]byte(contentS), &content); err == nil {
					p.Content = content[ColumnNameContent]
				}
			} else {
				content := make(map[string]string)
				if err = json.Unmarshal([]byte(row), &content); err == nil {
					p.Content = content[ColumnNameContent]
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// checkConnection check is milvus connection is ready
func (c *milvusClient) checkConnection() error {
	// creates connection if not exists
	if c.client == nil {
		client, err := connect(c.cfg)
		if err != nil {
			zap.S().Error(err.Error())
			return fmt.Errorf("milvus is not initialized")
		}
		c.client = client
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// check connection status
	state, err := c.client.CheckHealth(ctx)
	if err != nil {
		return fmt.Errorf("client.CheckHealth error %s", err.Error())
	}
	if !state.IsHealthy {
		return fmt.Errorf("milvus is not ready  %s", strings.Join(state.Reasons, " "))
	}

	return nil
}

func connect(cfg *MilvusConfig) (milvus.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return milvus.NewClient(ctx, milvus.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		RetryRateLimit: &milvus.RetryRateLimitOption{
			MaxRetry:   2,
			MaxBackoff: 2 * time.Second,
		},
	})
}
