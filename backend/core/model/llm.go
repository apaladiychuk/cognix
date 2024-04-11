package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type LLM struct {
	tableName   struct{}    `pg:"llm"`
	ID          int64       `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	ModelID     string      `json:"model_id,omitempty"`
	TenantID    uuid.UUID   `json:"tenant_id,omitempty"`
	Url         string      `json:"url,omitempty"`
	ApiKey      string      `json:"-"`
	Endpoint    string      `json:"endpoint,omitempty"`
	CreatedDate time.Time   `json:"created_date,omitempty"`
	UpdatedDate pg.NullTime `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate pg.NullTime `json:"deleted_date,omitempty" pg:",use_zero"`
}
