package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

type executor struct {
	connectorRepo repository.ConnectorRepository
	streamClient  messaging.Client
}

func (e *executor) run(ctx context.Context) error {
	ch, err := e.streamClient.Listen(model.TopicExecutor)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case msg := <-ch:
				_, err = e.runConnector(ctx, msg)
				if err != nil {
					zap.S().Errorf("Failed to run connector: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (e *executor) runConnector(ctx context.Context, msg *messaging.Message) (context.Context, error) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	var trigger connector.Trigger
	if err := json.Unmarshal(msg.Body, &trigger); err != nil {
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
	}
}
