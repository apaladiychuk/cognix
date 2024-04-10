package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type Credential struct {
	tableName      struct{}    `pg:"credentials"`
	ID             int64       `json:"id"`
	UserID         uuid.UUID   `json:"user_id"`
	TenantID       uuid.UUID   `json:"tenant_id"`
	Source         SourceType  `json:"source"`
	CreatedDate    time.Time   `json:"created_date"`
	UpdatedDate    pg.NullTime `json:"updated_date" pg:",use_zero"`
	DeletedDate    pg.NullTime `json:"deleted_date" pg:",use_zero"`
	Shared         bool        `json:"shared" pg:",use_zero"`
	CredentialJson JSONMap     `json:"credential_json" pg:"type:jsonb"`
}
