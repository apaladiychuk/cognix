package messaging

import (
	"cognix.ch/api/v2/core/utils"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"time"
)

type (
	Config struct {
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
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	},
		NewClient,
	),
)
