package main

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"time"
)

type Executor struct {
	connectorRepo repository.ConnectorRepository
	docRepo       repository.DocumentRepository
	msgClient     messaging.Client
	chunking      ai.Chunking
	milvusClinet  storage.MilvusClient
}

func (e *Executor) run(ctx context.Context, topic, subscriptionName string, task messaging.MessageHandler) {
	if err := e.msgClient.Listen(ctx, topic, subscriptionName, task); err != nil {
		zap.S().Errorf("failed to listen[%s]: %v", topic, err)
	}
	return
}

func (e *Executor) runEmbedding(ctx context.Context, msg *proto.Message) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	payload := msg.GetBody().GetChunking()
	if payload == nil {
		zap.S().Errorf("Failed to get embedding payload")
		return nil
	}

	return nil
}

func (e *Executor) runConnector(ctx context.Context, msg *proto.Message) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	trigger := msg.GetBody().GetTrigger()

	if trigger == nil {
		return fmt.Errorf("failed to get trigger payload")
	}
	connectorModel, err := e.connectorRepo.GetByID(ctx, trigger.GetId())
	if err != nil {
		return err
	}
	// todo move to connector
	if err = e.connectorRepo.InvalidateConnector(ctx, connectorModel); err != nil {
		return err
	}
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return err
	}
	resultCh := connectorWF.Execute(ctx, trigger.Params)

	if err = e.milvusClinet.CreateSchema(ctx, connectorModel.CollectionName()); err != nil {
		return fmt.Errorf("error creating schema: %v", err)
	}

	for result := range resultCh {
		var loopErr error
		if result.SaveContent {
			e.saveContent(ctx, result)
		}
		doc := e.handleResult(connectorModel, result)

		if doc.ID.IntPart() != 0 {
			loopErr = e.docRepo.Update(ctx, doc)
		} else {
			loopErr = e.docRepo.Create(ctx, doc)
		}

		if loopErr != nil {
			err = loopErr
			doc.Status = model.StatusFailed
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}

		if loopErr = e.msgClient.Publish(ctx, model.TopicChunking, &proto.Body{
			Payload: &proto.Body_Chunking{Chunking: &proto.ChunkingData{
				Url:        result.URL,
				DocumentId: doc.ID.IntPart(),
				FileType:   result.GetType(),
			}},
		}); loopErr != nil {
			err = loopErr
			doc.Status = model.StatusFailed
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
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

func (e *Executor) saveContent(ctx context.Context, response *connector.Response) {

}

func (e *Executor) handleResult(connectorModel *model.Connector, result *connector.Response) *model.Document {
	doc, ok := connectorModel.DocsMap[result.URL]
	if !ok {
		doc = &model.Document{
			DocumentID:  result.URL,
			ConnectorID: connectorModel.ID,
			Link:        result.URL,
			CreatedDate: time.Now().UTC(),
			Status:      model.StatusInProgress,
		}
		connectorModel.DocsMap[result.URL] = doc
	}

	doc.Status = model.StatusInProgress

	return doc
}

func NewExecutor(connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	streamClient messaging.Client,
	chunking ai.Chunking,
	milvusClinet storage.MilvusClient,
) *Executor {
	return &Executor{
		connectorRepo: connectorRepo,
		docRepo:       docRepo,
		msgClient:     streamClient,
		chunking:      chunking,
		milvusClinet:  milvusClinet,
	}
}
