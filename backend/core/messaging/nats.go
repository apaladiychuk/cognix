package messaging

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
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

func (c *client) Publish(ctx context.Context, topic string, body *proto.Body) error {
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

func (c *client) Listen(ctx context.Context, topic, subscriptionName string, handler MessageHandler) error {
	out := make(chan *proto.Message)
	subscription, err := c.conn.Subscribe(topic,
		func(msg *nats.Msg) {
			var message proto.Message
			if err := json.Unmarshal(msg.Data, &message); err != nil {
				zap.S().Errorf("Error unmarshalling message: %s", string(msg.Data))
				return
			}
			out <- &message
		})
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-out:
			if err = handler(ctx, msg); err != nil {
				zap.S().Errorf("handler message error: %s", err.Error())
			}
		case <-ctx.Done():
			break
		}
	}
	return subscription.Unsubscribe()
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
