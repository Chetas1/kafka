// Package kafka wraps confluent-kafka-go (librdkafka) into thin Producer /
// Consumer types backed by config.Config. SASL/PLAIN auth is hard-wired
// because that's what the configured broker expects; switch to SASL_SSL
// for production by overriding the security.protocol map below.
package kafka

import (
	"errors"
	"fmt"
	"time"

	"github.com/Chetas1/kafka/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// flushTimeout bounds Producer.Close so a stalled broker doesn't hang
// graceful shutdown forever.
const flushTimeout = 15 * time.Second

// Producer publishes string messages to a Kafka topic.
type Producer struct {
	producer *kafka.Producer
	topic    string
}

// NewProducer connects to the configured broker. The returned *Producer is
// not safe for concurrent use across goroutines that call Close — callers
// must serialize Close.
func NewProducer(cfg *config.Config) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Broker,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.username":     cfg.Kafka.Username,
		"sasl.password":     cfg.Kafka.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("new kafka producer: %w", err)
	}
	return &Producer{producer: p, topic: cfg.KafkaProducer.Topic}, nil
}

// Produce sends a single message and waits for a delivery report so the caller
// learns of broker-side errors (e.g. unknown topic, auth failure).
func (kp *Producer) Produce(message string) error {
	deliveryCh := make(chan kafka.Event, 1)
	defer close(deliveryCh)

	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryCh)
	if err != nil {
		return fmt.Errorf("produce: %w", err)
	}

	ev := <-deliveryCh
	m, ok := ev.(*kafka.Message)
	if !ok {
		return fmt.Errorf("unexpected delivery event: %T", ev)
	}
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery: %w", m.TopicPartition.Error)
	}
	return nil
}

// Close flushes any in-flight messages (bounded by flushTimeout) and releases
// the underlying producer. Returns a non-nil error if messages were left
// unflushed when the timeout expired.
func (kp *Producer) Close() error {
	remaining := kp.producer.Flush(int(flushTimeout / time.Millisecond))
	kp.producer.Close()
	if remaining > 0 {
		return fmt.Errorf("flush timed out with %d undelivered messages", remaining)
	}
	return nil
}

// ErrConsumerClosed is returned by (*Consumer).Consume when the consumer
// has been closed while reading.
var ErrConsumerClosed = errors.New("consumer closed")
