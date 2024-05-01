package model

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Persona struct {
	tableName       struct{}        `pg:"personas"`
	ID              decimal.Decimal `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	LlmID           decimal.Decimal `json:"llm_id,omitempty"`
	DefaultPersona  bool            `json:"default_persona,omitempty" pg:",use_zero"`
	Description     string          `json:"description,omitempty" pg:",use_zero"`
	TenantID        uuid.UUID       `json:"tenant_id,omitempty"`
	IsVisible       bool            `json:"is_visible,omitempty" pg:",use_zero"`
	DisplayPriority int             `json:"display_priority,omitempty"`
	StarterMessages json.RawMessage `json:"starter_messages,omitempty" pg:",use_zero"`
	LLM             *LLM            `json:"llm,omitempty" pg:"rel:has-one"`
	Prompt          *Prompt         `json:"prompt,omitempty" pg:"rel:has-one,fk:id,join_fk:persona_id"`
	CreatedDate     time.Time       `json:"created_date,omitempty"`
	UpdatedDate     pg.NullTime     `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate     pg.NullTime     `json:"deleted_date,omitempty" pg:",use_zero"`
}
