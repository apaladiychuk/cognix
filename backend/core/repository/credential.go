package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type (
	CredentialRepository interface {
		GetAll(c context.Context, tenantID, userID, source string) ([]*model.Credential, error)
		GetByID(c context.Context, id int64, tenantID, userID string) (*model.Credential, error)
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

func (r *credentialRepository) GetAll(c context.Context, tenantID, userID, source string) ([]*model.Credential, error) {
	credentials := make([]*model.Credential, 0)
	stm := r.db.WithContext(c).Model(&credentials).
		Where("tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("shared = ?", true), nil
		})
	if source != "" {
		stm = stm.Where("source = ?", source)
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

func (r *credentialRepository) GetByID(c context.Context, id int64, tenantID, userID string) (*model.Credential, error) {
	var credential model.Credential
	if err := r.db.WithContext(c).Model(&credential).
		Where("id = ?", id).
		Where("tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).WhereOr("shared = ?", true), nil
		}).
		First(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find credential [%s]", err.Error())
	}
	return &credential, nil
}
