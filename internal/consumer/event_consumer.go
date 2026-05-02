// Package consumer is the application-layer event handler. It receives raw
// strings from the kafka transport package and turns them into business
// actions. Today's implementation just logs; replace `Process` with the real
// handler when the schema stabilizes.
package consumer

import (
	"log"

	"github.com/Chetas-Patil/kafka/config"
)

// EventConsumer is the application-layer interface for processing messages.
type EventConsumer interface {
	Process(message string) error
}

type eventConsumer struct {
	cfg config.Config
}

// NewEventConsumer constructs an EventConsumer bound to the given config.
func NewEventConsumer(cfg config.Config) EventConsumer {
	return &eventConsumer{cfg: cfg}
}

// Process handles a single message. Returning an error stops the consumer
// loop; nil acknowledges the message.
func (e *eventConsumer) Process(message string) error {
	log.Printf("event_consumer: received message %q", message)
	return nil
}
