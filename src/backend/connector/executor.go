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
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"time"
)

type Executor struct {
	cfg            *Config
	connectorRepo  repository.ConnectorRepository
	credentialRepo repository.CredentialRepository
	docRepo        repository.DocumentRepository
	msgClient      messaging.Client
	chunking       ai.Chunking
	minioClient    storage.MinIOClient
	milvusClient   storage.MilvusClient
	oauthClient    *resty.Client
	downloadClient *resty.Client
}

// run this method listen messages from nats
func (e *Executor) run(streamName, topic string, task messaging.MessageHandler) {
	if err := e.msgClient.Listen(context.Background(), streamName, topic, task); err != nil {
		zap.S().Errorf("failed to listen[%s]: %v", topic, err)
	}
	return
}

// runConnector run connector from nats message
func (e *Executor) runConnector(ctx context.Context, msg jetstream.Msg) error {

	//ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header()))
	var trigger proto.ConnectorRequest

	if err := proto2.Unmarshal(msg.Data(), &trigger); err != nil {
		zap.S().Errorf("Error unmarshalling message: %s", err.Error())
		return err
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

	// create new instance of connector by connector model
	connectorWF, err := connector.New(connectorModel)
	if err != nil {
		return err
	}
	if trigger.Params == nil {
		trigger.Params = make(map[string]string)
	}
	// add file limit parameter to connector
	trigger.Params[connector.ParamFileLimit] = fmt.Sprintf("%d", e.cfg.FileLimit*connector.GB)
	// execute connector
	resultCh := connectorWF.Execute(ctx, trigger.Params)
	// read result from channel
	for result := range resultCh {
		var loopErr error
		// save content in minio
		if result.SaveContent {
			if err = e.saveContent(ctx, result); err != nil {
				loopErr = err
				zap.S().Errorf("failed to save content: %v", err)
			}
		}
		// find or create document from result
		doc := e.handleResult(connectorModel, result)
		// create or update document in database
		if doc.ID.IntPart() != 0 {
			loopErr = e.docRepo.Update(ctx, doc)
		} else {
			loopErr = e.docRepo.Create(ctx, doc)
		}

		if loopErr != nil {
			err = loopErr
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}

		// send message to chunking service
		if loopErr = e.msgClient.Publish(ctx, e.msgClient.StreamConfig().ChunkerStreamSubject,
			&proto.ChunkingData{
				Url:            result.URL,
				DocumentId:     doc.ID.IntPart(),
				FileType:       result.GetType(),
				CollectionName: connectorModel.CollectionName(),
			}); loopErr != nil {
			err = loopErr
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}
	}
	// remove documents that were removed from source
	var ids []int64
	for _, doc := range connectorModel.DocsMap {
		if doc.IsExists || doc.ID.IntPart() == 0 {
			continue
		}
		ids = append(ids, doc.ID.IntPart())
	}
	if len(ids) > 0 {
		if loopErr := e.docRepo.ArchiveRestore(ctx, false, ids...); loopErr != nil {
			err = loopErr
		}
	}

	if err != nil {
		connectorModel.LastAttemptStatus = model.ConnectorStatusError
	} else {
		connectorModel.LastAttemptStatus = model.ConnectorStatusSuccess
	}
	connectorModel.LastSuccessfulIndexDate = pg.NullTime{time.Now().UTC()}
	connectorModel.LastUpdate = pg.NullTime{time.Now().UTC()}
	if len(ids) > 0 {
		if err = e.milvusClient.Delete(ctx, connectorModel.CollectionName(), ids...); err != nil {
			//return err
		}
	}
	if err = e.connectorRepo.Update(ctx, connectorModel); err != nil {
		return err
	}
	return nil
}

func (e *Executor) saveContent(ctx context.Context, response *connector.Response) error {

	downloadResp, err := e.downloadClient.R().Get(response.URL)

	if err != nil || downloadResp.IsError() {
		return fmt.Errorf("download fail: %v]", err)
	}
	defer downloadResp.RawBody().Close()

	url, _, err := e.minioClient.Upload(ctx, response.Name, response.MimeType, downloadResp.RawBody())
	if err != nil {
		zap.S().Errorf("Failed to upload file: %v", err)
		return err
	}
	response.URL = url
	return nil
}

func (e *Executor) handleResult(connectorModel *model.Connector, result *connector.Response) *model.Document {
	doc, ok := connectorModel.DocsMap[result.SourceID]
	if !ok {
		doc = &model.Document{
			SourceID:     result.SourceID,
			ConnectorID:  connectorModel.ID,
			Link:         result.URL,
			CreationDate: time.Now().UTC(),
		}
		connectorModel.DocsMap[result.SourceID] = doc
	}
	return doc
}

// refreshToken  refresh OAuth token and store credential in database
func (e *Executor) refreshToken(ctx context.Context, cm *model.Connector) error {
	provider, ok := model.ConnectorAuthProvider[cm.Type]
	if !ok {
		return nil
	}
	token, ok := cm.ConnectorSpecificConfig["token"]
	if !ok {
		return fmt.Errorf("wrong token")
	}

	response, err := e.oauthClient.R().SetContext(ctx).
		SetBody(token).Post(fmt.Sprintf("/api/oauth/%s/refresh_token", provider))

	if err != nil || response.IsError() {
		return fmt.Errorf("failed to refresh token: %v : %v", err, response.Error())
	}
	var payload struct {
		Data oauth2.Token `json:"data"`
	}

	if err = json.Unmarshal(response.Body(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshl token: %v : %v", err, response.Error())
	}
	cm.ConnectorSpecificConfig["token"] = payload.Data
	if err = e.connectorRepo.Update(ctx, cm); err != nil {
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
	milvusClient storage.MilvusClient,
) *Executor {
	return &Executor{
		cfg:            cfg,
		connectorRepo:  connectorRepo,
		credentialRepo: credentialRepo,
		docRepo:        docRepo,
		msgClient:      streamClient,
		chunking:       chunking,
		minioClient:    minioClient,
		milvusClient:   milvusClient,
		oauthClient: resty.New().
			SetTimeout(time.Minute).
			SetBaseURL(cfg.OAuthURL),
		downloadClient: resty.New().
			SetTimeout(time.Minute).
			SetDoNotParseResponse(true),
	}
}
