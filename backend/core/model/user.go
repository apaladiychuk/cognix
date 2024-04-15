package model

import "github.com/google/uuid"

const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
)

type User struct {
	ID         uuid.UUID   `json:"id"`
	TenantID   uuid.UUID   `json:"tenant_id"`
	UserName   string      `json:"user_name"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	ExternalID string      `json:"-"`
	Roles      StringSlice `json:"roles" pg:",array"`
	Tenant     *Tenant     `json:"tenant,omitempty" pg:"rel:has-one"`
}

func (u *User) HasRoles(role ...string) bool {
	for _, r := range role {
		for _, ur := range u.Roles {
			if ur == r {
				return true
			}
		}
	}
	return false
}
