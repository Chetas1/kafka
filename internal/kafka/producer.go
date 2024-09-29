package kafka

import (
	"fmt"
	"github/chetasp/kafka/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(config *config.Config) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.Broker,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.username":     config.Kafka.Username,
		"sasl.password":     config.Kafka.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	return &KafkaProducer{producer: p, topic: config.KafkaProducer.Topic}, nil
}

// Produce sends a message to the Kafka topic
func (kp *KafkaProducer) Produce(message string) error {
	kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	kp.producer.Flush(15 * 1000)
	return nil
}

// Close closes the producer
func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
