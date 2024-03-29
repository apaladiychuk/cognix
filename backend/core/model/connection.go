package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type Connector struct {
	tableName               struct{}    `pg:"connectors,omitempty"`
	ID                      int64       `json:"id,omitempty"`
	CredentialID            int64       `json:"credential_id,omitempty" pg:",use_zero"`
	Name                    string      `json:"name,omitempty"`
	Source                  string      `json:"source,omitempty"`
	InputType               string      `json:"input_type,omitempty"`
	ConnectorSpecificConfig JSONMap     `json:"connector_specific_config,omitempty"`
	RefreshFreq             int         `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID   `json:"user_id,omitempty"`
	TenantID                uuid.UUID   `json:"tenant_id,omitempty"`
	Shared                  bool        `json:"shared,omitempty" pg:",use_zero"`
	Disabled                bool        `json:"disabled,omitempty" pg:",use_zero"`
	LastSuccessfulIndexTime pg.NullTime `json:"last_successful_index_time,omitempty" pg:",use_zero"`
	LastAttemptStatus       string      `json:"last_attempt_status,omitempty"`
	TotalDocsIndexed        int         `json:"total_docs_indexed" pg:",use_zero"`
	CreatedDate             time.Time   `json:"created_date,omitempty"`
	UpdatedDate             pg.NullTime `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate             pg.NullTime `json:"deleted_date,omitempty" pg:",use_zero"`
}
