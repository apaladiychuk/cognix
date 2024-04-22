package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type executor struct {
	connectorRepo repository.ConnectorRepository
	streamClient  messaging.Client
}

func (e *executor) run(ctx context.Context) error {
	ch, err := e.streamClient.Listen(connector.TopicExecutor)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-ch:
			_, err = e.runConnector(ctx, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (e *executor) runConnector(ctx context.Context, msg *messaging.Message) (context.Context, error) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	var trigger connector.Trigger
	if err := msg.Body.ToStruct(&trigger); err != nil {
		return ctx, err
	}
	connectorModel, err := e.connectorRepo.GetByID(ctx, trigger.ID)
	if err != nil {
		return ctx, err
	}
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return ctx, err
	}
	return ctx, connectorWF.Execute(ctx, trigger.Params)
}

func NewExecutor(connectorRepo repository.ConnectorRepository,
	streamClient messaging.Client) *executor {
	return &executor{
		connectorRepo: connectorRepo,
		streamClient:  streamClient,
		tracer:        otel.Tracer("connector"),
	}
}
