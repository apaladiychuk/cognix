package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

const (
	IngestionApi   SourceType = "ingestion_api"
	Slack          SourceType = "slack"
	WEB            SourceType = "web"
	GoogleDrive    SourceType = "google_drive"
	GMAIL          SourceType = "gmail"
	Requesttracker SourceType = "requesttracker"
	Github         SourceType = "github"
	Gitlab         SourceType = "gitlab"
	Guru           SourceType = "guru"
	Bookstack      SourceType = "bookstack"
	Confluence     SourceType = "confluence"
	Slab           SourceType = "slab"
	Jira           SourceType = "jira"
	Productboard   SourceType = "productboard"
	File           SourceType = "file"
	Notion         SourceType = "notion"
	Zulip          SourceType = "zulip"
	Linear         SourceType = "linear"
	Hubspot        SourceType = "hubspot"
	Document360    SourceType = "document360"
	Gong           SourceType = "gong"
	GoogleSites    SourceType = "google_sites"
	Zendesk        SourceType = "zendesk"
	Loopio         SourceType = "loopio"
	Sharepoint     SourceType = "sharepoint"
)

type SourceType string

type Connector struct {
	tableName               struct{}    `pg:"connectors,omitempty"`
	ID                      int64       `json:"id,omitempty"`
	CredentialID            int64       `json:"credential_id,omitempty" pg:",use_zero"`
	Name                    string      `json:"name,omitempty"`
	Source                  SourceType  `json:"source,omitempty"`
	InputType               string      `json:"input_type,omitempty"`
	ConnectorSpecificConfig JSONMap     `json:"connector_specific_config,omitempty"`
	RefreshFreq             int         `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID   `json:"user_id,omitempty"`
	TenantID                uuid.UUID   `json:"tenant_id,omitempty"`
	Shared                  bool        `json:"shared,omitempty" pg:",use_zero"`
	Disabled                bool        `json:"disabled,omitempty" pg:",use_zero"`
	LastSuccessfulIndexTime pg.NullTime `json:"last_successful_index_time,omitempty" pg:",use_zero"`
	LastAttemptStatus       string      `json:"last_attempt_status,omitempty"`
	TotalDocsIndexed        int         `json:"total_docs_indexed" pg:",use_zero"`
	CreatedDate             time.Time   `json:"created_date,omitempty"`
	UpdatedDate             pg.NullTime `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate             pg.NullTime `json:"deleted_date,omitempty" pg:",use_zero"`
}
