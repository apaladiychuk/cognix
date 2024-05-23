package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type (
	DocumentSet struct {
		tableName   struct{}                    `pg:"document_sets"`
		ID          decimal.Decimal             `json:"id,omitempty"`
		UserID      uuid.UUID                   `json:"user_id,omitempty"`
		Name        string                      `json:"name,omitempty"`
		Description string                      `json:"description,omitempty"`
		IsUpToDate  bool                        `json:"is_up_to_date,omitempty" pg:",use_zero"`
		CreatedDate time.Time                   `json:"created_date,omitempty"`
		UpdatedDate pg.NullTime                 `json:"updated_date,omitempty" pg:",use_zero"`
		DeletedDate pg.NullTime                 `json:"deleted_date,omitempty" pg:",use_zero"`
		Pairs       []*DocumentSetConnectorPair `json:"pairs,omitempty" pg:"rel:has-manu"`
	}

	DocumentSetConnectorPair struct {
		tableName     struct{}        `pg:"document_set_connector_pairs"`
		DocumentSetID decimal.Decimal `json:"document_set_id,omitempty"`
		ConnectorID   decimal.Decimal `json:"connector_id,omitempty"`
		IsCurrent     bool            `json:"is_current,omitempty" pg:",use_zero"`
	}
)
