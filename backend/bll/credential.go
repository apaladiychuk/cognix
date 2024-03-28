package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"gopkg.in/guregu/null.v4"
	"time"
)

type (
	CredentialBL interface {
		GetAll(ctx context.Context, user *model.User, source string) ([]*model.Credential, error)
		GetByID(ctx context.Context, user *model.User, id int) (*model.Credential, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateCredentialParam) (*model.Credential, error)
		Update(ctx context.Context, id int, user *model.User, param *parameters.UpdateCredentialParam) (*model.Credential, error)
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

func (c *credentialBL) GetAll(ctx context.Context, user *model.User, source string) ([]*model.Credential, error) {
	return c.credentialRepo.GetAll(ctx, user.TenantID.String(), user.ID.String(), source)
}

func (c *credentialBL) GetByID(ctx context.Context, user *model.User, id int) (*model.Credential, error) {
	return c.credentialRepo.GetByID(ctx, id, user.TenantID.String(), user.ID.String())
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

func (c *credentialBL) Update(ctx context.Context, id int, user *model.User, param *parameters.UpdateCredentialParam) (*model.Credential, error) {
	credential, err := c.credentialRepo.GetByID(ctx, id, user.TenantID.String(), user.ID.String())
	if err != nil {
		return nil, err
	}
	if credential.UserID != user.ID {
		return nil, utils.ErrorPermission.New("you are not credential owner.")
	}
	credential.CredentialJson = param.CredentialJson
	credential.Shared = param.Shared
	credential.UpdatedDate = null.TimeFrom(time.Now().UTC())
	if err = c.credentialRepo.Update(ctx, credential); err != nil {
		return nil, err
	}
	return credential, nil
}
