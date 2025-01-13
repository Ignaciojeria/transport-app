package subscriptionwrapper

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Middleware func(next MessageProcessor) MessageProcessor

type MessageProcessor func(ctx context.Context, m *pubsub.Message) (statusCode int, err error)
