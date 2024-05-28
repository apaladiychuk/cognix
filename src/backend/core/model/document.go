package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Document struct {
	tableName   struct{}        `pg:"documents"`
	ID          decimal.Decimal `json:"id,omitempty"`
	DocumentID  string          `json:"document_id,omitempty"`
	ConnectorID decimal.Decimal `json:"connector_id,omitempty"`
	Link        string          `json:"link,omitempty" pg:"link"`
	Signature   string          `json:"signature,omitempty" pg:",use_zero"`
	CreatedDate time.Time       `json:"created_date,omitempty"`
	UpdatedDate pg.NullTime     `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate pg.NullTime     `json:"deleted_date,omitempty" pg:",use_zero"`
	IsExists    bool            `json:"-" pg:"-"`
}

type DocumentResponse struct {
	ID          decimal.Decimal `json:"id,omitempty"`
	MessageID   decimal.Decimal `json:"message_id,omitempty"`
	Link        string          `json:"link,omitempty"`
	DocumentID  string          `json:"document_id,omitempty"`
	Content     string          `json:"content,omitempty"`
	UpdatedDate time.Time       `json:"updated_date,omitempty"`
}
type DocumentFeedback struct {
	tableName    struct{}        `pg:"document_feedbacks"`
	ID           decimal.Decimal `json:"id,omitempty"`
	DocumentID   decimal.Decimal `json:"document_id,omitempty"`
	UserID       uuid.UUID       `json:"user_id,omitempty"`
	DocumentRank int             `json:"document_rank,omitempty" pg:",use_zero"`
	UpVotes      bool            `json:"up_votes,omitempty" pg:",use_zero"`
	Feedback     string          `json:"feedback,omitempty" pg:",use_zero"`
}
