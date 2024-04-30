package messaging

import (
	"cognix.ch/api/v2/core/utils"
	"encoding/json"
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
	Message struct {
		Header map[string]string `json:"header"`
		Body   json.RawMessage   `json:"body"`
	}

	Subscription struct {
		ch           chan *Message
		subscription *nats.Subscription
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
		return newNatsClient(cfg.Nats)
	case providerPulsar:
		return NewPulsar(cfg.Pulsar)
	}
	return nil, fmt.Errorf("unknown provider %s", cfg.Provider)
}
