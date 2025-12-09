package eventprocessing

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/infrastructure/eventprocessing/gcp"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"

	"cloud.google.com/go/pubsub/v2"
)

// =============================================================
// GCP IMPLEMENTATION OF THE EVENT BUS
// =============================================================
type GcpPublisherManager struct {
	client *pubsub.Client
}

// IoC: register this manager + the GCP client
func init() {
	ioc.Registry(
		NewGcpPublisherManager,
		gcp.NewClient,
	)
}

// Constructor for the GCP Publisher (no topic defined here)
func NewGcpPublisherManager(c *pubsub.Client) PublisherManager {
	return &GcpPublisherManager{client: c}
}

// =============================================================
// Publish → dynamic topic, CloudEvents-encoded payload
// =============================================================
func (p *GcpPublisherManager) Publish(
	ctx context.Context,
	topic string,
	evt DomainEvent,
) error {

	// DomainEvent → CloudEvent
	ce := evt.ToCloudEvent()

	// Encode CloudEvent JSON
	bytes, err := json.Marshal(ce)
	if err != nil {
		return err
	}

	// Build attributes from CloudEvent extensions
	attrs := map[string]string{}
	for k, v := range ce.Context.GetExtensions() {
		attrs[k] = fmt.Sprintf("%v", v)
	}

	// Dynamic topic resolution (very cheap)
	pubTopic := p.client.Publisher(topic)

	// Publish message
	_, err = pubTopic.Publish(ctx, &pubsub.Message{
		Data:       bytes,
		Attributes: attrs,
	}).Get(ctx)

	return err
}
