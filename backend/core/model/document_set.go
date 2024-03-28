package model

import "github.com/google/uuid"

type (
	DocumentSet struct {
		tableName   struct{}  `pg:"document_sets"`
		ID          int       `json:"id,omitempty"`
		UserID      uuid.UUID `json:"user_id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		IsUpToDate  bool      `json:"is_up_to_date,omitempty"`
	}

	DocumentSetConnectorPair struct {
		tableName     struct{} `pg:"document_set_connector_pairs"`
		ID            int      `json:"id,omitempty"`
		DocumentSetID int      `json:"document_set_id,omitempty"`
		ConnectorID   int      `json:"connector_id,omitempty"`
		IsCurrent     bool     `json:"is_current,omitempty"`
	}
)
