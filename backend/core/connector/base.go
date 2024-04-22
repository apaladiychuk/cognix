package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
)

const (
	TopicExecutor   = "executor"
	ConnectorTracer = "connector"
)

type Base struct {
	collectionName string
	model          *model.Connector
}

type Connector interface {
	Config(connector *model.Connector) (Connector, error)
	Execute(ctx context.Context, param model.JSONMap) error
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

type Trigger struct {
	ID     int64         `json:"id"`
	Params model.JSONMap `json:"params"`
}
type nopConnector struct{}

func (n *nopConnector) Config(connector *model.Connector) (Connector, error) {
	return n, nil
}

func (n *nopConnector) Execute(ctx context.Context, param model.JSONMap) error {
	return nil
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
