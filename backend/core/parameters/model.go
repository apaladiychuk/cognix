package parameters

import (
	"cognix.ch/api/v2/core/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type OAuthParam struct {
	Action   string `json:"action,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
}

type InviteParam struct {
	Email   string `json:"email"`
	Role    string `json:"role"`
	BaseURL string `json:"base_url"`
}

func (v InviteParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.BaseURL, validation.Required),
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleAdmin, model.RoleUser)),
	)
}

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
