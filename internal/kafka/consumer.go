package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/Chetas-Patil/kafka/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer subscribes to a single topic and dispatches messages to a handler.
type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

// NewConsumer connects to the configured broker and subscribes to the
// configured topic. Errors during connection or subscription are returned
// without leaking a half-initialised *kafka.Consumer.
func NewConsumer(cfg *config.Config) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Broker,
		"group.id":          cfg.KafkaConsumer.Group,
		"auto.offset.reset": "earliest",
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.username":     cfg.Kafka.Username,
		"sasl.password":     cfg.Kafka.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("new kafka consumer: %w", err)
	}

	if err := c.SubscribeTopics([]string{cfg.KafkaConsumer.Topic}, nil); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("subscribe %s: %w", cfg.KafkaConsumer.Topic, err)
	}

	return &Consumer{consumer: c, topic: cfg.KafkaConsumer.Topic}, nil
}

// Consume polls for messages and invokes handler synchronously per message.
// It returns when ctx is cancelled (normal shutdown) or when the consumer is
// closed (returns ErrConsumerClosed).
func (kc *Consumer) Consume(ctx context.Context, handler func(string)) error {
	const pollTimeoutMs = 100
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ev := kc.consumer.Poll(pollTimeoutMs)
		switch e := ev.(type) {
		case *kafka.Message:
			handler(string(e.Value))
		case kafka.Error:
			// Fatal client-side errors propagate up; transient errors are logged
			// upstream by callers (we don't want to spam from a hot loop).
			if e.IsFatal() {
				return fmt.Errorf("kafka consumer fatal: %w", e)
			}
		case nil:
			// Poll timeout — loop and re-check ctx.
		default:
			// Other events (offset commits, partition assignments, ...) are
			// no-ops for this thin wrapper.
		}
	}
}

// Close shuts down the underlying consumer. It is safe to call from a defer.
func (kc *Consumer) Close() error {
	if err := kc.consumer.Close(); err != nil && !errors.Is(err, ErrConsumerClosed) {
		return fmt.Errorf("close consumer: %w", err)
	}
	return nil
}
