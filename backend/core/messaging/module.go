package messaging

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
	"time"
)

type (
	Config struct {
		URL                 string `env:"NUTS_URL"`
		ConnectorStreamName string `env:"NUTS_STREAM_NAME" envDefault:"Connector"`
	}
	Message struct {
		Header map[string]string `json:"header"`
		Body   interface{}       `json:"body"`
	}
)

const (
	reconnectAttempts = 120
	reconnectWaitTime = 5 * time.Second
	streamMaxPending  = 256
)

var NutsModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	},
		NewClient,
	),
)
