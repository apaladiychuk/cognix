package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

type taskRunner func(ctx context.Context, msg *proto.Message) error

type executor struct {
	connectorRepo repository.ConnectorRepository
	docRepo       repository.DocumentRepository
	msgClient     messaging.Client
	embeddingCh   chan string
}

func (e *executor) run(ctx context.Context, topic, subscriptionName string, task taskRunner) error {
	ch, err := e.msgClient.Listen(ctx, topic, subscriptionName)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case msg := <-ch:
				err = task(ctx, msg)
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

func (e *executor) runEmbedding(ctx context.Context, msg *proto.Message) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	payload := msg.GetBody().GetEmbedding()
	if payload == nil {
		zap.S().Errorf("Failed to get embedding payload")
	}
	zap.S().Infof("process embedding %d == > %50s ", payload.GetDocumentId(), payload.Content)
	return nil
}

func (e *executor) runConnector(ctx context.Context, msg *proto.Message) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	trigger := msg.GetBody().GetTrigger()

	if trigger == nil {
		return fmt.Errorf("failed to get trigger payload")
	}
	connectorModel, err := e.connectorRepo.GetByID(ctx, trigger.GetId())
	if err != nil {
		return err
	}
	connectorWF, err := connector.New(connectorModel, e.msgClient)
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
		msgClient:     streamClient,
		embeddingCh:   make(chan string),
	}
}
