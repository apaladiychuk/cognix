package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const TopicExecutor = "executor"

type executor struct {
	connectorRepo repository.ConnectorRepository
	streamClient  messaging.Client
}

func (e *executor) New(ctx context.Context, connectorModel *model.Connector) Connector {
	switch connectorModel.Source {
	case model.SourceTypeWEB:
		return connector.NewWeb(connectorModel)
	}
}

func (e *executor) run(ctx context.Context) error {
	ch, err := e.streamClient.Listen(TopicExecutor)
	if err != nil {
		return err
	}
	select {
	case msg := <-ch:

	case <-ctx.Done():
		return ctx.Err()

	}

}

func runConnector(msg *messaging.Message) (context.Context, error) {
	ctx := context.Background()
	otel.GetTextMapPropagator().Extract(ctx, &propagation.MapCarrier{})
}
