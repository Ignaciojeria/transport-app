package eventprocessing

import (
	"context"
	"micartapro/app/shared/infrastructure/httpserver"

	"cloud.google.com/go/pubsub/v2"
	ioc "github.com/Ignaciojeria/ioc"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// MessageProcessor is the unified callback signature used by ALL brokers.
//
// Return contract:
//   - status < 500  → ACK  (message handled, do not retry)
//   - status >= 500 → NACK (retry)
//
// `error` is optional and informational.
// Only the HTTP-style status code controls ACK/NACK.
//
// This allows Pub/Sub, NATS JetStream, Kafka, HTTP push, etc.
// to share the same application-level processing logic.
type MessageProcessor func(ctx context.Context, event cloudevents.Event) int

// ProcessorMiddleware allows you to wrap processors with additional behavior.
// Example use cases:
//   - logging
//   - tracing
//   - panic recovery
//   - metrics
//   - tenant extraction
type ProcessorMiddleware func(MessageProcessor) MessageProcessor

// ApplyMiddlewares builds a final MessageProcessor with middleware.
// Middlewares are applied in order:
//
// m1 → m2 → m3 → processor
//
// Example:
//
//	wrapped := ApplyMiddlewares(mp, LoggingMiddleware, TracingMiddleware)
func ApplyMiddlewares(
	processor MessageProcessor,
	middlewares ...ProcessorMiddleware,
) MessageProcessor {

	// Apply in reverse so they wrap correctly.
	for i := len(middlewares) - 1; i >= 0; i-- {
		processor = middlewares[i](processor)
	}
	return processor
}

type Subscriber interface {
	Start(subscriptionName string, processor MessageProcessor, receiveSettings ReceiveSettings) error
}

func init() {
	ioc.Register(NewSubscriberStrategy)
}

func NewSubscriberStrategy(c *pubsub.Client, s httpserver.Server) Subscriber {
	return NewPubSubSubscriber(c, s)
}

// =============================================================
// DomainEvent → Anything that can be converted to CloudEvent
// =============================================================
type DomainEvent interface {
	ToCloudEvent(source string) cloudevents.Event
}

type PublishRequest struct {
	Topic       string
	Source      string
	OrderingKey string // Optional: para control de orden en cola
	Event       DomainEvent
}

// =============================================================
// PublisherManager (vendor-neutral interface)
// =============================================================
type PublisherManager interface {
	Publish(ctx context.Context, request PublishRequest) error
}

func init() {
	ioc.Register(NewPublisherStrategy)
}
func NewPublisherStrategy(c *pubsub.Client) PublisherManager {
	return NewGcpPublisherManager(c)
}
