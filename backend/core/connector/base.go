package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
)

type Base struct {
	collectionName string
	model          *model.Connector
	embeddingCh    chan string
}

type Connector interface {
	Config(connector *model.Connector) (Connector, error)
	Execute(ctx context.Context, param model.JSONMap) (*model.Connector, error)
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

type Trigger struct {
	ID     int64         `json:"id"`
	Params model.JSONMap `json:"params"`
}
type nopConnector struct {
	Base
}

func (n *nopConnector) Config(connector *model.Connector) (Connector, error) {
	return n, nil
}

func (n *nopConnector) Execute(ctx context.Context, param model.JSONMap) (*model.Connector, error) {
	return &model.Connector{}, nil
}

func New(connectorModel *model.Connector) (Connector, error) {
	switch connectorModel.Source {
	case model.SourceTypeWEB:
		return NewWeb(connectorModel)
	default:
		return &nopConnector{}, nil
	}
}

func (b *Base) Config(connector *model.Connector) {
	b.model = connector
	if connector.Shared {
		b.collectionName = fmt.Sprintf(model.CollectionTenant, connector.TenantID)
	} else {
		b.collectionName = fmt.Sprintf(model.CollectionUser, connector.UserID)
	}
	return
}
