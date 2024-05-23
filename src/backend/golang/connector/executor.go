package main

import (
	"bytes"
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"time"
)

type Executor struct {
	connectorRepo    repository.ConnectorRepository
	credentialRepo   repository.CredentialRepository
	docRepo          repository.DocumentRepository
	msgClient        messaging.Client
	chunking         ai.Chunking
	minioClient      storage.MinIOClient
	oauthClient      *resty.Client
	cancel           context.CancelFunc
	subscriptionName string
}

func (e *Executor) run(_ context.Context, topic, subscriptionName string, task messaging.MessageHandler) {
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	if err := e.msgClient.Listen(ctx, topic, e.subscriptionName, task); err != nil {
		zap.S().Errorf("failed to listen[%s]: %v", topic, err)
	}
	return
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
	zap.S().Infof("run connector : %s [%d]", connectorModel.Name, connectorModel.ID.IntPart())
	// refresh token if needed
	if err = e.refreshToken(ctx, connectorModel); err != nil {
		return err
	}

	// todo move to connector
	if err = e.connectorRepo.InvalidateConnector(ctx, connectorModel); err != nil {
		return err
	}
	// create new instance of connector by connector model
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return err
	}
	// execute connector
	resultCh := connectorWF.Execute(ctx, trigger.Params)
	// read result from channel
	zap.S().Debug(" wait for result ...")
	for result := range resultCh {
		zap.S().Debugf(" receive %s ", result.URL)
		var loopErr error
		// save content in minio
		if result.SaveContent {
			if err = e.saveContent(ctx, result); err != nil {
				loopErr = err
				zap.S().Errorf("failed to save content: %v", err)
			}
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
				Url:            result.URL,
				DocumentId:     doc.ID.IntPart(),
				FileType:       result.GetType(),
				CollectionName: connectorModel.CollectionName(),
			}},
		}); loopErr != nil {
			err = loopErr
			doc.Status = model.StatusFailed
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}
	}
	zap.S().Debugf(" channel closed ")
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

func (e *Executor) saveContent(ctx context.Context, response *connector.Response) error {
	url, _, err := e.minioClient.Upload(ctx, response.Name, response.MimeType, bytes.NewBuffer(response.Content))
	if err != nil {
		zap.S().Errorf("Failed to upload file: %v", err)
		return err
	}
	response.URL = url
	return nil
}

func (e *Executor) handleResult(connectorModel *model.Connector, result *connector.Response) *model.Document {
	doc, ok := connectorModel.DocsMap[result.URL]
	if !ok {
		doc = &model.Document{
			DocumentID:  result.SourceID,
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

// refreshToken  refresh OAuth token and store credential in database
func (e *Executor) refreshToken(ctx context.Context, cm *model.Connector) error {
	if cm.Credential == nil || cm.Credential.CredentialJson.Provider == model.ProviderCustom {
		return nil
	}
	response, err := e.oauthClient.R().SetContext(ctx).
		SetBody(cm.Credential.CredentialJson.Token).Post(fmt.Sprintf("/api/oauth/%s/refresh_token", cm.Credential.CredentialJson.Provider))

	if err != nil || response.IsError() {
		return fmt.Errorf("failed to refresh token: %v : %v", err, response.Error())
	}
	if err = json.Unmarshal(response.Body(), cm.Credential.CredentialJson.Token); err != nil {
		return fmt.Errorf("failed to unmarshl token: %v : %v", err, response.Error())
	}
	if err = e.credentialRepo.Update(ctx, cm.Credential); err != nil {
		return err
	}
	return nil
}

func NewExecutor(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	credentialRepo repository.CredentialRepository,
	docRepo repository.DocumentRepository,
	streamClient messaging.Client,
	chunking ai.Chunking,
	minioClient storage.MinIOClient,
	oauthClient *resty.Client,
) *Executor {
	return &Executor{
		subscriptionName: cfg.SubscriptionName,
		connectorRepo:    connectorRepo,
		credentialRepo:   credentialRepo,
		docRepo:          docRepo,
		msgClient:        streamClient,
		chunking:         chunking,
		minioClient:      minioClient,
		oauthClient:      oauthClient,
	}
}
