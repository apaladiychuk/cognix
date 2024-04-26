package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type EmbeddingModel struct {
	tableName struct{} `pg:"embedding_models"`

	ID            decimal.Decimal `json:"id,omitempty"`
	TenantID      uuid.UUID       `json:"tenant_id,omitempty"`
	ModelID       string          `json:"model_id,omitempty"`
	ModelName     string          `json:"model_name,omitempty"`
	ModelDim      int             `json:"model_dim,omitempty" pg:",use_zero"`
	Normalize     bool            `json:"normalize,omitempty" pg:",use_zero"`
	QueryPrefix   string          `json:"query_prefix,omitempty"`
	PassagePrefix string          `json:"passage_prefix,omitempty"`
	IndexName     string          `json:"index_name,omitempty"`
	URL           string          `json:"url,omitempty"`
	IsActive      bool            `json:"is_active,omitempty" pg:",use_zero"`
	CreatedDate   time.Time       `json:"created_date,omitempty"`
	UpdatedDate   pg.NullTime     `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate   pg.NullTime     `json:"deleted_date,omitempty" pg:",use_zero"`
}
