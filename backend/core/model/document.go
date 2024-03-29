package model

import (
	"time"
)

type Document struct {
	tableName        struct{}  `pg:"documents"`
	ID               int64     `json:"id,omitempty"`
	DocumentID       string    `json:"document_id,omitempty"`
	ConnectorID      int64     `json:"connector_id,omitempty"`
	Boost            int       `json:"boost,omitempty"`
	Hidden           bool      `json:"hidden,omitempty"`
	SemanticID       string    `json:"semantic_id,omitempty"`
	Link             string    `json:"link,omitempty"`
	UpdatedDate      time.Time `json:"updated_date,omitempty"`
	FromIngestionAPI bool      `json:"from_ingestion_api,omitempty"`
	Signature        string    `json:"signature,omitempty"`
}
