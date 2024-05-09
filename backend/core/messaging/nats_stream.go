package messaging

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	_ "github.com/nats-io/nats.go/jetstream"
	_ "go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"sync"

	_ "go.opentelemetry.io/otel"
)

type clientStream struct {
	conn                *nats.Conn
	stream              nats.JetStreamContext
	connectorStreamName string
	subscriptions       []*Subscription
	once                sync.Once
}

func (c *clientStream) Close() {
	c.once.Do(func() {
		for _, sub := range c.subscriptions {
			close(sub.ch)
			if err := sub.subscription.Unsubscribe(); err != nil {
				zap.S().Errorf("unsubscribe %s ", err.Error())
			}
		}
	})
}

func (c *clientStream) Publish(ctx context.Context, topic string, body *proto.Body) error {
	message, err := buildMessage(ctx, body)
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

func (c *clientStream) Listen(_ context.Context, topic, subscriptionName string) (<-chan *proto.Message, error) {
	out := make(chan *proto.Message)
	subscription, err := c.stream.Subscribe(fmt.Sprintf("%s.%s", c.connectorStreamName, topic),
		func(msg *nats.Msg) {
			defer func() {
				if err := msg.Ack(); err != nil {
					zap.S().Warnf("Ack failed: %s", err.Error())
				}
			}()
			var message proto.Message
			if err := proto2.Unmarshal(msg.Data, &message); err != nil {
				zap.S().Errorf("Error unmarshalling message: %s", string(msg.Data))
				return
			}

			out <- &message
		})
	if err != nil {
		return nil, err
	}
	c.subscriptions = append(c.subscriptions, &Subscription{
		ch:           out,
		subscription: subscription,
	})
	return out, nil
}

func NewClientStream(cfg *natsConfig) (Client, error) {
	conn, err := nats.Connect(
		cfg.URL,
		nats.Name("Cognix"),
		nats.MaxReconnects(reconnectAttempts),
		nats.ReconnectWait(reconnectWaitTime),
	)
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}
	js, err := conn.JetStream(nats.PublishAsyncMaxPending(streamMaxPending))
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}

	stream, err := js.StreamInfo(cfg.ConnectorStreamName)
	// stream not found, create it
	if stream == nil {
		zap.S().Infof("Creating stream: %s", cfg.ConnectorStreamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     cfg.ConnectorStreamName,
			Subjects: []string{cfg.ConnectorStreamName + ".*"},
		})
		if err != nil {
			zap.S().Errorf("Error creating stream: %s", err.Error())
			return nil, err
		}
	}

	return &clientStream{
		conn:                conn,
		stream:              js,
		connectorStreamName: cfg.ConnectorStreamName,
		once:                sync.Once{},
	}, nil
}

type NopClient struct {
}

func (n *NopClient) Publish(ctx context.Context, topic string, body *proto.Body) error {
	return nil
}

func (n *NopClient) Listen(ctx context.Context, topic, subscriptionName string) (<-chan *proto.Message, error) {
	return make(chan *proto.Message), nil
}

func (n *NopClient) Close() {
	//TODO implement me
	panic("implement me")
}

func NewMopClient() Client {
	return &NopClient{}
}
