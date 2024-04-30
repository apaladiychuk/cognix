package messaging

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"sync"
)

type client struct {
	conn                *nats.Conn
	connectorStreamName string
	subscriptions       []*Subscription
	once                sync.Once
}

func (c *client) Close() {
	c.once.Do(func() {
		for _, sub := range c.subscriptions {
			close(sub.ch)
			if err := sub.subscription.Unsubscribe(); err != nil {
				zap.S().Errorf("unsubscribe %s ", err.Error())
			}
		}
	})
}

func (c *client) Publish(ctx context.Context, topic string, body interface{}) error {
	message, err := buildMessage(ctx, body)
	if err != nil {
		return err
	}
	err = c.conn.Publish(topic, message)
	if err != nil {
		return err
	}
	//zap.S().Infof("Published message with ack: %s", pubAck.Domain)
	return nil
}

func (c *client) Listen(_ context.Context, topic, subscriptionName string) (<-chan *Message, error) {
	out := make(chan *Message)
	subscription, err := c.conn.Subscribe(topic,
		func(msg *nats.Msg) {
			var message Message
			if err := json.Unmarshal(msg.Data, &message); err != nil {
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

func newNatsClient(cfg *natsConfig) (Client, error) {
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

	return &client{
		conn:                conn,
		connectorStreamName: cfg.ConnectorStreamName,
		once:                sync.Once{},
	}, nil
}

func buildMessage(ctx context.Context, data interface{}) ([]byte, error) {
	header := make(propagation.MapCarrier)
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, &header)
	msg := &Message{
		Header: header,
		Body:   body,
	}
	return json.Marshal(msg)
}
