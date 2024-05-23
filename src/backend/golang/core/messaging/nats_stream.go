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
)

type clientStream struct {
	conn                *nats.Conn
	js                  jetstream.JetStream
	stream              jetstream.Stream
	connectorStreamName string
	cancel              context.CancelFunc
	ctx                 context.Context
	wg                  *sync.WaitGroup
}

func (c *clientStream) Close() {
	c.wg.Add(1)
	c.cancel()
	c.wg.Wait()
}

func (c *clientStream) Publish(ctx context.Context, topic string, body *proto.Body) error {
	message, err := buildMessage(ctx, body)
	if err != nil {
		return err
	}
	// todo here we must define
	_, err = c.js.Publish(ctx, fmt.Sprintf("%s.%s", c.connectorStreamName, topic), message)
	//,
	//		nats.AckWait(time.Minute*2)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientStream) Listen(ctx context.Context, topic, subscriptionName string, handler MessageHandler) error {

	//cons, _ := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
	//	Durable:       "subscriptionName",
	//	FilterSubject: fmt.Sprintf("%s.%s", c.connectorStreamName, topic),
	//	AckPolicy:     jetstream.AckExplicitPolicy,
	//	AckWait:       time.Minute,
	//	MaxWaiting:    1,
	//})
	//
	//cons.Consume(func(msg jetstream.Msg) {
	//	zap.S().Infof("Received message: %s %s ", msg.Subject(), msg.Reply())
	//	err := msg.Ack()
	//	if err != nil {
	//		zap.S().Errorf("Error acknowledging message: %s", err.Error())
	//	}
	//	var message proto.Message
	//	if err := proto2.Unmarshal(msg.Data(), &message); err != nil {
	//		zap.S().Errorf("Error unmarshalling message: %s", err.Error())
	//		return
	//	}
	//	if err := handler(ctx, &message); err != nil {
	//		zap.S().Errorf("Error handling message: %s", err.Error())
	//	}
	//	//err := msg.DoubleAck(ctx)
	//
	//})
	_, _ = c.conn.QueueSubscribe(fmt.Sprintf("%s.%s", c.connectorStreamName, topic), "event-processor", func(msg *nats.Msg) {
		zap.S().Infof("Received message: %s", msg.Reply)
		var message proto.Message
		if err := proto2.Unmarshal(msg.Data, &message); err != nil {
			zap.S().Errorf("Error unmarshalling message: %s", err.Error())
			return
		}
		if err := handler(ctx, &message); err != nil {
			zap.S().Errorf("Error handling message: %s", err.Error())
		}
		err := msg.Ack()
		if err != nil {
			zap.S().Errorf("Error acknowledging message: %s", err.Error())
		}
	})

	//cons, _ := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
	//	Durable:       "subscriptionName",
	//	FilterSubject: fmt.Sprintf("%s.%s", c.connectorStreamName, topic),
	//	AckPolicy:     jetstream.AckExplicitPolicy,
	//	AckWait:       10 * time.Minute,
	//	MaxWaiting:    1,
	//})

	//cc, _ := cons.Consume(func(msg jetstream.Msg) {
	//	zap.S().Infof("Received message: %s", msg.Reply)
	//	var message proto.Message
	//	if err := proto2.Unmarshal(msg.Data(), &message); err != nil {
	//		zap.S().Errorf("Error unmarshalling message: %s", err.Error())
	//		return
	//	}
	//	if err := handler(ctx, &message); err != nil {
	//		zap.S().Errorf("Error handling message: %s", err.Error())
	//	}
	//	zap.S().Infof("do ack")
	//	err := msg.Ack()
	//	if err != nil {
	//		zap.S().Errorf("Error acknowledging message: %s", err.Error())
	//	}
	//})
	//for {
	//	select {
	//	case <-ctx.Done():
	//		break
	//	default:
	//
	//	}
	//}
	<-c.ctx.Done()
	//cc.Stop()

	//c.js.DeleteConsumer(ctx, c.connectorStreamName, "subscriptionName")
	zap.S().Info("finish")
	c.wg.Done()
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
		Name:      cfg.ConnectorStreamName,
		Retention: jetstream.WorkQueuePolicy,
		Storage:   jetstream.FileStorage,
		Subjects:  []string{cfg.ConnectorStreamName + ".>"},
	})
	if err != nil {
		zap.S().Errorf("Error creating stream: %s", err.Error())
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &clientStream{
		conn:                conn,
		stream:              stream,
		js:                  js,
		connectorStreamName: cfg.ConnectorStreamName,
		ctx:                 ctx,
		cancel:              cancel,
		wg:                  &sync.WaitGroup{},
	}, nil
}
