package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type (
	ConnectorRepository interface {
		GetActive(ctx context.Context) ([]*model.Connector, error)
		GetAllByUser(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Connector, error)
		GetByIDAndUser(ctx context.Context, tenantID, userID uuid.UUID, id int64) (*model.Connector, error)
		GetByID(ctx context.Context, id int64) (*model.Connector, error)
		GetBySource(ctx context.Context, tenantID, userID uuid.UUID, source model.SourceType) (*model.Connector, error)
		Create(ctx context.Context, connector *model.Connector) error
		Update(ctx context.Context, connector *model.Connector) error
		InvalidateConnector(ctx context.Context, connector *model.Connector) error
	}
	connectorRepository struct {
		db *pg.DB
	}
)

func (r *connectorRepository) InvalidateConnector(ctx context.Context, connector *model.Connector) error {
	if _, err := r.db.WithContext(ctx).Model(&model.Document{}).
		Set("status = ?", model.StatusPending).
		Where("connector_id  = ?", connector.ID.IntPart()).Update(); err != nil {
		return utils.Internal.Wrap(err, "Error while invalidating connector")
	}
	return nil
}

func (r *connectorRepository) GetBySource(ctx context.Context, tenantID, userID uuid.UUID, source model.SourceType) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(ctx).Model(&connector).
		Where("source = ?", source).
		Where("tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("shared = ?", true), nil
		}).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "ca not find connector")
	}
	return &connector, nil
}

func NewConnectorRepository(db *pg.DB) ConnectorRepository {
	return &connectorRepository{db: db}
}

func (r *connectorRepository) GetAllByUser(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Connector, error) {
	connectors := make([]*model.Connector, 0)
	if err := r.db.WithContext(ctx).Model(&connectors).
		Where("tenant_id = ?", tenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("shared = ?", true), nil
		}).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connectors")
	}
	return connectors, nil
}

func (r *connectorRepository) GetByIDAndUser(ctx context.Context, tenantID, userID uuid.UUID, id int64) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(ctx).Model(&connector).
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

func (r *connectorRepository) GetByID(ctx context.Context, id int64) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(ctx).Model(&connector).
		Relation("Docs").
		Relation("Credential").
		Where("connector.id = ?", id).
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connector")
	}
	connector.DocsMap = make(map[string]*model.Document)
	for _, doc := range connector.Docs {
		connector.DocsMap[doc.DocumentID] = doc
	}
	return &connector, nil
}

func (r *connectorRepository) Create(ctx context.Context, connector *model.Connector) error {
	stm := r.db.WithContext(ctx).Model(connector)
	if !connector.CredentialID.Valid {
		stm = stm.ExcludeColumn("credential_id")
	}
	if _, err := stm.Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create connector")
	}
	return nil
}

func (r *connectorRepository) Update(ctx context.Context, connector *model.Connector) error {
	if _, err := r.db.WithContext(ctx).Model(connector).Where("id = ?", connector.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update connector")
	}
	return nil
}

func (r *connectorRepository) GetActive(ctx context.Context) ([]*model.Connector, error) {
	connectors := make([]*model.Connector, 0)
	if err := r.db.WithContext(ctx).
		Model(&connectors).
		Where("disabled = false").Select(); err != nil {
		return nil, utils.Internal.Wrap(err, "can not load connectors")
	}
	return connectors, nil

}
