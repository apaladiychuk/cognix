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
	"time"
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
	if err = e.connectorRepo.InvalidateConnector(ctx, connectorModel); err != nil {
		return err
	}
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return err
	}

	resultCh := connectorWF.Execute(ctx, trigger.Params)

	for result := range resultCh {
		doc, errr := e.handleResult(ctx, connectorModel, result)
		if errr != nil {
			doc.Status = model.StatusFailed
			doc.IsUpdated = true
			err = errr
		}
		if !doc.IsUpdated {
			continue
		}
		if doc.ID.IntPart() != 0 {
			errr = e.docRepo.Update(ctx, doc)
		} else {
			errr = e.docRepo.Create(ctx, doc)
		}
		if err != nil {
			err = errr
			zap.S().Errorf("Failed to update document: %v", err)
		}
	}

	if err != nil {
		connectorModel.LastAttemptStatus = model.StatusFailed
	} else {
		connectorModel.LastAttemptStatus = model.StatusSuccess
	}
	if err = e.connectorRepo.Update(ctx, connectorModel); err != nil {
		return err
	}
	return nil
}

func (e *executor) handleResult(ctx context.Context, connectorModel *model.Connector, result *proto.TriggerResponse) (*model.Document, error) {
	doc, ok := connectorModel.DocsMap[result.GetDocumentId()]
	if !ok {
		doc = &model.Document{
			DocumentID:  result.GetDocumentId(),
			ConnectorID: connectorModel.ID,
			Link:        result.GetUrl(),
			Signature:   result.GetSignature(),
			CreatedDate: time.Now().UTC(),
			Status:      model.StatusInProgress,
			IsUpdated:   true,
		}
	} else {
		if doc.Signature != result.GetSignature() {
			doc.Signature = result.GetSignature()
			doc.Status = model.StatusInProgress
			doc.IsUpdated = true
		}
	}
	if !doc.IsUpdated {
		doc.Status = model.StatusSuccess
		return doc, nil
	}

	return doc, e.msgClient.Publish(ctx, model.TopicEmbedding,
		&proto.Body{Payload: &proto.Body_Embedding{Embedding: &proto.EmbeddingRequest{
			Id:         connectorModel.ID.IntPart(),
			DocumentId: doc.ID.IntPart(),
			Key:        doc.DocumentID,
			Content:    result.Content,
		}}},
	)
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
