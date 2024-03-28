package model

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Connector struct {
	tableName               struct{}  `pg:"connectors,omitempty"`
	ID                      int       `json:"id,omitempty"`
	CredentialID            int       `json:"credential_id,omitempty" pg:",use_zero"`
	Name                    string    `json:"name,omitempty"`
	Source                  string    `json:"source,omitempty"`
	InputType               string    `json:"input_type,omitempty"`
	ConnectorSpecificConfig JSONMap   `json:"connector_specific_config,omitempty"`
	RefreshFreq             int       `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID `json:"user_id,omitempty"`
	TenantID                uuid.UUID `json:"tenant_id,omitempty"`
	Shared                  bool      `json:"shared,omitempty" pg:",use_zero"`
	Disabled                bool      `json:"disabled,omitempty" pg:",use_zero"`
	LastSuccessfulIndexTime null.Time `json:"last_successful_index_time,omitempty" pg:",use_zero"`
	LastAttemptStatus       string    `json:"last_attempt_status,omitempty"`
	TotalDocsIndexed        int       `json:"total_docs_indexed"`
	CreatedDate             time.Time `json:"created_date,omitempty"`
	UpdatedDate             null.Time `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate             null.Time `json:"deleted_date,omitempty" pg:",use_zero"`
}
