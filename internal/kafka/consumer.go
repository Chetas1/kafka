package kafka

import (
	"fmt"
	"github/chetasp/kafka/config"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config *config.Config) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.Broker,
		"group.id":          config.KafkaConsumer.Group,
		"auto.offset.reset": "earliest",
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.username":     config.Kafka.Username,
		"sasl.password":     config.Kafka.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	c.SubscribeTopics([]string{config.KafkaConsumer.Topic}, nil)
	return &KafkaConsumer{consumer: c, topic: config.KafkaConsumer.Topic}, nil
}

// Consume reads messages from the Kafka topic
func (kc *KafkaConsumer) Consume(handler func(string)) error {
	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err == nil {
			handler(string(msg.Value))
		} else {
			log.Printf("consumer error: %v (%v)", err, msg)
		}
	}
}

// Close closes the consumer
func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
