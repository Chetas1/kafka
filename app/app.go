// Package app wires producer + consumer + event handler into a single
// runtime. App.Run produces a small burst of test messages, then blocks on
// the consumer until ctx is cancelled.
package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Chetas-Patil/kafka/config"
	"github.com/Chetas-Patil/kafka/internal/consumer"
	kafkahelper "github.com/Chetas-Patil/kafka/internal/kafka"
)

// App is the runtime composition of producer, consumer, and event handler.
type App struct {
	cfg      *config.Config
	producer *kafkahelper.Producer
	consumer *kafkahelper.Consumer
	events   consumer.EventConsumer
}

// InitializeApplication constructs producer + consumer, returning a wrapped
// error and freeing the producer if the consumer fails to construct.
func InitializeApplication(cfg config.Config) (*App, error) {
	producer, err := kafkahelper.NewProducer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("init producer: %w", err)
	}

	cons, err := kafkahelper.NewConsumer(&cfg)
	if err != nil {
		if cerr := producer.Close(); cerr != nil {
			log.Printf("close producer after consumer init failure: %v", cerr)
		}
		return nil, fmt.Errorf("init consumer: %w", err)
	}

	return &App{
		cfg:      &cfg,
		producer: producer,
		consumer: cons,
		events:   consumer.NewEventConsumer(cfg),
	}, nil
}

// Run produces a small burst of demo messages, then runs the consumer loop
// until ctx is cancelled. Returns the first non-cancellation error from
// either producer or consumer.
func (a *App) Run(ctx context.Context) error {
	const demoBurst = 10
	for i := 0; i < demoBurst; i++ {
		message := fmt.Sprintf("Message-%d", i)
		if err := a.producer.Produce(message); err != nil {
			log.Printf("produce %d: %v", i, err)
		}
	}

	consumerErr := make(chan error, 1)
	go func() {
		consumerErr <- a.consumer.Consume(ctx, func(message string) {
			if err := a.events.Process(message); err != nil {
				log.Printf("process: %v", err)
			}
		})
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received")
	case err := <-consumerErr:
		if err != nil && err != context.Canceled {
			return fmt.Errorf("consumer: %w", err)
		}
	}
	return nil
}

// Close flushes the producer and stops the consumer. Errors from each are
// logged; the first is returned so the caller can react.
func (a *App) Close() error {
	var firstErr error
	if err := a.producer.Close(); err != nil {
		log.Printf("close producer: %v", err)
		firstErr = err
	}
	if err := a.consumer.Close(); err != nil {
		log.Printf("close consumer: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
