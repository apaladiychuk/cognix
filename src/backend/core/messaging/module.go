package messaging

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"context"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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
		Nats *natsConfig
		//Pulsar *pulsarConfig
		Stream *StreamConfig
	}
	natsConfig struct {
		URL string `env:"NATS_URL"`
	}
	// StreamConfig contains variables for configure streams
	StreamConfig struct {
		ConnectorStreamName    string `env:"CONNECTOR_STREAM_NAME,required"`
		ConnectorStreamSubject string `env:"CONNECTOR_STREAM_SUBJECT,required"`
		ChunkerStreamName      string `env:"CHUNKER_STREAM_NAME,required"`
		ChunkerStreamSubject   string `env:"CHUNKER_STREAM_SUBJECT,required"`
	}
	Subscription struct {
		ch           chan *proto.Message
		subscription *nats.Subscription
	}
	MessageHandler func(ctx context.Context, msg jetstream.Msg) error
	Client         interface {
		Publish(ctx context.Context, topic string, body proto2.Message) error
		Listen(ctx context.Context, streamName, topic string, handler MessageHandler) error
		StreamConfig() *StreamConfig
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
			//Pulsar: &pulsarConfig{},
			Nats:   &natsConfig{},
			Stream: &StreamConfig{},
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
	//return newNatsClient(cfg.Nats)
	return NewClientStream(cfg)
	//switch cfg.Provider {
	//case providerNats:
	//	return NewClientStream(cfg.Nats)
	//case providerPulsar:
	//	return NewPulsar(cfg.Pulsar)
	//}
	//return nil, fmt.Errorf("unknown provider %s", cfg.Provider)
}
