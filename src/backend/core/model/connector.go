package model

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

const (
	CollectionTenant = "tenant_%s"
	CollectionUser   = "user_%s"

	ConnectorStatusActive        = "Active"
	ConnectorStatusPending       = "Pending"
	ConnectorStatusWorking       = "Working"
	ConnectorStatusSuccess       = "Completed Successfully"
	ConnectorStatusError         = "Completed with Errors"
	ConnectorStatusDisabled      = "Disabled"
	ConnectorStatusUnableProcess = "Unable to Process"
)

type Connector struct {
	tableName               struct{}             `pg:"connectors"`
	ID                      decimal.Decimal      `json:"id,omitempty"`
	CredentialID            decimal.NullDecimal  `json:"credential_id,omitempty"`
	Name                    string               `json:"name,omitempty"`
	Type                    SourceType           `json:"source,omitempty" pg:"type"`
	ConnectorSpecificConfig JSONMap              `json:"connector_specific_config,omitempty"`
	RefreshFreq             int                  `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID            `json:"user_id,omitempty"`
	TenantID                uuid.NullUUID        `json:"tenant_id,omitempty"`
	Disabled                bool                 `json:"disabled,omitempty" pg:",use_zero"`
	LastSuccessfulIndexDate pg.NullTime          `json:"last_successful_index_date,omitempty" pg:",use_zero"`
	LastAttemptStatus       string               `json:"last_attempt_status,omitempty"`
	TotalDocsIndexed        int                  `json:"total_docs_indexed" pg:",use_zero"`
	CreationDate            time.Time            `json:"creation_date,omitempty"`
	LastUpdate              pg.NullTime          `json:"last_update,omitempty" pg:",use_zero"`
	DeletedDate             pg.NullTime          `json:"deleted_date,omitempty" pg:",use_zero"`
	Credential              *Credential          `json:"credential,omitempty" pg:"rel:has-one,fk:credential_id"`
	Docs                    []*Document          `json:"docs,omitempty" pg:"rel:has-many"`
	DocsMap                 map[string]*Document `json:"docs_map,omitempty" pg:"-"`
	EmbeddingModel          *EmbeddingModel      `json:"-" pg:"rel:has-one"`
}

func (c *Connector) CollectionName() string {
	return CollectionName(c.UserID, c.TenantID)
}
func CollectionName(userID uuid.UUID, tenantID uuid.NullUUID) string {
	if tenantID.Valid {
		return strings.ReplaceAll(fmt.Sprintf(CollectionTenant, tenantID), "-", "")
	}
	return strings.ReplaceAll(fmt.Sprintf(CollectionUser, userID), "-", "")

}
