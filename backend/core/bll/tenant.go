package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
)

type (
	TenantBL interface {
		GetUsers(ctx context.Context, user *model.User) ([]*model.User, error)
	}
	tenantBL struct {
		tenantRepo repository.TenantRepository
	}
)

func (b *tenantBL) GetUsers(ctx context.Context, user *model.User) ([]*model.User, error) {
	if len(user.Roles) == 0 || user.Roles[0] == model.RoleUser {
		return nil, utils.ErrorPermission.New("access denied")
	}
	return b.tenantRepo.GetUsers(ctx, user.TenantID)
}

func NewTenantBL(tenantRepo repository.TenantRepository) TenantBL {
	return &tenantBL{tenantRepo: tenantRepo}
}
