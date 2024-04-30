package messaging

import (
	"context"
	"encoding/json"
	"github.com/apache/pulsar-client-go/pulsar"
	"go.uber.org/zap"
	"time"
)

type (
	pulsarConfig struct {
		URL               string `env:"PULSAR_URL"`
		OperationTimeout  int    `env:"OPERATION_TIMEOUT" envDefault:"30"`
		ConnectionTimeout int    `env:"CONNECTION_TIMEOUT" envDefault:"30"`
	}
	pulsarClient struct {
		conn       pulsar.Client
		producers  map[string]pulsar.Producer
		subscriber map[string]pulsar.Consumer
	}
)

func (p *pulsarClient) Publish(ctx context.Context, topic string, body interface{}) error {
	msg, err := buildMessage(ctx, body)
	if err != nil {
		return err
	}
	producer, ok := p.producers[topic]
	if !ok {
		producer, err = p.conn.CreateProducer(pulsar.ProducerOptions{
			Topic: topic,
		})
		if err != nil {
			return err
		}
		p.producers[topic] = producer
	}

	_, err = producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: msg,
	})
	return err
}

func (p *pulsarClient) Listen(ctx context.Context, topic, subscriptionName string) (<-chan *Message, error) {
	consumer, err := p.conn.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subscriptionName,
		Type:             pulsar.Shared,
	})
	if err != nil {
		return nil, err
	}

	msgCh := make(chan *Message, 1)

	go func() {
		defer consumer.Close()
		for {
			// may block here
			select {
			case <-ctx.Done():
				close(msgCh)
				break
			default:

			}
			msg, err := consumer.Receive(ctx)
			if err != nil {
				zap.S().Errorf("Receive message error: %s", err.Error())
				break
			}
			var message Message
			if err := json.Unmarshal(msg.Payload(), &message); err != nil {
				zap.S().Errorf("Error unmarshalling message: %s", string(msg.Payload()))
				continue
			}
			msgCh <- &message
			if err = consumer.Ack(msg); err != nil {
				zap.S().Errorf("Ack message error: %s", err.Error())
			}
		}
		if err = consumer.Unsubscribe(); err != nil {
			zap.S().Errorf("Unsubscribe message error: %s", err.Error())
		}
	}()

	return msgCh, nil
}

func (p *pulsarClient) Close() {
	for _, producer := range p.producers {
		producer.Close()
	}
	p.conn.Close()
}

func NewPulsar(cfg *pulsarConfig) (Client, error) {
	coon, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               cfg.URL,
		OperationTimeout:  time.Duration(cfg.OperationTimeout) * time.Second,
		ConnectionTimeout: time.Duration(cfg.ConnectionTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &pulsarClient{
		conn:       coon,
		producers:  make(map[string]pulsar.Producer),
		subscriber: make(map[string]pulsar.Consumer),
	}, nil
}
