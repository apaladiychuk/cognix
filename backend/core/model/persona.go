package model

import "github.com/google/uuid"

type Persona struct {
	tableName       struct{}  `pg:"personas"`
	ID              int64     `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	LlmID           int       `json:"llm_id,omitempty"`
	DefaultPersona  bool      `json:"default_persona,omitempty" pg:",use_zero"`
	Description     string    `json:"description,omitempty"`
	TenantID        uuid.UUID `json:"tenant_id,omitempty"`
	IsVisible       bool      `json:"is_visible,omitempty" pg:",use_zero"`
	DisplayPriority int       `json:"display_priority,omitempty"`
	StarterMessages JSON      `json:"starter_messages,omitempty" pg:",use_zero"`
}
