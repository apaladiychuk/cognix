package model

import (
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
