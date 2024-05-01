package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type (
	CredentialRepository interface {
		GetAll(c context.Context, tenantID, userID uuid.UUID, param *parameters.GetAllCredentialsParam) ([]*model.Credential, error)
		GetByID(c context.Context, id int64, tenantID, userID uuid.UUID, relations ...string) (*model.Credential, error)
		Create(c context.Context, cred *model.Credential) error
		Update(c context.Context, cred *model.Credential) error
	}
	credentialRepository struct {
		db *pg.DB
	}
)

func NewCredentialRepository(db *pg.DB) CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) GetAll(c context.Context, tenantID, userID uuid.UUID, param *parameters.GetAllCredentialsParam) ([]*model.Credential, error) {
	credentials := make([]*model.Credential, 0)
	stm := r.db.WithContext(c).Model(&credentials).
		Where("tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("shared = ?", true), nil
		})
	if param.Source != "" {
		stm = stm.Where("source = ?", param.Source)
	}
	if !param.Archived {
		stm = stm.Where("deleted_date is null")
	}

	if err := stm.Select(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find credentials [%s]", err.Error())
	}
	return credentials, nil
}

func (r *credentialRepository) Create(c context.Context, cred *model.Credential) error {
	if _, err := r.db.WithContext(c).Model(cred).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create credential")
	}
	return nil
}

func (r *credentialRepository) Update(c context.Context, cred *model.Credential) error {
	if _, err := r.db.WithContext(c).Model(cred).Where("id = ?", cred.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update credential")
	}
	return nil
}

func (r *credentialRepository) GetByID(c context.Context, id int64, tenantID, userID uuid.UUID, relations ...string) (*model.Credential, error) {
	var credential model.Credential
	stm := r.db.WithContext(c).Model(&credential).
		Where("credential.id = ?", id).
		Where("credential.tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("credential.user_id = ?", userID).WhereOr("credential.shared = ?", true), nil
		})
	for _, relation := range relations {
		stm = stm.Relation(relation)
	}
	if err := stm.First(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find credential [%s]", err.Error())
	}
	return &credential, nil
}
