package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
	"strings"
)

type Base struct {
	collectionName string
	model          *model.Connector
	resultCh       chan *proto.TriggerResponse
}

type Connector interface {
	CollectionName() string
	Execute(ctx context.Context, param map[string]string) chan *proto.TriggerResponse
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

type nopConnector struct {
	Base
}

func (n *nopConnector) Execute(ctx context.Context, param map[string]string) chan *proto.TriggerResponse {
	ch := make(chan *proto.TriggerResponse)
	return ch
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
	b.resultCh = make(chan *proto.TriggerResponse, 10)
	if connector.Shared {
		b.collectionName = strings.ReplaceAll(fmt.Sprintf(model.CollectionTenant, connector.TenantID), "-", "")
	} else {
		b.collectionName = strings.ReplaceAll(fmt.Sprintf(model.CollectionUser, connector.UserID), "-", "")
	}

	return
}

func (b *Base) CollectionName() string {
	return b.collectionName
}
