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
	ExternalID string      `json:"external_id"`
	Roles      StringSlice `json:"roles" pg:",array"`
}
