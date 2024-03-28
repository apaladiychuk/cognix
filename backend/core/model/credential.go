package model

import (
	"github.com/google/uuid"
	null "gopkg.in/guregu/null.v4"
	"time"
)

type Credential struct {
	tableName      struct{}  `pg:"credentials"`
	ID             int       `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	Source         string    `json:"source"`
	CreatedDate    time.Time `json:"created_date"`
	UpdatedDate    null.Time `json:"updated_date" pg:",use_zero"`
	DeletedDate    null.Time `json:"deleted_date" pg:",use_zero"`
	Shared         bool      `json:"shared" pg:",use_zero"`
	CredentialJson JSONMap   `json:"credential_json" pg:"type:jsonb"`
}
