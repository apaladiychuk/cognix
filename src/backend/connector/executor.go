package main

import (
	"bytes"
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"strings"

	"io"
	"time"
)

type Executor struct {
	cfg            *Config
	connectorRepo  repository.ConnectorRepository
	docRepo        repository.DocumentRepository
	msgClient      messaging.Client
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
	startTime := time.Now()
	//ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header()))
	var trigger proto.ConnectorRequest

	if err := proto2.Unmarshal(msg.Data(), &trigger); err != nil {
		zap.S().Errorf("Error unmarshalling message: %s", err.Error())
		return err
	}
	// read connector model with documents, embedding model
	connectorModel, err := e.connectorRepo.GetByID(ctx, trigger.GetId())
	if err != nil {
		return err
	}
	defer func() {
		zap.S().Infof("connector %s completed. elapsed time: %d ms", connectorModel.Name, time.Since(startTime)/time.Millisecond)
	}()

	zap.S().Infof("receive message : %s [%d]", connectorModel.Name, connectorModel.ID.IntPart())
	// refresh token if needed
	connectorModel.Status = model.ConnectorStatusWorking

	// create new instance of connector by connector model
	connectorWF, err := connector.New(connectorModel, e.connectorRepo, e.cfg.OAuthURL)
	if err != nil {
		return err
	}
	if trigger.Params == nil {
		trigger.Params = make(map[string]string)
	}
	// execute connector
	resultCh := connectorWF.Execute(ctx, trigger.Params)
	// read result from channel
	hasSemanticMessage := false
	for result := range resultCh {
		var loopErr error
		// empty result when channel was closed.
		if result.SourceID == "" {
			break
		}
		hasSemanticMessage = true

		// save content in minio
		if result.Content != nil {
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

		if _, ok := model.WhisperFileTypes[result.FileType]; ok {
			// send message to whisper
			whisperDate := proto.WhisperData{
				Url:            result.URL,
				DocumentId:     doc.ID.IntPart(),
				ConnectorId:    connectorModel.ID.IntPart(),
				FileType:       result.FileType,
				CollectionName: connectorModel.CollectionName(),
				ModelName:      connectorModel.User.EmbeddingModel.ModelID,
				ModelDimension: int32(connectorModel.User.EmbeddingModel.ModelDim),
			}
			zap.S().Infof("send message to whisper %s - %s", connectorModel.Name, result.URL)
			if loopErr = e.msgClient.Publish(ctx,
				e.msgClient.StreamConfig().WhisperStreamName,
				e.msgClient.StreamConfig().WhisperStreamSubject,
				&whisperDate); loopErr != nil {
				err = loopErr
				zap.S().Errorf("Failed to publish whisper : %v", loopErr)
				continue
			}

		} else {
			// send message to semantic
			semanticData := proto.SemanticData{
				Url:            result.URL,
				DocumentId:     doc.ID.IntPart(),
				ConnectorId:    connectorModel.ID.IntPart(),
				FileType:       result.FileType,
				CollectionName: connectorModel.CollectionName(),
				ModelName:      connectorModel.User.EmbeddingModel.ModelID,
				ModelDimension: int32(connectorModel.User.EmbeddingModel.ModelDim),
			}
			zap.S().Infof("send message to semantic %s - %s", connectorModel.Name, result.URL)
			if loopErr = e.msgClient.Publish(ctx,
				e.msgClient.StreamConfig().SemanticStreamName,
				e.msgClient.StreamConfig().SemanticStreamSubject,
				&semanticData); loopErr != nil {
				err = loopErr
				zap.S().Errorf("Failed to publish semantic: %v", loopErr)
				continue
			}
		}
	}
	if errr := e.deleteUnusedFiles(ctx, connectorModel); err != nil {
		zap.S().Errorf("deleting unused files: %v", errr)
		if err == nil {
			err = errr
		}
	}
	if err != nil {
		zap.S().Errorf("failed to update documents: %v", err)
		connectorModel.Status = model.ConnectorStatusUnableProcess
	} else {
		if !hasSemanticMessage {
			connectorModel.Status = model.ConnectorStatusSuccess
		}
	}
	connectorModel.LastUpdate = pg.NullTime{time.Now().UTC()}

	if err = e.connectorRepo.Update(ctx, connectorModel); err != nil {
		return err
	}
	return nil
}

func (e *Executor) deleteUnusedFiles(ctx context.Context, connector *model.Connector) error {
	var ids []int64
	for _, doc := range connector.DocsMap {
		if doc.IsExists || doc.ID.IntPart() == 0 {
			continue
		}
		filepath := strings.Split(doc.URL, ":")
		if len(filepath) == 3 && filepath[0] == "minio" {
			if err := e.minioClient.DeleteObject(ctx, filepath[1], filepath[2]); err != nil {
				return err
			}
		}
		ids = append(ids, doc.ID.IntPart())
	}
	if len(ids) > 0 {
		if err := e.milvusClient.Delete(ctx, connector.CollectionName(), ids...); err != nil {
			return err
		}
		return e.docRepo.DeleteByIDS(ctx, ids...)
	}
	return nil
}
func (e *Executor) saveContent(ctx context.Context, response *connector.Response) error {

	var reader io.Reader
	//  download file if url presented.
	if response.Content.URL != "" {
		fileResponse, err := e.downloadClient.R().
			SetDoNotParseResponse(true).
			Get(response.Content.URL)
		defer fileResponse.RawBody().Close()
		if err = utils.WrapRestyError(fileResponse, err); err != nil {
			return err
		}
		reader = fileResponse.RawBody()
	} else {
		// create reader from raw content
		reader = bytes.NewReader(response.Content.Body)
	}

	fileName, _, err := e.minioClient.Upload(ctx, response.Content.Bucket, response.Name, response.MimeType, reader)
	if err != nil {
		return err
	}
	zap.S().Debugf("save fileName %s response name %s ", fileName, response.Name)
	response.URL = fmt.Sprintf("minio:%s:%s", response.Content.Bucket, fileName)
	return nil
}

func (e *Executor) handleResult(connectorModel *model.Connector, result *connector.Response) *model.Document {
	doc, ok := connectorModel.DocsMap[result.SourceID]
	if !ok {
		doc = &model.Document{
			SourceID:     result.SourceID,
			ConnectorID:  connectorModel.ID,
			URL:          result.URL,
			Signature:    result.Signature,
			CreationDate: time.Now().UTC(),
		}
		connectorModel.DocsMap[result.SourceID] = doc
	} else {
		doc.URL = result.URL
		doc.LastUpdate = pg.NullTime{time.Now().UTC()}
	}

	return doc
}

func NewExecutor(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	streamClient messaging.Client,
	minioClient storage.MinIOClient,
	milvusClient storage.MilvusClient,
) *Executor {
	return &Executor{
		cfg:           cfg,
		connectorRepo: connectorRepo,
		docRepo:       docRepo,
		msgClient:     streamClient,
		minioClient:   minioClient,
		milvusClient:  milvusClient,
		oauthClient: resty.New().
			SetTimeout(time.Minute).
			SetBaseURL(cfg.OAuthURL),
		downloadClient: resty.New().
			SetTimeout(time.Minute).
			SetDoNotParseResponse(true),
	}
}
