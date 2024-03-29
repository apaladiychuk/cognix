package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type Prompt struct {
	tableName        struct{}    `pg:"prompts"`
	ID               int64       `json:"id,omitempty"`
	PersonaID        int64       `json:"persona_id,omitempty"`
	UserID           uuid.UUID   `json:"user_id,omitempty"`
	Name             string      `json:"name,omitempty"`
	Description      string      `json:"description,omitempty"`
	SystemPrompt     string      `json:"system_prompt,omitempty"`
	TaskPrompt       string      `json:"task_prompt,omitempty"`
	IncludeCitations bool        `json:"include_citations,omitempty"`
	DatetimeAware    bool        `json:"datetime_aware,omitempty"`
	DefaultPrompt    bool        `json:"default_prompt,omitempty"`
	CreatedDate      time.Time   `json:"created_date,omitempty"`
	DeletedDate      pg.NullTime `json:"deleted_date,omitempty"`
}
