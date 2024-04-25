package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

const (
	SourceTypeIngestionApi   SourceType = "ingestion_api"
	SourceTypeSlack          SourceType = "slack"
	SourceTypeWEB            SourceType = "web"
	SourceTypeGoogleDrive    SourceType = "google_drive"
	SourceTypeGMAIL          SourceType = "gmail"
	SourceTypeRequesttracker SourceType = "requesttracker"
	SourceTypeGithub         SourceType = "github"
	SourceTypeGitlab         SourceType = "gitlab"
	SourceTypeGuru           SourceType = "guru"
	SourceTypeBookstack      SourceType = "bookstack"
	SourceTypeConfluence     SourceType = "confluence"
	SourceTypeSlab           SourceType = "slab"
	SourceTypeJira           SourceType = "jira"
	SourceTypeProductboard   SourceType = "productboard"
	SourceTypeFile           SourceType = "file"
	SourceTypeNotion         SourceType = "notion"
	SourceTypeZulip          SourceType = "zulip"
	SourceTypeLinear         SourceType = "linear"
	SourceTypeHubspot        SourceType = "hubspot"
	SourceTypeDocument360    SourceType = "document360"
	SourceTypeGong           SourceType = "gong"
	SourceTypeGoogleSites    SourceType = "google_sites"
	SourceTypeZendesk        SourceType = "zendesk"
	SourceTypeLoopio         SourceType = "loopio"
	SourceTypeSharepoint     SourceType = "sharepoint"

	CollectionTenant = "tenant:%s"
	CollectionUser   = "user:%s"

	StatusFailed  = "failed"
	StatusSuccess = "success"
)

type SourceType string

type Connector struct {
	tableName               struct{}             `pg:"connectors,omitempty"`
	ID                      int64                `json:"id,omitempty"`
	CredentialID            int64                `json:"credential_id,omitempty"`
	Name                    string               `json:"name,omitempty"`
	Source                  SourceType           `json:"source,omitempty"`
	InputType               string               `json:"input_type,omitempty"`
	ConnectorSpecificConfig JSONMap              `json:"connector_specific_config,omitempty"`
	RefreshFreq             int                  `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID            `json:"user_id,omitempty"`
	TenantID                uuid.UUID            `json:"tenant_id,omitempty"`
	Shared                  bool                 `json:"shared,omitempty" pg:",use_zero"`
	Disabled                bool                 `json:"disabled,omitempty" pg:",use_zero"`
	LastSuccessfulIndexTime pg.NullTime          `json:"last_successful_index_time,omitempty" pg:",use_zero"`
	LastAttemptStatus       string               `json:"last_attempt_status,omitempty"`
	TotalDocsIndexed        int                  `json:"total_docs_indexed" pg:",use_zero"`
	CreatedDate             time.Time            `json:"created_date,omitempty"`
	UpdatedDate             pg.NullTime          `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate             pg.NullTime          `json:"deleted_date,omitempty" pg:",use_zero"`
	Credential              *Credential          `json:"credential,omitempty" pg:"rel:has-one,fk:credential_id"`
	Docs                    []*Document          `json:"docs,omitempty" pg:"rel:has-many"`
	DocsMap                 map[string]*Document `json:"docs_map,omitempty" pg:"-"`
}
