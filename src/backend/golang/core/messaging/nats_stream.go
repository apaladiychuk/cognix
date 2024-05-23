package messaging

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	_ "github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"sync"
	"time"
)

type clientStream struct {
	conn                *nats.Conn
	js                  jetstream.JetStream
	stream              jetstream.Stream
	connectorStreamName string
	subscriptions       []*Subscription
	once                sync.Once
	cancel              context.CancelFunc
}

func (c *clientStream) Close() {
	c.cancel()

}

func (c *clientStream) Publish(ctx context.Context, topic string, body *proto.Body) error {
	message, err := buildMessage(ctx, body)
	if err != nil {
		return err
	}
	// todo here we must define
	pubAck, err := c.js.Publish(ctx, fmt.Sprintf("%s.%s", c.connectorStreamName, topic), message)
	//,
	//		nats.AckWait(time.Minute*2)
	if err != nil {
		return err
	}
	zap.S().Infof("Published message with ack: %s - %s - %d %v", pubAck.Stream, pubAck.Domain, pubAck.Sequence, pubAck.Duplicate)
	return nil
}

func (c *clientStream) Listen(_ context.Context, topic, subscriptionName string, handler MessageHandler) error {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	cons, _ := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          subscriptionName,
		FilterSubject: fmt.Sprintf("%s.%s", c.connectorStreamName, topic),
		AckPolicy:     jetstream.AckAllPolicy,
		AckWait:       10 * time.Minute,
		MaxWaiting:    1,
	})
	cc, _ := cons.Consume(func(msg jetstream.Msg) {
		zap.S().Infof("Received message: %s", msg.Reply)
		var message proto.Message
		if err := proto2.Unmarshal(msg.Data(), &message); err != nil {
			zap.S().Errorf("Error unmarshalling message: %s", err.Error())
			return
		}
		if err := handler(ctx, &message); err != nil {
			zap.S().Errorf("Error handling message: %s", err.Error())
		}
		zap.S().Infof("do ack")
		err := msg.Ack()
		if err != nil {
			zap.S().Errorf("Error acknowledging message: %s", err.Error())
		}
	})
	for {
		select {
		case <-ctx.Done():
			break
		default:

		}
	}
	cc.Stop()
	zap.S().Info("finish")
	return nil
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

	js, err := jetstream.New(conn)
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}

	stream, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     cfg.ConnectorStreamName,
		Subjects: []string{cfg.ConnectorStreamName + ".>"},
	})
	if err != nil {
		zap.S().Errorf("Error creating stream: %s", err.Error())
		return nil, err
	}

	return &clientStream{
		conn:                conn,
		stream:              stream,
		js:                  js,
		connectorStreamName: cfg.ConnectorStreamName,
		once:                sync.Once{},
	}, nil
}
