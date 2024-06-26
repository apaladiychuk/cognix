package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"time"
)

type Server struct {
	renewInterval time.Duration
	connectorRepo repository.ConnectorRepository
	docRepo       repository.DocumentRepository
	messenger     messaging.Client
	scheduler     gocron.Scheduler
	streamCfg     *messaging.StreamConfig
	cfg           *Config
}

func NewServer(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	messenger messaging.Client,
	messagingCfg *messaging.Config) (*Server, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Server{connectorRepo: connectorRepo,
		docRepo:       docRepo,
		renewInterval: time.Duration(cfg.RenewInterval) * time.Second,
		cfg:           cfg,
		messenger:     messenger,
		streamCfg:     messagingCfg.Stream,
		scheduler:     s,
	}, nil
}

func (s *Server) run(ctx context.Context) error {
	zap.S().Infof("Schedule reload task")
	go s.schedule()
	zap.S().Infof("Start listener ...")
	go s.loadFromDatabase()
	return nil
}

// loadFromDatabase load connectors from database and run if needed
func (s *Server) loadFromDatabase() error {
	ctx := context.Background()
	if !s.messenger.IsOnline() {
		zap.S().Infof("Messenger is offline.")
		return nil
	}
	zap.S().Infof("Loading connectors from db")
	connectors, err := s.connectorRepo.GetActive(ctx)
	if err != nil {
		zap.S().Errorf("Load connectors failed: %v", err)
		return err
	}
	for _, connector := range connectors {
		if err = NewTrigger(s.messenger, s.connectorRepo, s.docRepo, connector, s.cfg.FileSizeLimit, s.cfg.OAuthURL).Do(ctx); err != nil {
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
