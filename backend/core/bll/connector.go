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
	ConnectorBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error)
	}
	connectorBL struct {
		connectorRepo  repository.ConnectorRepository
		credentialRepo repository.CredentialRepository
	}
)

func NewConnectorBL(connectorRepo repository.ConnectorRepository, credentialRepo repository.CredentialRepository) ConnectorBL {
	return &connectorBL{
		connectorRepo:  connectorRepo,
		credentialRepo: credentialRepo,
	}
}

func (c *connectorBL) Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error) {
	cred, err := c.credentialRepo.GetByID(ctx, param.CredentialID, user.TenantID, user.ID)
	if err != nil {
		return nil, err
	}
	if cred.Source != model.SourceType(param.Source) {
		return nil, utils.InvalidInput.New("wrong credential source")
	}
	connector := model.Connector{
		CredentialID:            param.CredentialID,
		Name:                    param.Name,
		Source:                  model.SourceType(param.Source),
		InputType:               param.InputType,
		ConnectorSpecificConfig: param.ConnectorSpecificConfig,
		RefreshFreq:             param.RefreshFreq,
		UserID:                  user.ID,
		TenantID:                user.TenantID,
		Shared:                  param.Shared,
		Disabled:                param.Disabled,
		CreatedDate:             time.Now().UTC(),
	}
	if err = c.connectorRepo.Create(ctx, &connector); err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *connectorBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error) {
	connector, err := c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	cred, err := c.credentialRepo.GetByID(ctx, param.CredentialID, user.TenantID, user.ID)
	if err != nil {
		return nil, err
	}
	if cred.Source != connector.Source {
		return nil, utils.InvalidInput.New("wrong credential source")
	}
	connector.ConnectorSpecificConfig = param.ConnectorSpecificConfig
	connector.CredentialID = param.CredentialID
	connector.Name = param.Name
	connector.InputType = param.InputType
	connector.RefreshFreq = param.RefreshFreq
	connector.Shared = param.Shared
	connector.Disabled = param.Disabled
	connector.UpdatedDate = pg.NullTime{time.Now().UTC()}

	//	sql.NullTime{
	//Time:  time.Now().UTC(),
	//Valid: true,
	//	} //  null.TimeFrom(time.Now().UTC())
	if err = c.connectorRepo.Update(ctx, connector); err != nil {
		return nil, err
	}
	return connector, nil
}

func (c *connectorBL) GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error) {
	return c.connectorRepo.GetAll(ctx, user.TenantID, user.ID)
}

func (c *connectorBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error) {
	return c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
}
