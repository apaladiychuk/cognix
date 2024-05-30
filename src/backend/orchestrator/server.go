package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"time"
)

type Server struct {
	renewInterval   time.Duration
	connectorRepo   repository.ConnectorRepository
	messenger       messaging.Client
	scheduleTrigger Trigger
	scheduler       gocron.Scheduler
	streamCfg       *messaging.StreamConfig
}

func NewServer(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	messenger messaging.Client,
	messagingCfg *messaging.Config,
	scheduleTrigger Trigger) (*Server, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Server{connectorRepo: connectorRepo,
		renewInterval:   time.Duration(cfg.RenewInterval) * time.Second,
		messenger:       messenger,
		streamCfg:       messagingCfg.Stream,
		scheduler:       s,
		scheduleTrigger: scheduleTrigger}, nil
}

func (s *Server) run(ctx context.Context) error {
	zap.S().Infof("Schedule reload task")
	go s.schedule()
	zap.S().Infof("Start listener ...")
	go s.listen(context.Background())
	return nil
}

// loadFromDatabase load connectors from database and run if needed
func (s *Server) loadFromDatabase() error {
	ctx := context.Background()
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

func (s *Server) schedule() error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.renewInterval),
		gocron.NewTask(s.loadFromDatabase),
		gocron.WithName("reload from database"),
	)
	if err != nil {
		return err
	}
	s.scheduler.Start()
	return nil

}

// listen nats channel with updated connectors
func (s *Server) listen(ctx context.Context) {

	if err := s.loadFromDatabase(); err != nil {
		return
	}
	//if err := s.messenger.Listen(ctx, model.TopicUpdateConnector, model.SubscriptionOrchestrator, s.handleTriggerRequest); err != nil {
	//	zap.S().Errorf("failed to listen: %v", err)
	//}
}

func (s *Server) scheduleConnector(ctx context.Context, trigger *proto.ConnectorRequest) error {
	conn, err := s.connectorRepo.GetByID(ctx, trigger.GetId())
	if err != nil {
		return err
	}
	return s.scheduleTrigger.Do(ctx, conn)
}
