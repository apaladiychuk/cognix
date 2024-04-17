package messaging

import (
	"context"
	"github.com/nats-io/nats.go"
)

type Client interface {
	Publish(ctx context.Context, topic string, msg *Message) error
	Listen(ctx context.Context, topic string) (<-chan *Message, error)
}

type client struct {
	conn nuts.Client
}

func (c *client) Publish(ctx context.Context, topic string, msg *Message) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) Listen(ctx context.Context, topic string) (<-chan *Message, error) {
	//TODO implement me
	panic("implement me")
}

func NewClient(cfg *Config) (Client, error) {
	conn, err := nats.Connect(
		cfg.URL,
		nats.Name("Fennec channel"),
		nats.MaxReconnects(reconnectAttempts),
		nats.ReconnectWait(reconnectWaitTime),
	)
	if err != nil {
		return nil, err
	}
	return &client{
		conn: conn,
	}, nil
}
