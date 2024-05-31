package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

const (
	ConnectorSchedulerSpan = "connector-scheduler"
)

type (
	trigger struct {
		messenger      messaging.Client
		connectorRepo  repository.ConnectorRepository
		tracer         trace.Tracer
		connectorModel *model.Connector
		fileSizeLimit  int
	}
)

func (t *trigger) Do(ctx context.Context) error {
	// if connector is new or
	// todo we need figure out how to use multiple  orchestrators instances
	// one approach could be that this method will extract top x rows from the database
	// and it will book them
	if t.connectorModel.LastSuccessfulIndexDate.IsZero() ||
		t.connectorModel.LastSuccessfulIndexDate.Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second).Before(time.Now().UTC()) {
		ctx, span := t.tracer.Start(ctx, ConnectorSchedulerSpan)
		defer span.End()
		span.SetAttributes(attribute.Int64(model.SpanAttributeConnectorID, t.connectorModel.ID.IntPart()))
		span.SetAttributes(attribute.String(model.SpanAttributeConnectorSource, string(t.connectorModel.Type)))
		zap.S().Infof("run connector %s", t.connectorModel.ID)

		if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
			span.RecordError(err)
			return err
		}
		connWF, err := connector.New(t.connectorModel)
		if err != nil {
			return err
		}
		if err = connWF.PrepareTask(ctx, t); err != nil {
			span.RecordError(err)
			if errr := t.updateStatus(ctx, model.ConnectorStatusError); errr != nil {
				span.RecordError(errr)
			}
			return err
		}

	}
	return nil
}

// RunChunker send message to chunker service
func (t *trigger) RunChunker(ctx context.Context, data *proto.ChunkingData) error {
	if err := t.updateStatus(ctx, model.ConnectorStatusWorking); err != nil {
		return err
	}
	return t.messenger.Publish(ctx, t.messenger.StreamConfig().ChunkerStreamSubject, data)
}

// RunConnector send message to connector service
func (t *trigger) RunConnector(ctx context.Context, data *proto.ConnectorRequest) error {
	data.Params[connector.ParamFileLimit] = fmt.Sprintf("%d", t.fileSizeLimit)

	if err := t.updateStatus(ctx, model.ConnectorStatusWorking); err != nil {
		return err
	}
	return t.messenger.Publish(ctx, t.messenger.StreamConfig().ConnectorStreamSubject, data)
}
func (t *trigger) UpToDate(ctx context.Context) error {
	return t.updateStatus(ctx, model.ConnectorStatusSuccess)
}

func NewTrigger(messenger messaging.Client,
	connectorRepo repository.ConnectorRepository,
	connectorModel *model.Connector,
	fileSizeLimit int) *trigger {
	return &trigger{
		messenger:      messenger,
		connectorRepo:  connectorRepo,
		connectorModel: connectorModel,
		fileSizeLimit:  fileSizeLimit,
		tracer:         otel.Tracer(model.TracerConnector),
	}
}

// update status of connector in database
func (t *trigger) updateStatus(ctx context.Context, status string) error {
	t.connectorModel.LastAttemptStatus = status
	t.connectorModel.LastUpdate = pg.NullTime{time.Now().UTC()}
	return t.connectorRepo.Update(ctx, t.connectorModel)
}
