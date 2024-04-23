package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
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
	if conn.LastSuccessfulIndexTime.IsZero() ||
		conn.LastSuccessfulIndexTime.Add(time.Duration(conn.RefreshFreq)*time.Second).Before(time.Now()) {
		ctx, span := t.tracer.Start(ctx, ConnectorSchedulerSpan)
		span.SetAttributes(attribute.Int64(model.SpanAttributeConnectorID, conn.ID))
		span.SetAttributes(attribute.String(model.SpanAttributeConnectorSource, string(conn.Source)))
		zap.S().Infof("run connector %s", conn.ID)
		return t.messenger.Publish(ctx, model.TopicExecutor,
			connector.Trigger{
				ID: conn.ID,
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
