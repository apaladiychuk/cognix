package model

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/oauth2"
	"time"
)

const (
	ProviderCustom    OAuthProvider = "custom"
	ProviderMicrosoft OAuthProvider = "microsoft"
)

var ConnectorAuthProvider = map[SourceType]OAuthProvider{
	SourceTypeOneDrive: ProviderMicrosoft,
	SourceTypeMsTeams:  ProviderMicrosoft,
}

// OAuthProvider represents enum for oauth providers
type OAuthProvider string

type Credential struct {
	tableName      struct{}        `pg:"credentials"`
	ID             decimal.Decimal `json:"id"`
	UserID         uuid.UUID       `json:"user_id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	Source         SourceType      `json:"source"`
	CreatedDate    time.Time       `json:"created_date"`
	UpdatedDate    pg.NullTime     `json:"updated_date" pg:",use_zero"`
	DeletedDate    pg.NullTime     `json:"deleted_date" pg:",use_zero"`
	Shared         bool            `json:"shared" pg:",use_zero"`
	CredentialJson *CredentialJson `json:"credential_json" pg:"type:jsonb"`
	Connectors     []*Connector    `json:"connectors" pg:"rel:has-many"`
}

type CredentialJson struct {
	Provider OAuthProvider `json:"provider"`
	Token    *oauth2.Token `json:"token"`
	Custom   JSONMap       `json:"custom"`
}

func (j *CredentialJson) ToStruct(dest interface{}) error {
	buf, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, dest)
}

func (j *CredentialJson) FromStruct(src interface{}) error {
	buf, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, j)
}
