package parameters

import (
	"cognix.ch/api/v2/core/model"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/shopspring/decimal"
)

type LoginParam struct {
	RedirectURL string `form:"redirect_url"`
}

type OAuthParam struct {
	Action   string `json:"action,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
}

type InviteParam struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (v InviteParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleAdmin, model.RoleUser)),
	)
}

type ArchivedParam struct {
	Archived bool `form:"archived"`
}
type GetAllCredentialsParam struct {
	ArchivedParam
	Source string `form:"source"`
}

type CreateCredentialParam struct {
	Source         string                `json:"source"`
	Shared         bool                  `json:"shared"`
	CredentialJson *model.CredentialJson `json:"credential_json"`
}

func (v CreateCredentialParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Source, validation.Required,
			validation.By(func(value interface{}) error {
				if st, ok := model.AllSourceTypes[model.SourceType(v.Source)]; !ok || !st.IsImplemented {
					return fmt.Errorf("invalid source type")
				}
				return nil
			})))
}

type UpdateCredentialParam struct {
	Shared         bool                  `json:"shared"`
	CredentialJson *model.CredentialJson `json:"credential_json"`
}

type CreateConnectorParam struct {
	CredentialID            decimal.NullDecimal `json:"credential_id,omitempty"`
	Name                    string              `json:"name,omitempty"`
	Source                  string              `json:"source,omitempty"`
	ConnectorSpecificConfig model.JSONMap       `json:"connector_specific_config,omitempty"`
	RefreshFreq             int                 `json:"refresh_freq,omitempty"`
	Shared                  bool                `json:"shared,omitempty"`
	Disabled                bool                `json:"disabled,omitempty"`
}

func (v CreateConnectorParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Source, validation.Required,
			validation.By(func(value interface{}) error {
				if st, ok := model.AllSourceTypes[model.SourceType(v.Source)]; !ok || !st.IsImplemented {
					return fmt.Errorf("invalid source type")
				}
				return nil
			})))
}

type UpdateConnectorParam struct {
	CredentialID            decimal.NullDecimal `json:"credential_id,omitempty"`
	Name                    string              `json:"name,omitempty"`
	ConnectorSpecificConfig model.JSONMap       `json:"connector_specific_config,omitempty"`
	RefreshFreq             int                 `json:"refresh_freq,omitempty"`
	Shared                  bool                `json:"shared,omitempty"`
	Disabled                bool                `json:"disabled,omitempty"`
}

type AddUserParam struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (v AddUserParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleUser, model.RoleAdmin)),
	)
}

type EditUserParam struct {
	Role string `json:"role"`
}

func (v EditUserParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleUser, model.RoleAdmin)),
	)
}
