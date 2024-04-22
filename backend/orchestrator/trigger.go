package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"time"
)

type (
	Trigger interface {
		Do(ctx context.Context) error
	}
	cronTrigger struct {
		connectorRepo repository.ConnectorRepository
		messenger     messaging.Client
	}
)

func (t *cronTrigger) Do(ctx context.Context) error {
	connectors, err := t.connectorRepo.GetActive(ctx)
	if err != nil {
		return err
	}
	for _, connector := range connectors {

	}
}

func (t *cronTrigger) runConnector(ctx context.Context, connector *model.Connector) error {
	// if connector is new or
	if !connector.LastSuccessfulIndexTime.IsZero() &&
		connector.LastSuccessfulIndexTime.Add(time.Duration(connector.RefreshFreq)*time.Second).Before(time.Now()) {

	}

}
func NewCronTrigger(connectorRepo repository.ConnectorRepository,
	messenger messaging.Client) Trigger {
	return &cronTrigger{
		connectorRepo: connectorRepo,
		messenger:     messenger,
	}
}
