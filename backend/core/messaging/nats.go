package messaging

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	_ "github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	_ "go.opentelemetry.io/otel/propagation"

	_ "go.opentelemetry.io/otel"
)

type Client interface {
	Publish(ctx context.Context, topic string, body interface{}) error
	Listen(ctx context.Context, topic string) (<-chan *Message, error)
}

type client struct {
	conn *nats.Conn
}

func (c *client) Publish(ctx context.Context, topic string, body interface{}) error {
	message := c.buildMessage(ctx, body)
}

func (c *client) Listen(ctx context.Context, topic string) (<-chan *Message, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) buildMessage(ctx context.Context, body interface{}) *Message {
	data := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, &data)
	return &Message{
		Header: data,
		Body:   body,
	}
}

func NewClient(cfg *Config) (Client, error) {
	conn, err := nats.Connect(
		cfg.URL,
		nats.Name("Cognix"),
		nats.MaxReconnects(reconnectAttempts),
		nats.ReconnectWait(reconnectWaitTime),
	)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}
	return &client{
		conn: conn,
	}, nil
}
