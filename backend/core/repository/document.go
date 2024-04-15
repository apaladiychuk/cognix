package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type (
	DocumentRepository interface {
		FindByConnectorID(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error)
		Create(ctx context.Context, document ...*model.Document) error
	}
	documentRepository struct {
		db *pg.DB
	}
)

func (r *documentRepository) FindByConnectorID(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error) {
	documents := make([]*model.Document, 0)
	if err := r.db.WithContext(ctx).Model(&documents).
		Join("INNER JOIN connectors c ON c.id = connector_id").
		Where("connector_id = ?", connectorID).
		Where("c.tenant_id = ?", user.TenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("c.user_id = ? ", user.ID).
				WhereOr("c.shared = ?", true), nil
		}).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find documents ")
	}
	return documents, nil
}

func (r *documentRepository) Create(ctx context.Context, document ...*model.Document) error {
	if _, err := r.db.WithContext(ctx).Model(&document).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not insert document")
	}
	return nil
}

func NewDocumentRepository(db *pg.DB) DocumentRepository {
	return &documentRepository{db: db}
}
