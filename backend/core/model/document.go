package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type Document struct {
	tableName        struct{}    `pg:"documents"`
	ID               int64       `json:"id,omitempty"`
	DocumentID       string      `json:"document_id,omitempty"`
	ConnectorID      int64       `json:"connector_id,omitempty"`
	Boost            int         `json:"boost,omitempty"`
	Hidden           bool        `json:"hidden,omitempty"`
	SemanticID       string      `json:"semantic_id,omitempty"`
	Link             string      `json:"link,omitempty"`
	FromIngestionAPI bool        `json:"from_ingestion_api,omitempty"`
	Signature        string      `json:"signature,omitempty"`
	CreatedDate      time.Time   `json:"created_date,omitempty"`
	UpdatedDate      pg.NullTime `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate      pg.NullTime `json:"deleted_date,omitempty" pg:",use_zero"`
}

type DocumentFeedback struct {
	tableName    struct{}  `pg:"document_feedbacks"`
	ID           int64     `json:"id,omitempty"`
	DocumentID   int64     `json:"document_id,omitempty"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	DocumentRank int       `json:"document_rank,omitempty" pg:",use_zero"`
	UpVotes      bool      `json:"up_votes,omitempty" pg:",use_zero"`
	Feedback     string    `json:"feedback,omitempty" pg:",use_zero"`
}
