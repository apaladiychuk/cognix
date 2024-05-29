package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
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
	Trigger interface {
		Do(ctx context.Context, conn *model.Connector) error
	}
	cronTrigger struct {
		messenger messaging.Client
		tracer    trace.Tracer
	}
)

func (t *cronTrigger) Do(ctx context.Context, conn *model.Connector) error {
	// if connector is new or
	// todo we need figure out how to use multiple  orchestrators instances
	// one approach could be that this method will extract top x rows from the database
	// and it will book them
	if conn.LastSuccessfulIndexTime.IsZero() ||
		conn.LastSuccessfulIndexTime.Add(time.Duration(conn.RefreshFreq)*time.Second).Before(time.Now().UTC()) {
		ctx, span := t.tracer.Start(ctx, ConnectorSchedulerSpan)
		span.SetAttributes(attribute.Int64(model.SpanAttributeConnectorID, conn.ID.IntPart()))
		span.SetAttributes(attribute.String(model.SpanAttributeConnectorSource, string(conn.Source)))
		zap.S().Infof("run connector %s", conn.ID)
		return t.messenger.Publish(ctx, t.messenger.StreamConfig().ConnectorStreamSubject,
			&proto.ConnectorRequest{
				Id: conn.ID.IntPart(),
			})
	}
	return nil
}

func NewCronTrigger(messenger messaging.Client) Trigger {
	return &cronTrigger{
		messenger: messenger,
		tracer:    otel.Tracer(model.TracerConnector),
	}
}
