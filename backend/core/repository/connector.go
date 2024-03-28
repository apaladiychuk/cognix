package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type (
	ConnectorRepository interface {
		GetAll(c context.Context, tenantID, userID string) ([]*model.Connector, error)
		GetByID(c context.Context, id int, tenantID, userID string) (*model.Connector, error)
		Create(c context.Context, connector *model.Connector) error
		Update(c context.Context, connector *model.Connector) error
	}
	connectorRepository struct {
		db *pg.DB
	}
)

func NewConnectorRepository(db *pg.DB) ConnectorRepository {
	return &connectorRepository{db: db}
}

func (r *connectorRepository) GetAll(c context.Context, tenantID, userID string) ([]*model.Connector, error) {
	connectors := make([]*model.Connector, 0)
	if err := r.db.WithContext(c).Model(&connectors).
		Where("tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("shared = ?", true), nil
		}).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connectors")
	}
	return connectors, nil
}

func (r *connectorRepository) GetByID(c context.Context, id int, tenantID, userID string) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(c).Model(&connector).
		Where("tenant_id = ?", tenantID).
		Where("id = ?", id).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("shared = ?", true), nil
		}).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connector")
	}
	return &connector, nil
}

func (r *connectorRepository) Create(c context.Context, connector *model.Connector) error {
	//TODO implement me
	panic("implement me")
}

func (r *connectorRepository) Update(c context.Context, connector *model.Connector) error {
	//TODO implement me
	panic("implement me")
}
