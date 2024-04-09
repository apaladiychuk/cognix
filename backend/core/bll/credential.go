package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"time"
)

type (
	CredentialBL interface {
		GetAll(ctx context.Context, user *model.User, param *parameters.GetAllCredentialsParam) ([]*model.Credential, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Credential, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateCredentialParam) (*model.Credential, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateCredentialParam) (*model.Credential, error)
	}
	credentialBL struct {
		credentialRepo repository.CredentialRepository
	}
)

func NewCredentialBL(credentialRepo repository.CredentialRepository) CredentialBL {
	return &credentialBL{
		credentialRepo: credentialRepo,
	}
}

func (c *credentialBL) GetAll(ctx context.Context, user *model.User, param *parameters.GetAllCredentialsParam) ([]*model.Credential, error) {
	return c.credentialRepo.GetAll(ctx, user.TenantID, user.ID, param)
}

func (c *credentialBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Credential, error) {
	return c.credentialRepo.GetByID(ctx, id, user.TenantID, user.ID)
}

func (c *credentialBL) Create(ctx context.Context, user *model.User, param *parameters.CreateCredentialParam) (*model.Credential, error) {
	credential := model.Credential{
		UserID:         user.ID,
		TenantID:       user.TenantID,
		Source:         param.Source,
		CreatedDate:    time.Now().UTC(),
		Shared:         param.Shared,
		CredentialJson: param.CredentialJson,
	}
	if err := c.credentialRepo.Create(ctx, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

func (c *credentialBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateCredentialParam) (*model.Credential, error) {
	credential, err := c.credentialRepo.GetByID(ctx, id, user.TenantID, user.ID)
	if err != nil {
		return nil, err
	}
	if credential.UserID != user.ID {
		return nil, utils.ErrorPermission.New("you are not credential owner.")
	}
	credential.CredentialJson = param.CredentialJson
	credential.Shared = param.Shared
	credential.UpdatedDate = pg.NullTime{time.Now().UTC()}
	if err = c.credentialRepo.Update(ctx, credential); err != nil {
		return nil, err
	}
	return credential, nil
}
