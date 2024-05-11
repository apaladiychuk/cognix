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

func (e *Executor) run(ctx context.Context, topic, subscriptionName string, task messaging.MessageHandler) error {
	return e.msgClient.Listen(ctx, topic, subscriptionName, task)
}

func (e *Executor) runEmbedding(ctx context.Context, msg *proto.Message) error {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header))
	payload := msg.GetBody().GetEmbedding()
	if payload == nil {
		zap.S().Errorf("Failed to get embedding payload")
		return nil
	}
	zap.S().Infof("process embedding %d == > %50s ", payload.GetDocumentId(), payload.Content)
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
	if err = e.connectorRepo.InvalidateConnector(ctx, connectorModel); err != nil {
		return err
	}
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return err
	}
	embedding := ai.NewEmbeddingParser(&model.EmbeddingModel{ModelID: "text-embedding-ada-002"})
	resultCh := connectorWF.Execute(ctx, trigger.Params)

	if err = e.milvusClinet.CreateSchema(ctx, connectorWF.CollectionName()); err != nil {
		return fmt.Errorf("error creating schema: %v", err)
	}

	for result := range resultCh {
		var loopErr error
		doc := e.handleResult(ctx, connectorModel, result)
		if !doc.IsUpdated {
			continue
		}
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
		chunks, loopErr := e.chunking.Split(ctx, result.Content)
		if loopErr != nil {
			err = loopErr
			doc.Status = model.StatusFailed
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}
		embeddingResponse, loopErr := embedding.Parse(ctx, &proto.EmbeddingRequest{
			DocumentId: doc.ID.IntPart(),
			Key:        doc.DocumentID,
			Content:    chunks,
		})
		if loopErr != nil {
			err = loopErr
			doc.Status = model.StatusFailed
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}
		milvusPayload := make([]*storage.MilvusPayload, 0, len(embeddingResponse.Payload))
		for _, payload := range embeddingResponse.Payload {
			milvusPayload = append(milvusPayload, &storage.MilvusPayload{
				DocumentID: embeddingResponse.GetDocumentId(),
				Chunk:      payload.GetChunk(),
				Content:    payload.GetContent(),
				Vector:     payload.GetVector(),
			})
		}
		if loopErr = e.milvusClinet.Save(ctx, connectorWF.CollectionName(), milvusPayload...); loopErr != nil {
			err = loopErr
			doc.Status = model.StatusFailed
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}

		//todo for production
		//if loopErr = e.msgClient.Publish(ctx, model.TopicEmbedding,
		//	&proto.Body{MilvusPayload: &proto.Body_Embedding{Embedding: &proto.EmbeddingRequest{
		//		DocumentId: doc.ID.IntPart(),
		//		Key:        doc.DocumentID,
		//		Content:    chunks,
		//	}}},
		//); err != nil {
		//	err = loopErr
		//	zap.S().Errorf("Failed to update document: %v", err)
		//	doc.Status = model.StatusFailed
		//	doc.Signature = ""
		//	continue
		//}
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

func (e *Executor) handleResult(ctx context.Context, connectorModel *model.Connector, result *proto.TriggerResponse) *model.Document {
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
		connectorModel.DocsMap[result.GetDocumentId()] = doc
	} else {
		if doc.Signature != result.GetSignature() {
			doc.Signature = result.GetSignature()
			doc.Status = model.StatusInProgress
			doc.IsUpdated = true
		}
	}
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
