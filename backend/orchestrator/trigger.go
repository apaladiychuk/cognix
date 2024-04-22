package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
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
		Do(ctx context.Context) error
	}
	cronTrigger struct {
		connectorRepo repository.ConnectorRepository
		messenger     messaging.Client
		tracer        trace.Tracer
	}
)

func (t *cronTrigger) Do(ctx context.Context) error {
	connectors, err := t.connectorRepo.GetActive(ctx)
	if err != nil {
		return err
	}
	for _, connector := range connectors {
		if err = t.runConnector(ctx, connector); err != nil {
			zap.S().Errorf("run connector %d failed: %v", connector.ID, err)
		}
	}
}

func (t *cronTrigger) runConnector(ctx context.Context, conn *model.Connector) error {
	// if connector is new or
	if !conn.LastSuccessfulIndexTime.IsZero() &&
		conn.LastSuccessfulIndexTime.Add(time.Duration(conn.RefreshFreq)*time.Second).Before(time.Now()) {
		ctx, span := t.tracer.Start(ctx, ConnectorSchedulerSpan)
		span.SetAttributes(attribute.Int64(model.SpanAttributeConnectorID, conn.ID))
		span.SetAttributes(attribute.String(model.SpanAttributeConnectorSource, string(conn.Source)))
		return t.messenger.Publish(ctx, connector.TopicExecutor,
			connector.Trigger{
				ID: conn.ID,
			})
	}

}
func NewCronTrigger(connectorRepo repository.ConnectorRepository,
	messenger messaging.Client) Trigger {
	return &cronTrigger{
		connectorRepo: connectorRepo,
		messenger:     messenger,
		tracer:        otel.Tracer(connector.ConnectorTracer),
	}
}
