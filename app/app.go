package app

import (
	"fmt"
	"github/chetasp/crowdstrike/config"
	"github/chetasp/crowdstrike/internal/consumer"
	kafkahelper "github/chetasp/crowdstrike/internal/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type IApp interface {
	Run(config config.Config) error
}

type App struct {
	config   *config.Config
	producer *kafkahelper.KafkaProducer
	consumer *kafkahelper.KafkaConsumer
	sigchan  chan os.Signal
}

// InitializeApplication initializes the application with the given configuration
func InitializeApplication(cfg config.Config) (*App, error) {
	producer, err := kafkahelper.NewProducer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	consumer, err := kafkahelper.NewConsumer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	app := &App{
		config:   &cfg,
		producer: producer,
		consumer: consumer,
		sigchan:  make(chan os.Signal, 1),
	}

	return app, nil
}

// Run starts the application
func (a *App) Run(config config.Config) error {

	eventConsumer := consumer.NewEventConsumer(config)

	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message-%d", i)
		if err := a.producer.Produce(message); err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}

	signal.Notify(a.sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := a.consumer.Consume(func(message string) {
			eventConsumer.Process(message)
		})
		if err != nil {
			log.Fatalf("Error in consuming messages: %v", err)
		}
	}()

	<-a.sigchan
	fmt.Println("Shutting down...")

	a.producer.Close()
	a.consumer.Close()

	return nil
}
