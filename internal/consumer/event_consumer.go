package consumer

import (
	"fmt"
	"github/chetasp/crowdstrike/config"
)

type EventConsumer interface {
	Process(message string) error
}

type eventConsumer struct {
	config config.Config
}

func NewEventConsumer(config config.Config) EventConsumer {
	return &eventConsumer{
		config: config,
	}
}

func (e *eventConsumer) Process(message string) error {
	fmt.Printf("received kafka messgae %s", message)
	return nil
}
