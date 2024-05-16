package model

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
	SourceTypeMsTeams        SourceType = "msteams"
)

type (
	SourceType            string
	SourceTypeDescription struct {
		ID            SourceType `json:"id"`
		Name          string     `json:"name"`
		IsImplemented bool       `json:"isImplemented"`
	}
)

var AllSourceTypes = map[SourceType]SourceTypeDescription{
	SourceTypeFile:        {SourceTypeFile, "File", false},
	SourceTypeWEB:         {SourceTypeWEB, "Web", true},
	SourceTypeSlack:       {SourceTypeSlack, "Slack", false},
	SourceTypeGoogleDrive: {SourceTypeGoogleDrive, "Google Drive", false},
	SourceTypeGMAIL:       {SourceTypeGMAIL, "Gmail", false},
	SourceTypeSharepoint:  {SourceTypeSharepoint, "Sharepoint", false},
}
