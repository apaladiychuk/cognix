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
	docRepo       repository.DocumentRepository
	streamClient  messaging.Client
	embeddingCh   chan string
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
				err = e.runConnector(ctx, msg)
				if err != nil {
					zap.S().Errorf("Failed to run connector: %v", err)
				}
			case <-ctx.Done():
				close(e.embeddingCh)
				return
			}
		}
	}()

	return nil
}

func (e *executor) runEmbedding() {
	for text := range e.embeddingCh {
		zap.S().Infof("sending embedded message: %s ...", text[:20])
	}
}

func (e *executor) runConnector(ctx context.Context, msg *messaging.Message) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	var trigger connector.Trigger
	if err := json.Unmarshal(msg.Body, &trigger); err != nil {
		return err
	}
	connectorModel, err := e.connectorRepo.GetByID(ctx, trigger.ID)
	if err != nil {
		return err
	}
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return err
	}

	connectorModel, err = connectorWF.Execute(ctx, trigger.Params)
	if err != nil {
		connectorModel.LastAttemptStatus = model.StatusFailed
	} else {
		connectorModel.LastAttemptStatus = model.StatusSuccess
	}
	if err = e.connectorRepo.UpdateStatistic(ctx, connectorModel); err != nil {
		return err
	}
	return nil
}

func NewExecutor(connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	streamClient messaging.Client) *executor {
	return &executor{
		connectorRepo: connectorRepo,
		docRepo:       docRepo,
		streamClient:  streamClient,
		embeddingCh:   make(chan string),
	}
}
