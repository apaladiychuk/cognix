package messaging

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const (
	providerNats   = "nats"
	providerPulsar = "pulsar"
)

type (
	Config struct {
		Provider string `env:"MESSAGING_PROVIDER" default:"nats"`
		Nats     *natsConfig
		Pulsar   *pulsarConfig
	}
	natsConfig struct {
		URL                 string `env:"NATS_URL"`
		ConnectorStreamName string `env:"NATS_STREAM_NAME" envDefault:"Connector"`
	}
	Subscription struct {
		ch           chan *proto.Message
		subscription *nats.Subscription
	}
	MessageHandler func(ctx context.Context, msg *proto.Message) error
	Client         interface {
		Publish(ctx context.Context, topic string, body *proto.Body) error
		Listen(ctx context.Context, topic, subscriptionName string, handler MessageHandler) error
		Close()
	}
)

const (
	reconnectAttempts = 120
	reconnectWaitTime = 5 * time.Second
	streamMaxPending  = 256
)

var NatsModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{
			Pulsar: &pulsarConfig{},
			Nats:   &natsConfig{},
		}
		err := utils.ReadConfig(&cfg)
		if err != nil {
			zap.S().Errorf(err.Error())
			return nil, err
		}
		return &cfg, nil
	},
		NewClient,
	),
)

func NewClient(cfg *Config) (Client, error) {
	switch cfg.Provider {
	case providerNats:
		return NewClientStream(cfg.Nats)
	case providerPulsar:
		return NewPulsar(cfg.Pulsar)
	}
	return nil, fmt.Errorf("unknown provider %s", cfg.Provider)
}
