package messaging

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
	"log"
	"time"
)

type (
	pulsarConfig struct {
		URL               string        `env:"PULSAR_URL"`
		OperationTimeout  time.Duration `env:"OPERATION_TIMEOUT" envDefault:"30"`
		ConnectionTimeout time.Duration `env:"CONNECTION_TIMEOUT" envDefault:"30"`
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

func (p *pulsarClient) Listen(_ context.Context, topic string) (<-chan *Message, error) {
	consumer, err := p.conn.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: uuid.New().String(),
		Type:             pulsar.Shared,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	for i := 0; i < 10; i++ {
		// may block here
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
			msg.ID(), string(msg.Payload()))

		consumer.Ack(msg)
	}

	if err := consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
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
		OperationTimeout:  cfg.OperationTimeout * time.Second,
		ConnectionTimeout: cfg.ConnectionTimeout * time.Second,
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
