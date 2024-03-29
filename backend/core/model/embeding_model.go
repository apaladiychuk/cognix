package model

import "github.com/google/uuid"

type EmbeddingModel struct {
	tableName struct{} `pg:"embedding_models"`

	ID            int64     `json:"id,omitempty"`
	TenantID      uuid.UUID `json:"tenant_id,omitempty"`
	ModelID       string    `json:"model_id,omitempty"`
	ModelName     string    `json:"model_name,omitempty"`
	ModelDim      int       `json:"model_dim,omitempty"`
	Normalize     bool      `json:"normalize,omitempty"`
	QueryPrefix   string    `json:"query_prefix,omitempty"`
	PassagePrefix string    `json:"passage_prefix,omitempty"`
	IndexName     string    `json:"index_name,omitempty"`
	URL           string    `json:"url,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
}
