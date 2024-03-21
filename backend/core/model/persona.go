package model

import "github.com/google/uuid"

type Persona struct {
	tableName        struct{} `pg:"personas"`
	Id               int64
	Name             string
	Llm_id           int
	DefaultPersona   bool `pg:",use_zero"`
	Description      string
	Tenant_id        uuid.UUID
	Search_type      string
	Is_visible       bool `pg:",use_zero"`
	Display_priority int
	Starter_messages JSON `pg:",use_zero"`
}
