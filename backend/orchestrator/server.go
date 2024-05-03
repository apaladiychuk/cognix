package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"go.uber.org/zap"
)

type Server struct {
	connectorRepo   repository.ConnectorRepository
	messenger       messaging.Client
	scheduleTrigger Trigger
}

func NewServer(connectorRepo repository.ConnectorRepository,
	messenger messaging.Client,
	scheduleTrigger Trigger) *Server {
	return &Server{connectorRepo: connectorRepo,
		messenger:       messenger,
		scheduleTrigger: scheduleTrigger}
}

func (s *Server) run(ctx context.Context) error {

	zap.S().Infof("Start listener ...")
	go s.listen(context.Background())
	return nil
}

func (s *Server) onStart(ctx context.Context) error {
	connectors, err := s.connectorRepo.GetActive(ctx)
	if err != nil {
		return err
	}
	for _, connector := range connectors {
		if err = s.scheduleTrigger.Do(ctx, connector); err != nil {
			zap.S().Errorf("run connector %d failed: %v", connector.ID, err)
		}
	}
	return nil
}

func (s *Server) listen(ctx context.Context) error {

	if err := s.onStart(ctx); err != nil {
		return err
	}

	ch, err := s.messenger.Listen(ctx, model.TopicUpdateConnector, model.SubscriptionOrchestrator)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-ch:
			trigger := msg.GetBody().GetTrigger()
			if trigger == nil {
				zap.S().Errorf("Received message with empty trigger")
				continue
			}
			if err = s.scheduleConnector(context.Background(), trigger); err != nil {
				zap.S().Errorf("error scheduling connector[%d] : %v", trigger.GetId(), err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (s *Server) scheduleConnector(ctx context.Context, trigger *proto.TriggerRequest) error {
	conn, err := s.connectorRepo.GetByID(ctx, trigger.GetId())
	if err != nil {
		return err
	}
	return s.scheduleTrigger.Do(ctx, conn)
}
