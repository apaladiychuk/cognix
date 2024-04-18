package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	_ "github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	_ "go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"log"

	_ "go.opentelemetry.io/otel"
)

type Client interface {
	Publish(ctx context.Context, topic string, body interface{}) error
	Listen(ctx context.Context, topic string) (<-chan *Message, error)
}

type client struct {
	conn                *nats.Conn
	stream              nats.JetStreamContext
	connectorStreamName string
}

func (c *client) Publish(ctx context.Context, topic string, body interface{}) error {
	message, err := c.buildMessage(ctx, body)
	if err != nil {
		return err
	}
	pubAck, err := c.stream.Publish(fmt.Sprintf("%s.%s", c.connectorStreamName, topic), message)
	if err != nil {
		return err
	}
	zap.S().Infof("Published message with ack: %s", pubAck.Domain)
	return nil
}

func (c *client) Listen(ctx context.Context, topic string) (<-chan *Message, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) buildMessage(ctx context.Context, body interface{}) ([]byte, error) {
	data := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, &data)
	msg := &Message{
		Header: data,
		Body:   body,
	}
	return json.Marshal(msg)
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
	js, err := conn.JetStream(nats.PublishAsyncMaxPending(streamMaxPending))
	if err != nil {
		return nil, err
	}

	stream, err := js.StreamInfo(cfg.ConnectorStreamName)
	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", cfg.ConnectorStreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     cfg.ConnectorStreamName,
			Subjects: []string{cfg.ConnectorStreamName + ".*"},
		})
		if err != nil {
			return nil, err
		}
	}

	return &client{
		conn:                conn,
		stream:              js,
		connectorStreamName: cfg.ConnectorStreamName,
	}, nil
}
