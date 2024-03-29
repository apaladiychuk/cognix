package parameters

import (
	"cognix.ch/api/v2/core/model"
)

type GetAllCredentialsParam struct {
	Source string `query:"source"`
}

type CreateCredentialParam struct {
	Source         string        `json:"source"`
	Shared         bool          `json:"shared"`
	CredentialJson model.JSONMap `json:"credential_json"`
}

type UpdateCredentialParam struct {
	Shared         bool          `json:"shared"`
	CredentialJson model.JSONMap `json:"credential_json"`
}

type CreateConnectorParam struct {
	CredentialID            int64         `json:"credential_id,omitempty"`
	Name                    string        `json:"name,omitempty"`
	Source                  string        `json:"source,omitempty"`
	InputType               string        `json:"input_type,omitempty"`
	ConnectorSpecificConfig model.JSONMap `json:"connector_specific_config,omitempty"`
	RefreshFreq             int           `json:"refresh_freq,omitempty"`
	Shared                  bool          `json:"shared,omitempty"`
	Disabled                bool          `json:"disabled,omitempty"`
}

type UpdateConnectorParam struct {
	CredentialID            int64         `json:"credential_id,omitempty"`
	Name                    string        `json:"name,omitempty"`
	InputType               string        `json:"input_type,omitempty"`
	ConnectorSpecificConfig model.JSONMap `json:"connector_specific_config,omitempty"`
	RefreshFreq             int           `json:"refresh_freq,omitempty"`
	Shared                  bool          `json:"shared,omitempty"`
	Disabled                bool          `json:"disabled,omitempty"`
}
