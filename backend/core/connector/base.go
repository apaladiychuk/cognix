package connector

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
)

type Base struct {
	collectionName string
	model          *model.Connector
	msgClient      messaging.Client
	resultCh       chan <-
}

type Connector interface {
	Execute(ctx context.Context, param map[string]string) (*model.Connector, error)
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

type nopConnector struct {
	Base
}

func (n *nopConnector) Execute(ctx context.Context, param map[string]string) (*model.Connector, error) {
	return &model.Connector{}, nil
}

func New(connectorModel *model.Connector, msgClient messaging.Client) (Connector, error) {
	switch connectorModel.Source {
	case model.SourceTypeWEB:
		return NewWeb(connectorModel, msgClient)
	default:
		return &nopConnector{}, nil
	}
}

func (b *Base) Config(connector *model.Connector, msgClient messaging.Client) {
	b.model = connector
	b.msgClient = msgClient
	if connector.Shared {
		b.collectionName = fmt.Sprintf(model.CollectionTenant, connector.TenantID)
	} else {
		b.collectionName = fmt.Sprintf(model.CollectionUser, connector.UserID)
	}
	return
}

func (b *Base) sendResult(ctx context.Context, payload *proto.EmbeddingRequest) error {
	return b.msgClient.Publish(ctx, model.TopicEmbedding, &proto.Body{Payload: &proto.Body_Embedding{Embedding: payload}})
}
